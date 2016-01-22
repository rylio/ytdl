package ytdl

import (
	"reflect"
	"testing"
)

var formats = []Format{
	FORMATS[18],
	FORMATS[22],
	FORMATS[34],
	FORMATS[37],
	FORMATS[133],
	FORMATS[139],
}

type FormatFilterTestCase struct {
	Key             FormatKey
	FilterValues    interface{}
	ExpectedFormats []Format
}

func TestFilterFormatsInclude(t *testing.T) {

	cases := []FormatFilterTestCase{
		FormatFilterTestCase{
			FormatExtensionKey,
			[]string{"mp4"},
			[]Format{formats[0], formats[1], formats[3], formats[4], formats[5]},
		},
		FormatFilterTestCase{
			FormatResolutionKey,
			[]string{"360p", "720p"},
			[]Format{formats[0], formats[1]},
		},
		FormatFilterTestCase{
			FormatItagKey,
			[]string{"22", "37"},
			[]Format{formats[1], formats[3]},
		},
		FormatFilterTestCase{
			FormatAudioBitrateKey,
			[]string{"96", "128"},
			[]Format{formats[0], formats[2]},
		},
		FormatFilterTestCase{
			FormatResolutionKey,
			[]string{""},
			[]Format{formats[5]},
		},
		FormatFilterTestCase{
			FormatAudioBitrateKey,
			[]string{"0"},
			[]Format{formats[4]},
		},
		FormatFilterTestCase{
			FormatResolutionKey,
			[]string{},
			nil,
		},
	}
	for _, v := range cases {
		f := FilterFormats(formats, v.Key, v.FilterValues.([]string))
		if !reflect.DeepEqual(v.ExpectedFormats, f) {
			t.Error("Format filter test case failed expected", v.ExpectedFormats, "got", f)
		}
	}
}
func TestFilterFormatsExclude(t *testing.T) {
	cases := []FormatFilterTestCase{
		FormatFilterTestCase{
			FormatExtensionKey,
			[]string{"mp4"},
			[]Format{formats[2]},
		},
		FormatFilterTestCase{
			FormatResolutionKey,
			[]string{"480p", "360p", "240p", ""},
			[]Format{formats[1], formats[3]},
		},
		FormatFilterTestCase{
			FormatResolutionKey,
			[]string{},
			formats,
		},
	}
	for _, v := range cases {
		f := FilterFormatsExclude(formats, v.Key, v.FilterValues.([]string))
		if !reflect.DeepEqual(v.ExpectedFormats, f) {
			t.Error("Format filter test case failed expected", v.ExpectedFormats, "got", f)
		}
	}
}

func TestFilterFormatsExtremes(t *testing.T) {

	cases := []FormatFilterTestCase{
		FormatFilterTestCase{
			FormatResolutionKey,
			true,
			[]Format{formats[3]},
		},
		FormatFilterTestCase{
			FormatResolutionKey,
			false,
			[]Format{formats[5]},
		},
		FormatFilterTestCase{
			FormatAudioBitrateKey,
			true,
			[]Format{formats[1], formats[3]},
		},
		FormatFilterTestCase{
			FormatAudioBitrateKey,
			false,
			[]Format{formats[4]},
		},
		FormatFilterTestCase{
			FormatItagKey,
			true,
			formats,
		},
	}

	for _, v := range cases {
		f := FilterFormatsExtremes(formats, v.Key, v.FilterValues.(bool))
		if !reflect.DeepEqual(v.ExpectedFormats, f) {
			t.Error("Format filter test case failed expected", v.ExpectedFormats, "got", f)
		}
	}
}
