package ytdl

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Client struct {
	Logger     zerolog.Logger
	HTTPClient *http.Client
}

var DefaultClient = &Client{
	HTTPClient: http.DefaultClient,
	Logger:     zerolog.Nop(),
}
