package ytdl

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/log"
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

type Format struct {
	Itag Itag

	url    string
	s      string
	sig    string
	stream string
	conn   string
	sp     string
}

func parseFormat(input string) (*Format, error) {
	query, err := url.ParseQuery(input)
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
		}
	}
	return &format, nil
}

// ValueForKey gets the format value for a format key, used for filtering
func (f *Format) ValueForKey(key FormatKey) interface{} {
	switch key {
	case FormatItagKey:
		return f.Itag.Itag
	case FormatExtensionKey:
		return f.Itag.Extension
	case FormatResolutionKey:
		return f.Itag.Resolution
	case FormatVideoEncodingKey:
		return f.Itag.VideoEncoding
	case FormatAudioEncodingKey:
		return f.Itag.AudioEncoding
	case FormatAudioBitrateKey:
		return f.Itag.AudioBitrate
	default:
		log.Debug().Msgf("Unknown format key: %v", key)
		return nil
	}
}

func (f *Format) CompareKey(other *Format, key FormatKey) int {
	switch key {
	case FormatResolutionKey:
		res := f.ValueForKey(key).(string)
		res1, res2 := 0, 0
		if res != "" {
			res1, _ = strconv.Atoi(res[0 : len(res)-2])
		}
		res = other.ValueForKey(key).(string)
		if res != "" {
			res2, _ = strconv.Atoi(res[0 : len(res)-2])
		}
		return res1 - res2
	case FormatAudioBitrateKey:
		return f.ValueForKey(key).(int) - other.ValueForKey(key).(int)
	case FormatFPSKey:
		if f.ValueForKey(key) == nil {
			return -1
		} else if other.ValueForKey(key) == nil {
			return 1
		} else {
			a, _ := strconv.Atoi(f.ValueForKey(key).(string))
			b, _ := strconv.Atoi(other.ValueForKey(key).(string))
			return a - b
		}
	default:
		return 0
	}
}
