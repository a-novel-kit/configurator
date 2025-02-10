package configurator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/a-novel-kit/configurator"
)

func TestLoader(t *testing.T) { //nolint:tparallel
	t.Setenv("ENV", "test")
	t.Setenv("VALUE", "tmp")

	jsonData := [][]byte{
		[]byte(`{"key-1": "value-1", "key-2": "value-2"}`),
		[]byte(`{"key-2": "other-value-2", "key-3": "other-value-3"}`),
		[]byte(`{"key-4": "value-4"}`),
	}

	yanlData := [][]byte{
		[]byte("key-1: value-1\nkey-2: value-2"),
		[]byte("key-2: other-value-2\nkey-3: other-value-3"),
		[]byte("key-4: value-4"),
	}

	withEnv := []byte(`{"key-1": "${VALUE}", "key-2": "value-2"}`)

	testCases := []struct {
		name string

		config configurator.LoaderConfig
		files  []configurator.ConfigFile

		expect any
	}{
		{
			name: "FilterENV",

			config: configurator.LoaderConfig{},

			files: []configurator.ConfigFile{
				configurator.NewConfig("foo", jsonData[0]),
				configurator.NewConfig("test", jsonData[1]),
				configurator.NewConfig("bar", jsonData[2]),
			},

			expect: map[string]any{
				"key-2": "other-value-2",
				"key-3": "other-value-3",
			},
		},
		{
			name: "FilterENV/NotFound",

			config: configurator.LoaderConfig{},

			files: []configurator.ConfigFile{
				configurator.NewConfig("foo", jsonData[0]),
				configurator.NewConfig("qux", jsonData[1]),
				configurator.NewConfig("bar", jsonData[2]),
			},

			expect: map[string]any(nil),
		},
		{
			name: "FilterENV/DefaultENV",

			config: configurator.LoaderConfig{},

			files: []configurator.ConfigFile{
				configurator.NewConfig("", jsonData[0]),
				configurator.NewConfig("test", jsonData[1]),
				configurator.NewConfig("bar", jsonData[2]),
			},

			expect: map[string]any{
				"key-1": "value-1",
				"key-2": "other-value-2",
				"key-3": "other-value-3",
			},
		},
		{
			name: "FilterENV/DefaultENV/OrderMatters",

			config: configurator.LoaderConfig{},

			files: []configurator.ConfigFile{
				configurator.NewConfig("test", jsonData[1]),
				configurator.NewConfig("", jsonData[0]),
				configurator.NewConfig("bar", jsonData[2]),
			},

			expect: map[string]any{
				"key-1": "value-1",
				"key-2": "value-2",
				"key-3": "other-value-3",
			},
		},
		{
			name: "OverrideEnv",

			config: configurator.LoaderConfig{Env: "foo"},

			files: []configurator.ConfigFile{
				configurator.NewConfig("foo", jsonData[0]),
				configurator.NewConfig("test", jsonData[1]),
				configurator.NewConfig("bar", jsonData[2]),
			},

			expect: map[string]any{
				"key-1": "value-1",
				"key-2": "value-2",
			},
		},
		{
			name: "CustomDeserializer",

			config: configurator.LoaderConfig{Deserializer: yaml.Unmarshal},

			files: []configurator.ConfigFile{
				configurator.NewConfig("foo", yanlData[0]),
				configurator.NewConfig("test", yanlData[1]),
				configurator.NewConfig("bar", yanlData[2]),
			},

			expect: map[string]any{
				"key-2": "other-value-2",
				"key-3": "other-value-3",
			},
		},
		{
			name: "ExpandENV/False",

			config: configurator.LoaderConfig{},

			files: []configurator.ConfigFile{
				configurator.NewConfig("", withEnv),
			},

			expect: map[string]any{
				"key-1": "${VALUE}",
				"key-2": "value-2",
			},
		},
		{
			name: "ExpandENV/True",

			config: configurator.LoaderConfig{ExpandEnv: true},

			files: []configurator.ConfigFile{
				configurator.NewConfig("", withEnv),
			},

			expect: map[string]any{
				"key-1": "tmp",
				"key-2": "value-2",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			loader := configurator.NewLoader[map[string]any](testCase.config)

			data, err := loader.Load(testCase.files...)
			require.NoError(t, err)
			require.Equal(t, testCase.expect, data)
		})
	}
}
