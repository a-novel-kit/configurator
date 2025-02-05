package supercontext_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/configurator/supercontext"
)

func TestExtract(t *testing.T) {
	value := "bar"
	ctx := context.WithValue(context.Background(), supercontext.CtxKey("foo"), value)

	t.Run("OK", func(t *testing.T) {
		extracted, err := supercontext.Extract[string](ctx, "foo")
		require.NoError(t, err)
		require.Equal(t, value, extracted)
	})

	t.Run("WrongType", func(t *testing.T) {
		_, err := supercontext.Extract[int](ctx, "foo")
		require.ErrorIs(t, err, supercontext.ErrUnsupportedContext)
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := supercontext.Extract[string](ctx, "bar")
		require.ErrorIs(t, err, supercontext.ErrUnsupportedContext)
	})
}
