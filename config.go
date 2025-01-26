package configurator

import (
	"os"

	"github.com/goccy/go-yaml"
)

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

var ENV = os.Getenv("ENV")

// LoadConfig loads the configuration object from the provided files, in order. Files are automatically filtered
// using their environment property, and the current ENV value.
func LoadConfig[Cfg any](files ...ConfigFile) *Cfg {
	// The final output, resulting from merging every file together.
	var out Cfg

	for _, file := range files {
		// Ensure the file can be loaded.
		if file.env == ENV || file.env == "" {
			// Allow yaml configuration to import values from their environment.
			expanded := os.ExpandEnv(string(file.file))
			// Assign the fields in the file to the output object. Missing fields will not be replaced, allowing
			// config files to be merged.
			if err := yaml.Unmarshal([]byte(expanded), &out); err != nil {
				panic(err)
			}
		}
	}

	return &out
}
