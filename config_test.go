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
			name: "Target",
			env:  "test",
			files: []configurator.ConfigFile{
				configurator.NewConfig("", configFiles["bothItems"]),
				// Override field in first entry.
				configurator.NewConfig("test", configFiles["item1"]),
				// Ignored.
				configurator.NewConfig("foo", configFiles["item1Alt"]),
			},
			expect: config{
				Item1: "foo",
				Item2: "two",
			},
		},

		{
			name: "NoDefaultValue",
			env:  "test",
			files: []configurator.ConfigFile{
				configurator.NewConfig("test", configFiles["item1"]),
			},
			expect: config{Item1: "foo"},
		},
		{
			name: "FileOrderMatters",
			env:  "test",
			files: []configurator.ConfigFile{
				configurator.NewConfig("test", configFiles["item1"]),
				configurator.NewConfig("", configFiles["bothItems"]),
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
