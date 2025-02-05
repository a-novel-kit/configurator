---
title: Home
icon: material-symbols:home-outline-rounded
---

# Configurator - Project setup files

```bash
go get -u github.com/a-novel-kit/configurator
```

This repository contains all the handlers used in A-Novel project, not big enough to justify their own packages.

## Config loader

It is common to have multiple configuration files for different environments in a project. While those are simple to
handle, it can become cumbersome. This is a simple interface to automate the process.

Given the following project tree.

```plaintext
.
└── config
    ├── config.prod.json
    ├── config.staging.json
    ├── config.json
    └── config.go
```

Using the following code:

```go
package config

import (
	_ "embed"
	"github.com/a-novel-kit/configurator"
)

//go:embed config.prod.json
var prodConfigFile []byte

//go:embed config.staging.json
var stagingConfigFile []byte

//go:embed config.json
var commonConfigFile []byte

type ConfigType struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
	/* ... */
}

var configLoader = configurator.NewLoader[ConfigType](configurator.LoaderConfig{})

var Config = configLoader.MustLoad(
	// Is always loaded. Place it first, so it doesn't override
	// specific environment configurations.
	configurator.NewConfig("", commonConfigFile),
	// Will only load if ENV=production.
	configurator.NewConfig("production", prodConfigFile),
	// Will only load if ENV=staging.
	configurator.NewConfig("staging", stagingConfigFile),
)
```

You can customize the loader to fit your requirements. Since the configuration is agnostic to the type of value
represented in the file, you can share a global `LoaderConfig` across all your loaders.

```go
package config

import (
	"github.com/a-novel-kit/configurator"
	"gopkg.in/yaml.v3"
)

var Loader = configurator.LoaderConfig{
	// Interpolate detected environment variables in the
	// file with their actual values. Disabled by default.
	ExpandEnv: true,
	// If you are a thug and don't put the ENV information
	// in the "ENV" variable, you can pass a custom env
	// value.
	Env: "custom-env-key",
	// It is common to use alternative file formats for
	// configuration, such as YAML. You can change the
	// default JSON decoder here.
	Deserializer: yaml.Unmarshal,
}
```
