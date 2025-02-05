package configurator

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoaderConfig sets a file loader.
type LoaderConfig struct {
	// Deserializer used to unmarshal the file content. Defaults to json.Unmarshal.
	Deserializer func(data []byte, dst any) error

	// If true, env variable declarations are automatically interpolated with their
	// actual value from the environment.
	ExpandEnv bool

	// The current environment. Uses the value in "ENV" variable by default.
	Env string
}

// ConfigFile represents an embedded .yaml file, that holds the configuration object.
type ConfigFile struct {
	// The content of the file, unformatted. We keep a []byte value so we can perform unmarshalling on it.
	file []byte
	// The environment associated with the file. A file with a given environment will only be loaded if the ENV
	// variable matches it. Files with an empty environment will always be loaded.
	env string
}

// NewConfig creates a new ConfigFile object, targeted to a specific environment. To create a global configuration file,
// leave the environment empty.
func NewConfig(env string, file []byte) ConfigFile {
	return ConfigFile{file, env}
}

// Loader is an interface used to load configuration files using environment information.
type Loader[T any] struct {
	config LoaderConfig
}

func (loader *Loader[T]) Load(files ...ConfigFile) (T, error) {
	// The final output, resulting from merging every file together.
	var out T

	for _, file := range files {
		// Ensure the file can be loaded.
		if file.env != loader.config.Env && file.env != "" {
			continue
		}

		// Interpolate environment variables with their actual value.
		data := file.file
		if loader.config.ExpandEnv {
			data = []byte(os.ExpandEnv(string(file.file)))
		}

		// Assign the fields in the file to the output object. Missing fields will not be replaced, allowing
		// config files to be merged.
		if err := loader.config.Deserializer(data, &out); err != nil {
			return out, fmt.Errorf("(Loader.Load) unmarshal file: %w", err)
		}
	}

	return out, nil
}

func (loader *Loader[T]) MustLoad(files ...ConfigFile) T {
	out, err := loader.Load(files...)
	if err != nil {
		panic(err)
	}

	return out
}

// NewLoader creates a new Loader object.
func NewLoader[T any](config LoaderConfig) *Loader[T] {
	if config.Deserializer == nil {
		config.Deserializer = json.Unmarshal
	}

	if config.Env == "" {
		config.Env = os.Getenv("ENV")
	}

	return &Loader[T]{config}
}
