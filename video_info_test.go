package ytdl

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestVideoInfo(t *testing.T) {
	testCases := map[string]bool{
		"https://www.youtube.com/watch?v=YQHsXMglC9A":            true,
		"https://www.youtube.com/watch?v=H-30B0cqh88":            true,
		"https://www.youtube.com/watch?v=qD8hOJoOGtk":            true,
		"https://www.youtube.com/":                               false,
		"https://www.youtube.com/watch?v=qHGTs1NSB1s":            true,
		"https://www.facebook.com/video.php?v=10153820411888896": false,
	}

	for k, v := range testCases {
		_, err := GetVideoInfo(k)
		if (err != nil && v) || (err == nil && !v) {
			t.Error("Failed test case:", k, err)
		}
	}
}

func TestDownloadVideo(t *testing.T) {
	info, err := GetVideoInfo("https://www.youtube.com/watch?v=FrG4TEcSuRg")
	if err != nil {
		t.Fatal(err)
	}
	format := info.Formats.Worst(FormatResolutionKey)[0]
	err = info.Download(format, ioutil.Discard)
	if err != nil {
		t.Error(err)
	}
}

func TestThumbnail(t *testing.T) {
	info, err := GetVideoInfo("https://www.youtube.com/watch?v=FrG4TEcSuRg")
	if err != nil {
		t.Fatal(err)
	}

	qualities := []ThumbnailQuality{
		ThumbnailQualityDefault,
		ThumbnailQualityHigh,
		ThumbnailQualityMaxRes,
		ThumbnailQualityMedium,
		ThumbnailQualitySD,
	}

	for _, v := range qualities {
		u := info.GetThumbnailURL(v)
		resp, err := http.Get(u.String())
		if err != nil {
			t.Error(err)
		} else if resp.StatusCode != 200 {
			t.Error("Invalid status code", resp.StatusCode, "for", v)
		}
		resp.Body.Close()
	}
}
