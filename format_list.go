package ytdl

import (
	"fmt"
	"mime"
	"sort"
	"strconv"
)

// FormatList is a slice of formats with filtering functionality
type FormatList []*Format

func (formats FormatList) Filter(key FormatKey, values []interface{}) FormatList {
	var dst FormatList
	for _, v := range values {
		for _, f := range formats {
			if interfaceToString(f.ValueForKey(key)) == interfaceToString(v) {
				dst = append(dst, f)
			}
		}
	}

	return dst
}

func (formats FormatList) Extremes(key FormatKey, best bool) FormatList {
	dst := formats.Copy()
	if len(dst) > 1 {
		dst.Sort(key, best)
		first := dst[0]
		var i int
		for i = 0; i < len(dst)-1; i++ {
			if first.CompareKey(dst[i+1], key) != 0 {
				break
			}
		}
		i++
		dst = dst[0:i]
	}
	return dst
}

func (formats FormatList) Best(key FormatKey) FormatList {
	return formats.Extremes(key, true)
}

func (formats FormatList) Worst(key FormatKey) FormatList {
	return formats.Extremes(key, false)
}

func (formats FormatList) Sort(key FormatKey, reverse bool) {
	var wrapper sort.Interface = formatsSortWrapper{formats, key}
	if reverse {
		wrapper = sort.Reverse(wrapper)
	}
	sort.Stable(wrapper)
}

func (formats FormatList) Subtract(other FormatList) FormatList {
	var dst FormatList
	for _, f := range formats {
		include := true
		for _, f2 := range other {
			if f2.Itag == f.Itag {
				include = false
				break
			}
		}
		if include {
			dst = append(dst, f)
		}
	}
	return dst
}

func (formats FormatList) Copy() FormatList {
	dst := make(FormatList, len(formats))
	copy(dst, formats)
	return dst
}

func (formats *FormatList) addByInfo(info formatInfo, adaptive bool) error {
	var err error
	var format *Format

	itag := getItag(info.Itag)
	if itag == nil {
		return fmt.Errorf("no itag found with number: %v", info.Itag)
	}

	switch {
	case info.Cipher != nil:
		format, err = parseFormat(*info.Cipher, adaptive)
		if err != nil {
			return fmt.Errorf("unable to parse cipher '%v': %w", info.Cipher, err)
		}
		format.Itag = *itag
	case info.SignatureCipher != nil:
		format, err = parseFormat(*info.SignatureCipher, adaptive)
		if err != nil {
			return fmt.Errorf("unable to parse cipher '%v': %w", info.SignatureCipher, err)
		}
		format.Itag = *itag
	default:
		format = &Format{
			Itag: *itag,
			url:  info.URL,
		}
	}
	format.Adaptive = adaptive
	if adaptive && info.Index != nil {
		format.AdaptiveStream = &AdaptiveStream{
			Index:             info.Index,
			Init:              info.Init,
			Codecs:            info.Codecs,
			Bitrate:           info.Bitrate,
			Width:             info.Width,
			Height:            info.Height,
			AudioSamplingRate: info.AudioSampleRate,
			AudioChannels:     info.AudioChannels,
			FrameRate:         strconv.Itoa(info.FPS),
		}
		if info.MimeType != "" {
			var params map[string]string
			format.AdaptiveStream.MimeType, params, err = mime.ParseMediaType(info.MimeType)
			if err != nil {
				return fmt.Errorf("failed to parse mime type: %s", info.MimeType)
			}
			format.AdaptiveStream.Codecs = params["codecs"]
		}
	}

	*formats = append(*formats, format)
	return nil
}

func (formats *FormatList) addByQueryString(input string, adaptive bool) error {
	format, err := parseFormat(input, adaptive)
	if err != nil {
		return err
	}
	format.Adaptive = adaptive
	*formats = append(*formats, format)
	return nil
}

type formatsSortWrapper struct {
	formats FormatList
	key     FormatKey
}

func (s formatsSortWrapper) Len() int {
	return len(s.formats)
}

func (s formatsSortWrapper) Less(i, j int) bool {
	return s.formats[i].CompareKey(s.formats[j], s.key) < 0
}

func (s formatsSortWrapper) Swap(i, j int) {
	s.formats[i], s.formats[j] = s.formats[j], s.formats[i]
}
