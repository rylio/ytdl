package ytdl

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
			val := interfaceToString(v)
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
