package ytdl

import (
	"github.com/rs/zerolog"
)

var log = zerolog.Nop()

// SetLogger sets a new logger
func SetLogger(logger zerolog.Logger) {
	log = logger
}
