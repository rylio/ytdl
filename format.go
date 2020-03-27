package ytdl

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// FormatKey is a string type containing a key in a video format map
type FormatKey string

// Available format Keys
const (
	FormatExtensionKey     FormatKey = "ext"
	FormatResolutionKey    FormatKey = "res"
	FormatVideoEncodingKey FormatKey = "videnc"
	FormatAudioEncodingKey FormatKey = "audenc"
	FormatItagKey          FormatKey = "itag"
	FormatAudioBitrateKey  FormatKey = "audbr"
	FormatFPSKey           FormatKey = "fps"
)

type Range struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Format struct {
	Itag
	Adaptive bool
	// FromDASH indicates that the stream
	// was extracted from the DASH manifest file
	FromDASH bool
	Index    *Range
	Init     *Range
	url      string
	s        string
	sig      string
	stream   string
	conn     string
	sp       string
}

func parseFormat(queryString string) (*Format, error) {
	query, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	format := Format{}

	for k, v := range query {
		switch k {
		case "itag":
			i, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse itag param: %w", err)
			}

			itag := getItag(i)
			if itag == nil {
				return nil, fmt.Errorf("no metadata found for itag: %v", i)
			}

			format.Itag = *itag
		case "url":
			format.url = v[0]
		case "s":
			format.s = v[0]
		case "sig":
			format.sig = v[0]
		case "stream":
			format.stream = v[0]
		case "conn":
			format.conn = v[0]
		case "sp":
			format.sp = v[0]
		case "index":
			format.Index, err = parseRange(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse index range")
			}
		case "init":
			format.Init, err = parseRange(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse init range")
			}
		}
	}
	return &format, nil
}

func parseRange(s string) (*Range, error) {
	sa := strings.Split(s, "-")
	if len(sa) != 2 {
		return nil, fmt.Errorf("Invalid range")
	}
	return &Range{Start: sa[0], End: sa[1]}, nil
}

// ValueForKey gets the format value for a format key, used for filtering
func (f *Format) ValueForKey(key FormatKey) interface{} {
	switch key {
	case FormatItagKey:
		return f.Itag.Number
	case FormatExtensionKey:
		return f.Extension
	case FormatResolutionKey:
		return f.Resolution
	case FormatVideoEncodingKey:
		return f.VideoEncoding
	case FormatAudioEncodingKey:
		return f.AudioEncoding
	case FormatAudioBitrateKey:
		return f.AudioBitrate
	default:
		return fmt.Errorf("Unknown format key: %v", key)
	}
}

func (f *Format) CompareKey(other *Format, key FormatKey) int {
	switch key {
	case FormatResolutionKey:
		return f.resolution() - other.resolution()
	case FormatAudioBitrateKey:
		return f.AudioBitrate - other.AudioBitrate
	case FormatFPSKey:
		return f.FPS - other.FPS
	default:
		return 0
	}
}

// width in pixels
func (f *Format) resolution() int {
	res := f.Itag.Resolution
	if len(res) < 2 {
		return 0
	}

	width, _ := strconv.Atoi(res[:len(res)-2])
	return width
}
