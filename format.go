package ytdl

import (
	"fmt"
	"mime"
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

func (r *Range) String() string {
	return r.Start + "-" + r.End
}

type Format struct {
	Itag
	Adaptive       bool
	AdaptiveStream *AdaptiveStream
	// FromDASH indicates that the stream
	// was extracted from the DASH manifest file
	FromDASH bool
	url      string
	s        string
	sig      string
	stream   string
	conn     string
	sp       string
}

// Type AdaptiveStream represents an adaptive stream
type AdaptiveStream struct {
	Index    *Range
	Init     *Range
	MimeType string
	Codecs   string
	Bitrate  int
	//Video stream specific
	Width     int
	Height    int
	FrameRate string
	// Audio stream specific
	AudioSamplingRate string
	AudioChannels     int
}

func parseFormat(queryString string, adaptive bool) (*Format, error) {
	query, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	format := Format{}
	if adaptive && format.AdaptiveStream == nil {
		format.AdaptiveStream = new(AdaptiveStream)
	}
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
			format.AdaptiveStream.Index, err = parseRange(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse index range")
			}
		case "init":
			format.AdaptiveStream.Init, err = parseRange(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse init range")
			}
		case "bitrate":
			format.AdaptiveStream.Bitrate, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse bitrate")
			}
		case "fps":
			format.AdaptiveStream.FrameRate = v[0]
		case "size":
			sa := strings.Split(v[0], "x")
			if len(sa) != 2 {
				return nil, fmt.Errorf("unable to parse size")
			}
			format.AdaptiveStream.Height, err = strconv.Atoi(sa[1])
			if err != nil {
				return nil, fmt.Errorf("unable to parse size")
			}
			format.AdaptiveStream.Width, err = strconv.Atoi(sa[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse size")
			}
		case "type":
			if format.AdaptiveStream == nil {
				continue
			}
			// type=video/mp4;+codecs="avc1.4d401e"
			var params map[string]string
			format.AdaptiveStream.MimeType, params, err = mime.ParseMediaType(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse media type")
			}
			format.AdaptiveStream.Codecs = params["codecs"]
		case "audio_channels":
			if format.AdaptiveStream == nil {
				continue
			}
			format.AdaptiveStream.AudioChannels, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse audio_channels")
			}
		case "audio_sample_rate":
			if format.AdaptiveStream == nil {
				continue
			}
			format.AdaptiveStream.AudioSamplingRate = v[0]
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
