package ytdl

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var formats = FormatList{
	FORMATS[18],
	FORMATS[22],
	FORMATS[34],
	FORMATS[37],
	FORMATS[133],
	FORMATS[139],
}

type formatListTestCase struct {
	Key             FormatKey
	FilterValues    interface{}
	ExpectedFormats FormatList
}

func TestFilter(t *testing.T) {
	cases := []formatListTestCase{
		{
			FormatExtensionKey,
			[]interface{}{"mp4"},
			FormatList{formats[0], formats[1], formats[3], formats[4], formats[5]},
		},
		{
			FormatResolutionKey,
			[]interface{}{"360p", "720p"},
			FormatList{formats[0], formats[1]},
		},
		{
			FormatItagKey,
			[]interface{}{"22", "37"},
			FormatList{formats[1], formats[3]},
		},
		{
			FormatAudioBitrateKey,
			[]interface{}{"96", "128"},
			FormatList{formats[0], formats[2]},
		},
		{
			FormatResolutionKey,
			[]interface{}{""},
			FormatList{formats[5]},
		},
		{
			FormatAudioBitrateKey,
			[]interface{}{"0"},
			FormatList{formats[4]},
		},
		{
			FormatResolutionKey,
			[]interface{}{},
			nil,
		},
	}

	for _, v := range cases {
		f := formats.Filter(v.Key, v.FilterValues.([]interface{}))
		if !reflect.DeepEqual(v.ExpectedFormats, f) {
			t.Error("Format filter test case failed expected", v.ExpectedFormats, "got", f)
		}
	}
}

func TestExtremes(t *testing.T) {

	cases := []formatListTestCase{
		{
			FormatResolutionKey,
			true,
			FormatList{formats[3]},
		},
		{
			FormatResolutionKey,
			false,
			FormatList{formats[5]},
		},
		{
			FormatAudioBitrateKey,
			true,
			FormatList{formats[1], formats[3]},
		},
		{
			FormatAudioBitrateKey,
			false,
			FormatList{formats[4]},
		},
		{
			FormatItagKey,
			true,
			formats,
		},
	}
	for _, v := range cases {
		f := formats.Extremes(v.Key, v.FilterValues.(bool))
		if !reflect.DeepEqual(v.ExpectedFormats, f) {
			t.Error("Format filter test case failed expected", v.ExpectedFormats, "got", f)
		}
	}
}

func TestSubtract(t *testing.T) {
	cases := []formatListTestCase{
		{
			FormatExtensionKey,
			[]interface{}{"mp4"},
			FormatList{formats[2]},
		},
		{
			FormatResolutionKey,
			[]interface{}{"480p", "360p", "240p", ""},
			FormatList{formats[1], formats[3]},
		},
		{
			FormatResolutionKey,
			[]interface{}{},
			formats,
		},
	}
	for _, v := range cases {
		f := formats.Subtract(formats.Filter(v.Key, v.FilterValues.([]interface{})))
		if !reflect.DeepEqual(v.ExpectedFormats, f) {
			t.Error("Format filter test case failed expected", v.ExpectedFormats, "got", f)
		}
	}
}

func TestSort(t *testing.T) {
	cases := []formatListTestCase{
		{
			FormatResolutionKey,
			formats,
			FormatList{
				formats[5],
				formats[4],
				formats[0],
				formats[2],
				formats[1],
				formats[3],
			},
		},
		{
			FormatAudioBitrateKey,
			formats,
			FormatList{
				formats[4],
				formats[5],
				formats[0],
				formats[2],
				formats[1],
				formats[3],
			},
		},
	}

	for _, v := range cases {
		sorted := v.FilterValues.(FormatList).Copy()
		sorted.Sort(v.Key, false)
		if !reflect.DeepEqual(v.ExpectedFormats, sorted) {
			t.Error("FormatList sort failed")
		}
	}
}

func TestCopy(t *testing.T) {
	if !reflect.DeepEqual(formats, formats.Copy()) {
		t.Error("Copying format list failed")
	}
}

func TestParseStreamList(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	file, err := os.Open("fixtures/stream_map.txt")
	require.NoError(err)
	defer file.Close()

	formats := FormatList{}
	formats.parseFormats(file)

	require.Len(formats, 2)
	format := formats[0]
	assert.Equal(22, format.Itag)
	assert.Equal("mp4", format.Extension)
	assert.Equal("720p", format.Resolution)
	assert.Equal("H.264", format.VideoEncoding)
	assert.Equal("aac", format.AudioEncoding)
	assert.Equal(192, format.AudioBitrate)

	meta := format.meta
	assert.Len(meta, 4)
	assert.Equal(`video/mp4; codecs="avc1.64001F, mp4a.40.2"`, meta["type"])
	assert.Len(meta["url"], 769)
}

func TestParseStreamListEmpty(t *testing.T) {
	formats := FormatList{}
	formats.parseFormats(strings.NewReader(""))
	assert.Len(t, formats, 0)
}
