package main

import (
	"reflect"
	"testing"
)

type filterTestCase struct {
	filterString   string
	key            string
	expectedValues []string
	exclude        bool
}

func TestParseFilter(t *testing.T) {
	cases := []filterTestCase{
		filterTestCase{
			"res:1080p",
			"res",
			[]string{"1080p"},
			false,
		},
	}
	for _, v := range cases {
		key, res, exclude, err := parseFilter(v.filterString)
		if err != nil {
			t.Error(err)
		} else {
			if key != v.key || !reflect.DeepEqual(res, v.expectedValues) || exclude != v.exclude {
				t.Error("Failed parsing filter", v.filterString)
			}
		}
	}
}
