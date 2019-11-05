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
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ":            true,
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

func TestGetDownloadURL(t *testing.T) {
	testCases := []string{
		"https://www.youtube.com/watch?v=FrG4TEcSuRg",
		"https://www.youtube.com/watch?v=jgVhBThJdXc",
		"https://www.youtube.com/watch?v=MXgnIP4rMoI",
		"https://www.youtube.com/watch?v=peBgUMT26jM",
		"https://www.youtube.com/watch?v=aQZDbBGBJsM",
		"https://www.youtube.com/watch?v=cRS4mS4gKwg",
		"https://www.youtube.com/watch?v=0fllyJTBsRU",
	}
	for _, url := range testCases {
		info, err := GetVideoInfo(url)
		if err != nil {
			t.Fatal(err)
		}
		format := info.Formats.Worst(FormatResolutionKey)[0]
		_, err = info.GetDownloadURL(format)
		if err != nil {
			t.Error("Failed test case:", url, err)
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
