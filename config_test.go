package configurator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/configurator"
)

func TestLoadConfig(t *testing.T) {
	configFiles := map[string][]byte{
		"item1":     []byte("item1: foo"),
		"item1Alt":  []byte("item1: qux"),
		"item2":     []byte("item2: bar"),
		"bothItems": []byte("item1: one\nitem2: two"),
	}

	type config struct {
		Item1 string `yaml:"item1"`
		Item2 string `yaml:"item2"`
	}

	testCases := []struct {
		name string

		env   string
		files []configurator.ConfigFile

		expect config
	}{
		{
			name: "DevEnv",
			env:  configurator.DevENV,
			files: []configurator.ConfigFile{
				configurator.GlobalConfig(configFiles["bothItems"]),
				// Override field in first entry.
				configurator.DevConfig(configFiles["item1"]),
				// Ignored.
				configurator.StagingConfig(configFiles["item1Alt"]),
			},
			expect: config{
				Item1: "foo",
				Item2: "two",
			},
		},
		{
			name: "StagingEnv",
			env:  configurator.StagingEnv,
			files: []configurator.ConfigFile{
				configurator.GlobalConfig(configFiles["bothItems"]),
				// Override field in first entry.
				configurator.StagingConfig(configFiles["item1Alt"]),
				// Ignored.
				configurator.DevConfig(configFiles["item1"]),
			},
			expect: config{
				Item1: "qux",
				Item2: "two",
			},
		},
		{
			name: "ProdEnv",
			env:  configurator.ProdENV,
			files: []configurator.ConfigFile{
				configurator.GlobalConfig(configFiles["bothItems"]),
				// Override field in first entry.
				configurator.ProdConfig(configFiles["item1Alt"]),
				// Ignored.
				configurator.DevConfig(configFiles["item1"]),
			},
			expect: config{
				Item1: "qux",
				Item2: "two",
			},
		},

		{
			name: "NoDefaultValue",
			env:  configurator.DevENV,
			files: []configurator.ConfigFile{
				configurator.DevConfig(configFiles["item1"]),
			},
			expect: config{Item1: "foo"},
		},
		{
			name: "FileOrderMatters",
			env:  configurator.DevENV,
			files: []configurator.ConfigFile{
				configurator.DevConfig(configFiles["item1"]),
				configurator.GlobalConfig(configFiles["bothItems"]),
			},
			expect: config{
				Item1: "one",
				Item2: "two",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configurator.ENV = tc.env

			cfg := configurator.LoadConfig[config](tc.files...)

			require.Equal(t, tc.expect, *cfg)
		})
	}
}
