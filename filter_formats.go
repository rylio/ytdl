package ytdl

import "strconv"

// FilterFormats filters out all formats whose key doesn't contain
// any of values. Formats are ordered by values
func FilterFormats(formats []Format, key FormatKey, values []string) []Format {
	if len(values) == 0 {
		return formats
	}
	filtered := []Format{}
	// filter values first for priority
	for _, value := range values {
		for _, format := range formats {
			v := format.ValueForKey(key)
			val := convertToString(v)
			if val == "" {
				val = "nil"
			}
			if value == val {
				filtered = append(filtered, format)
			}
		}
	}
	return filtered
}

func convertToString(val interface{}) string {
	switch val.(type) {
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(val.(int64), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(val.(uint64), 10)
	case string:
		return val.(string)
	default:
		return ""

	}
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
			val := convertToString(v)
			if val == "" {
				val = "nil"
			}
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

func FilterFormatsExtremes(formats []Format, key FormatKey, high bool) Format {
	return Format{}
}
