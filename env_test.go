package configurator_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/configurator"
)

func TestUnsetEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, configurator.ENV, configurator.DevENV)
			require.False(t, configurator.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
	})
}

func TestDevEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, configurator.ENV, configurator.DevENV)
			require.False(t, configurator.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
		Env: []string{"ENV=dev"},
	})
}

func TestStagingEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, configurator.ENV, configurator.StagingEnv)
			require.True(t, configurator.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
		Env: []string{"ENV=staging"},
	})
}

func TestProdEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, configurator.ENV, configurator.ProdENV)
			require.True(t, configurator.IsReleaseEnv())
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
		Env: []string{"ENV=prod"},
	})
}

func TestInvalidEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			fmt.Println(configurator.ENV)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
		},
		Env: []string{"ENV=foo"},
	})
}
