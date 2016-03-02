package main

import (
	"fmt"
	"github.com/otium/ytdl"
	"strings"
)

func parseFilter(filterString string) (func(ytdl.FormatList) ytdl.FormatList, error) {

	filterString = strings.TrimSpace(filterString)
	switch filterString {
	case "best", "worst":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.Extremes(ytdl.FormatResolutionKey, filterString == "best").Extremes(ytdl.FormatAudioBitrateKey, filterString == "best")
		}, nil
	case "best-video", "worst-video":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.Extremes(ytdl.FormatResolutionKey, strings.HasPrefix(filterString, "best"))
		}, nil
	case "best-audio", "worst-audio":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.Extremes(ytdl.FormatAudioBitrateKey, strings.HasPrefix(filterString, "best"))
		}, nil
	}
	err := fmt.Errorf("Invalid filter")
	split := strings.SplitN(filterString, ":", 2)
	if len(split) != 2 {
		return nil, err
	}
	key := ytdl.FormatKey(split[0])
	exclude := key[0] == '!'
	if exclude {
		key = key[1:]
	}
	value := strings.TrimSpace(split[1])
	if value == "best" || value == "worst" {
		return func(formats ytdl.FormatList) ytdl.FormatList {
			f := formats.Extremes(key, value == "best")
			if exclude {
				f = formats.Subtract(f)
			}
			return f
		}, nil
	}
	vals := strings.Split(value, ",")
	values := make([]interface{}, len(vals))
	for i, v := range vals {
		values[i] = strings.TrimSpace(v)
	}
	return func(formats ytdl.FormatList) ytdl.FormatList {
		f := formats.Filter(key, values)
		if exclude {
			f = formats.Subtract(f)
		}
		return f
	}, nil
}
