package ytdl

import (
	"net/http"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// newTestClient creates a new Client for the test case
func newTestClient(t *testing.T) Client {
	logger := zerolog.New(&testWriter{t: t}).Level(zerolog.DebugLevel).With().CallerWithSkipFrameCount(2).Logger()

	return Client{
		HTTPClient: http.DefaultClient,
		Logger:     logger,
	}
}

type testWriter struct {
	t *testing.T
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.t.Log(strings.TrimSpace(string(p)))
	return len(p), nil
}
