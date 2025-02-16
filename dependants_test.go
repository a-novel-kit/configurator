package configurator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/configurator"
)

func TestResolveDependants(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		mods map[string][]string
		deps map[string][]string

		expect    map[string][]string
		expectErr error
	}{
		{
			name: "Linear",

			mods: map[string][]string{
				"mod:1": {"dep:1", "dep:2"},
				"mod:2": {"dep:3"},
				"mod:3": {"dep:4"},
			},
			deps: map[string][]string{
				"mod:1": {},
				"mod:2": {},
				"mod:3": {},
			},

			expect: map[string][]string{
				"mod:1": {"dep:1", "dep:2"},
				"mod:2": {"dep:3"},
				"mod:3": {"dep:4"},
			},
		},
		{
			name: "Dependants",

			mods: map[string][]string{
				"mod:1": {"dep:1", "dep:2"}, // Depth: 1 -> mod:6
				"mod:2": {"dep:3"},          // Depth: 0
				"mod:3": {"dep:4"},          // Depth: 3 -> mod:4
				"mod:4": {"deps:5"},         // Depth: 2 -> mod:1(d:1), mod:2(d:0)
				"mod:5": {"deps:6"},         // Depth: 1 -> mod:6
				"mod:6": {"deps:7"},         // Depth: 0
			},
			deps: map[string][]string{
				"mod:1": {"mod:6"},
				"mod:2": {},
				"mod:3": {"mod:4"},
				"mod:4": {"mod:1", "mod:2"},
				"mod:5": {"mod:6"},
				"mod:6": {},
			},

			expect: map[string][]string{
				"mod:1": {"dep:1", "dep:2", "deps:7"},                             // Depth: 1 -> mod:6
				"mod:2": {"dep:3"},                                                // Depth: 0
				"mod:3": {"dep:4", "deps:5", "dep:1", "dep:2", "deps:7", "dep:3"}, // Depth: 3 -> mod:4
				"mod:4": {"deps:5", "dep:1", "dep:2", "deps:7", "dep:3"},          // Depth: 2 -> mod:1(d:1), mod:2(d:0)
				"mod:5": {"deps:6", "deps:7"},                                     // Depth: 1 -> mod:6
				"mod:6": {"deps:7"},                                               // Depth: 0
			},
		},
		{
			name: "Circular/Direct",

			mods: map[string][]string{
				"mod:1": {"dep:1", "dep:2"},
				"mod:2": {"dep:3"},
			},
			deps: map[string][]string{
				"mod:1": {"mod:2"},
				"mod:2": {"mod:1"},
			},

			expectErr: configurator.ErrCircularDependency,
		},
		{
			name: "Circular/HighLevelSeparation",

			mods: map[string][]string{
				"mod:1": {"dep:1", "dep:2"},
				"mod:2": {"dep:3"},
				"mod:3": {"dep:4"},
			},
			deps: map[string][]string{
				"mod:1": {"mod:2"},
				"mod:2": {"mod:3"},
				"mod:3": {"mod:1"},
			},

			expectErr: configurator.ErrCircularDependency,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			resolved, err := configurator.ResolveDependants(testCase.mods, testCase.deps)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resolved)
		})
	}
}
