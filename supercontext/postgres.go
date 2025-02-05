package supercontext

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

const (
	// PostgresDNSEnv is the env variable that should contain the Postgres connection string (Data Source Name or DSN).
	//
	// The DSN string should be in the format:
	//
	//	postgres://[user]:[password]@[host]:[port]/[database]?[...pgoptions]
	PostgresDNSEnv = "DSN"

	pgContextKey CtxKey = "pg"
)

var ErrNoDNS = errors.New("missing DSN environment variable")

// NewPGContext returns a new context with a shared PG connection.
func NewPGContext(ctx context.Context, migrations *embed.FS) (context.Context, error) {
	// Get the URL connection string from the environment.
	dsn := os.Getenv(PostgresDNSEnv)
	if dsn == "" {
		return nil, ErrNoDNS
	}

	// Open a connection to the database.
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Make a temporary assignation. If something goes wrong, it is unnecessary and misleading to assign a value
	// to the global variable.
	client := bun.NewDB(sqldb, pgdialect.New())

	// Wait for connection to be established.
	err := client.Ping()

	// Sometimes the database takes some time to be available, like when it is started at the same time as the
	// server (this usually happens under development / testing environments).
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = client.Ping()
	}

	// Database failed to provide a valid response in the expected time.
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if migrations != nil {
		// Apply migrations.
		mig := migrate.NewMigrations()
		if err = mig.Discover(migrations); err != nil {
			return nil, fmt.Errorf("discover mig: %w", err)
		}

		migrator := migrate.NewMigrator(client, mig)
		if err := migrator.Init(context.Background()); err != nil {
			return nil, fmt.Errorf("create migrator: %w", err)
		}

		if _, err = migrator.Migrate(context.Background()); err != nil {
			return nil, fmt.Errorf("apply mig: %w", err)
		}
	}

	ctxPG := context.WithValue(ctx, pgContextKey, bun.IDB(client))
	// Close clients on context termination.
	context.AfterFunc(ctxPG, func() {
		_ = client.Close()
		_ = sqldb.Close()
	})

	// Use type wrap so the assertion works properly when we try to extract the value.
	return ctxPG, nil
}

// NewPGContextTX creates a new Postgres context where the virtual database is replaced by a transaction.
//
// The parent context MUST contain a virtual Postgres instance, either created with NewPGContext or NewPGContextTX.
//
// The returned cancel function MUST be called with the argument set to true for the transaction to be committed.
// Omitting this call will automatically roll back the whole transaction.
func NewPGContextTX(ctx context.Context, opts *sql.TxOptions) (context.Context, func(commit bool) error, error) {
	pg, err := Extract[bun.IDB](ctx, pgContextKey)
	if err != nil {
		return nil, nil, fmt.Errorf("extract pg: %w", err)
	}

	tx, err := pg.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("begin tx: %w", err)
	}

	ctxTx, cancelFn := context.WithCancel(context.WithValue(ctx, pgContextKey, bun.IDB(tx)))
	context.AfterFunc(ctxTx, func() {
		// If context is canceled without calling the cancel function, abort.
		// If the cancel function was already called, this will return an error,
		// so we ignore it.
		_ = tx.Rollback()
	})

	cancelFnAugmented := func(commit bool) error {
		defer cancelFn()
		if commit {
			return tx.Commit()
		}
		return tx.Rollback()
	}

	return ctxTx, cancelFnAugmented, nil
}

// PGContext extracts the Postgres virtual database from the context.
func PGContext(ctx context.Context) (bun.IDB, error) {
	return Extract[bun.IDB](ctx, pgContextKey)
}
