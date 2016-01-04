package ytdl

import "fmt"

func reverseStringSlice(s []string) {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func interfaceToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
