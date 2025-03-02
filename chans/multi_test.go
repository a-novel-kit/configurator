package chans_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/configurator/chans"
)

func TestMultiChan(t *testing.T) {
	t.Parallel()

	multiChan := chans.NewMultiChan[string]()
	defer multiChan.Close()

	listener1 := multiChan.Register()
	listener2 := multiChan.Register()

	multiChan.Send("Hello")

	require.Equal(t, "Hello", <-listener1)
	require.Equal(t, "Hello", <-listener2)

	multiChan.Unregister(listener1)

	multiChan.Send("World")

	require.Equal(t, "World", <-listener2)
}
