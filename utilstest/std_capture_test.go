package utilstest_test

import (
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/configurator/utilstest"
)

func TestCaptureSTD(t *testing.T) {
	t.Parallel()

	patchedSTD, restore, err := utilstest.MonkeyPatchStderr()
	require.NoError(t, err)
	defer restore()

	logs, closer, err := utilstest.CaptureSTD(patchedSTD)
	require.NoError(t, err)
	defer utilstest.RequireCloser(t, closer)

	waiter := utilstest.WaitForLog(
		logs,
		utilstest.WithCaptureRegexpLog(regexp.MustCompile(`^foo bar$`)),
		time.Second,
	)

	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	log.Println("foo bar")

	res, err := waiter()
	require.NoError(t, err)
	require.Equal(t, "foo bar", res)
}
