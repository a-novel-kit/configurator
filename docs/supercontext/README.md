---
title: Contexts
icon: material-symbols:contextual-token-outline-rounded
category:
  - context
---

# Contexts

Context tools are available under the `supercontext` sub-package.

## Extract

A common context usage is to create derived contexts with extra values, using the std `context.WithValue`.

The `Extract` method is a convenient way to retrieve a value previously stored in the context. It returns specific
errors when the value is not present, and performs type checking automatically.

```go
package main

import (
	"context"
	"github.com/a-novel-kit/configurator/supercontext"
)

func main() {
	value := "bar"
	ctx := context.WithValue(context.Background(), supercontext.CtxKey("foo"), value)

	// Returns an error if the key "foo" is missing in the
	// context, or if the value is not of type string.
	extracted, err := supercontext.Extract[string](ctx, "foo")
}
```
