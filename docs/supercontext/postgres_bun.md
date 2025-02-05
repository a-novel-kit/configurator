---
title: PostgreSQL Context (Bun)
icon: devicon-plain:postgresql-wordmark
category:
  - context
---

# PostgreSQL Context (Bun)

Wraps a PostgreSQL connection in a context. Built around the [bun ORM](https://bun.uptrace.dev/).

```go
package main

import (
	"context"
	"github.com/a-novel-kit/configurator/supercontext"
)

func main() {
	// The "DSN" environment variable MUST be set and point to
	// a valid PostgreSQL database.
	ctx, _ := supercontext.NewPGContext(context.Background(), nil)

	// ...

	// Retrieve the database from anyone method that has access to the context.
	// Returns the generic bun.IDB interface.
	db, _ := supercontext.PGContext(ctx)
}
```

## Auto-migrations

You may also pass postgres [migration files](https://bun.uptrace.dev/guide/migrations.html) to the context creator,
so they are automatically applied.

_Project tree_

```plaintext
.
├── migrations
│   ├── 0001_initial.up.sql
│   ├── 0001_initial.down.sql
│   ├── ...
│   └── migrations.go
└── main.go
```

_migrations.go_

```go
package migrations

import "embed"

//go:embed *.sql
var Migrations embed.FS
```

_main.go_

```go
package main

import (
	"context"
	"github.com/a-novel-kit/configurator/supercontext"
	"github.com/myorg/mypkg/migrations"
)

func main() {
	// Migrations are automatically applied when the context is created.
	ctx, _ := supercontext.NewPGContext(context.Background(), &migrations.Migrations)
}
```

## Controlling transactions

One of the worst struggles, when building layered architectures, is handling transactions. You either handle
them outside the DAO layer (and leak DAO logic in your app), or are stuck without the ability to run multiple
dependant DAO methods in a single transaction.

The TX context solves this issue. It creates a new alternative context that handles the transaction object. Because
everything is located within the context, you don't need to worry about the actual implementation in upper layers.

Downside is, you have to manually cancel the context to commit the transaction. The plus side is, rollback happens
automatically when a parent context is canceled.

```go
package main

import (
	"context"
	"github.com/a-novel-kit/configurator/supercontext"
)

func main() {
	// A transactional context MUST be created from a parent context
	// that embeds a database connection.
	parentCTX, _ := supercontext.NewPGContext(context.Background(), nil)

	// ...

	// Here, you can open a transaction without dealing with DB-related logic.
	transactionalCTX, cancel, _ := supercontext.NewPGContextTX(parentCTX, nil)
	// Recommended to prevent transactions from handling. If you commit anyway,
	// the error returned from this call will be ignored.
    defer cancel(false)

	// ...

	// Methods that require the database will just use your context as always
	// (this is the benefit of using bun.IDB interface, rather than concrete types).
	db, _ := supercontext.PGContext(transactionalCTX)

	// ...

	// You'll have to manually commit the transaction.
	_ = cancel(true)
}
```
