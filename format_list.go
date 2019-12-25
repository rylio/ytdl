package ytdl

import (
	"bufio"
	"io"
	"sort"
	"strconv"

	"github.com/rs/zerolog/log"
)

// FormatList is a slice of formats with filtering functionality
type FormatList []Format

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
	wrapper := formatsSortWrapper{formats, key}
	if !reverse {
		sort.Stable(wrapper)
	} else {
		sort.Stable(sort.Reverse(wrapper))
	}
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

func (formats *FormatList) add(infos []formatInfo) {
	for _, info := range infos {

		format, ok := newFormat(info.Itag)
		if !ok {
			log.Debug().Msgf("No metadata found for itag: %v, skipping...", info.Itag)
		}

		if info.Cipher != nil {
			meta, err := parseQuery(*info.Cipher)
			if err != nil {
				log.Debug().Err(err).Msgf("Unable to parse cipher: %v", info.Cipher)
			}
			format.meta = meta
		} else {
			format.meta = map[string]string{
				"url": info.URL,
			}
		}

		*formats = append(*formats, format)
	}
}

func (formats *FormatList) parseFormats(rd io.Reader) {
	r := bufio.NewReader(rd)

	for {
		line, err := r.ReadString(',')

		if err == io.EOF {
			break
		}

		formats.parseFormat(line[:len(line)-1])
	}
}

func (formats *FormatList) parseFormat(input string) {
	meta, err := parseQuery(input)

	if err != nil {
		log.Debug().Err(err).Msgf("Unable to parse format: %v", input)
		return
	}

	itag, _ := strconv.Atoi(meta["itag"])
	format, ok := newFormat(itag)
	if !ok {
		log.Debug().Msgf("No metadata found for itag: %v, skipping...", itag)
		return
	}

	format.meta = meta
	*formats = append(*formats, format)
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
