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

	"github.com/a-novel-kit/configurator/chans"
)

const captureLogsSize = 512

// Read the content of the input. Each time a new line is encountered, send the current line to the channel output
// and start reading the next line. Returns the number of bytes read.
func readLines(in []byte, out chan<- string) int {
	var read int

	for i, c := range in {
		if c == '\n' {
			out <- string(in[read:i])

			read = i + 1
		}
	}

	return read
}

// MonkeyPatchStderr replaces stderr with a pipe. Returns the read end of the pipe, a function to restore stderr and an
// error if something went wrong.
func MonkeyPatchStderr() (*os.File, func(), error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("create pipe: %w", err)
	}

	stderr := os.Stderr
	os.Stderr = w

	return r, func() { os.Stderr = stderr }, nil
}

// CaptureSTD forwards push to the output into a listening channel.
// It returns a function to close the listening channel. An error is returned if initialization failed.
func CaptureSTD(file *os.File) (*chans.MultiChan[string], func() error, error) {
	var (
		err error
		n   int
	)

	output := chans.NewMultiChan[string]()

	// Cursor used to read lines.
	cursor := 0

	go func() {
		logs := make([]byte, 0, captureLogsSize)

		for {
			// Read next entry from the file.
			n, err = file.Read(logs[len(logs):cap(logs)])
			logs = logs[:len(logs)+n]

			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
					cursor += readLines(logs[cursor:], output.Chan())
					// Read the remaining part of the logs.
					output.Send(string(logs[cursor:]))
				} else {
					err = fmt.Errorf("read logs: %w", err)
				}

				return
			}

			if len(logs) == cap(logs) {
				// Add more capacity (let append pick how much).
				logs = append(logs, 0)[:len(logs)]
			}

			cursor += readLines(logs[cursor:], output.Chan())
		}
	}()

	return output, func() error {
		output.Close()

		return err
	}, nil
}

// RequireCloser requires the cleanup function returned by CaptureSTD to succeed.
func RequireCloser(t *testing.T, closer func() error) {
	t.Helper()
	require.NoError(t, closer())
}

type LogCaptureFN func(log string) bool

func WaitForLog(logs *chans.MultiChan[string], capture LogCaptureFN, timeout time.Duration) func() (string, error) {
	timer := time.NewTimer(timeout)
	listener := logs.Register()

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
			case log, ok := <-listener:
				if !ok {
					err = errors.New("channel closed")

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
		logs.Unregister(listener)

		return output, err
	}
}

func WithCaptureRegexpLog(re *regexp.Regexp) LogCaptureFN {
	return func(log string) bool {
		return re.MatchString(log)
	}
}
