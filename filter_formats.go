package ytdl

import "sort"

// FilterFormats filters out all formats whose key doesn't contain
// any of values. Formats are ordered by values
func FilterFormats(formats []Format, key FormatKey, values []string) []Format {
	if len(values) == 0 {
		return nil
	}
	filtered := []Format{}
	// filter values first for priority
	for _, value := range values {
		for _, format := range formats {
			v := format.ValueForKey(key)
			val := interfaceToString(v)
			if value == val {
				filtered = append(filtered, format)
			}
		}
	}
	return filtered
}

// FilterFormatsExclude excludes all formats whose passed key
// contains any of the passed values
func FilterFormatsExclude(formats []Format, key FormatKey, values []string) []Format {
	if len(values) == 0 {
		return formats
	}
	filtered := []Format{}
	for _, format := range formats {
		exclude := false
		for _, value := range values {
			v := format.ValueForKey(key)
			val := interfaceToString(v)
			if val == value {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, format)
		}
	}
	return filtered
}

func FilterFormatsExtremes(formats []Format, key FormatKey, high bool) []Format {
	cp := make([]Format, len(formats))
	copy(cp, formats)
	formats = cp
	sortedContainer := sortFormats{formats, key}
	if high {
		sort.Stable(sort.Reverse(sortedContainer))
	} else {
		sort.Stable(sortedContainer)
	}
	formats = sortedContainer.formats
	if len(formats) > 1 {
		first := formats[0]
		var i int
		for i = 0; i < len(formats)-1; i++ {
			if first.CompareKey(formats[i+1], key) != 0 {
				break
			}
		}
		i++
		formats = formats[0:i]
	}
	return formats
}

type sortFormats struct {
	formats []Format
	key     FormatKey
}

func (s sortFormats) Len() int {
	return len(s.formats)
}

func (s sortFormats) Less(i, j int) bool {
	return s.formats[i].CompareKey(s.formats[j], s.key) < 0
}

func (s sortFormats) Swap(i, j int) {
	s.formats[i], s.formats[j] = s.formats[j], s.formats[i]
}
