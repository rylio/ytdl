package ytdl

import (
	"fmt"
	"sort"
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

func (formats *FormatList) addByInfo(info formatInfo) error {
	var err error
	var format *Format

	itag := getItag(info.Itag)
	if itag == nil {
		return fmt.Errorf("no itag found with number: %v", info.Itag)
	}

	if info.Cipher != nil {
		format, err = parseFormat(*info.Cipher)
		if err != nil {
			return fmt.Errorf("unable to parse cipher '%v': %w", info.Cipher, err)
		}
		format.Itag = *itag
	} else {
		format = &Format{
			Itag: *itag,
			url:  info.URL,
		}
	}

	*formats = append(*formats, format)
	return nil
}

func (formats *FormatList) addByQueryString(input string) error {
	format, err := parseFormat(input)

	if err != nil {
		return err
	}

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
