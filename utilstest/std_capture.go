package utilstest

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func readLines(in []byte, out chan string) int {
	var read int

	for i, c := range in {
		if c == '\n' {
			out <- string(in[read:i])

			read = i + 1
		}
	}

	return read
}

func MonkeyPatchStderr() (*os.File, func(), error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("create pipe: %w", err)
	}

	stderr := os.Stderr
	os.Stderr = w

	return r, func() { os.Stderr = stderr }, nil
}

func CaptureSTD(file *os.File) (chan string, func() error, error) {
	var (
		err error
		n   int
	)

	outC := make(chan string)

	// Cursor used to read lines.
	cursor := 0

	go func() {
		logs := make([]byte, 0, 512)

		for {
			// Read next entry from the file.
			n, err = file.Read(logs[len(logs):cap(logs)])
			logs = logs[:len(logs)+n]

			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
					cursor += readLines(logs[cursor:], outC)
					// Read the remaining part of the logs.
					outC <- string(logs[cursor:])
				} else {
					err = fmt.Errorf("read logs: %w", err)
				}

				return
			}

			if len(logs) == cap(logs) {
				// Add more capacity (let append pick how much).
				logs = append(logs, 0)[:len(logs)]
			}

			cursor += readLines(logs[cursor:], outC)
		}
	}()

	return outC, func() error {
		close(outC)

		return err
	}, nil
}

func RequireCloser(t *testing.T, closer func() error) {
	t.Helper()
	require.NoError(t, closer())
}

type LogCaptureFN func(log string) bool

func WaitForLog(logs chan string, capture LogCaptureFN, timeout time.Duration) func() (string, error) {
	timer := time.NewTimer(timeout)

	var (
		output string
		err    error
	)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-timer.C:
				err = errors.New("timeout")

				return
			case log, ok := <-logs:
				if !ok {
					err = errors.New("log channel closed")

					return
				}

				if capture(log) {
					output = log

					return
				}
			}
		}
	}()

	return func() (string, error) {
		wg.Wait()

		return output, err
	}
}

func WithCaptureRegexpLog(re *regexp.Regexp) LogCaptureFN {
	return func(log string) bool {
		return re.MatchString(log)
	}
}
