package ytdl

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVideoInfo(t *testing.T) {

	tests := []struct {
		url       string
		assertion assert.BoolAssertionFunc
		duration  time.Duration
		published time.Time
		title     string
		author    string
	}{
		{
			url:       "https://www.youtube.com/",
			assertion: assert.False,
		},
		{
			url:       "https://www.facebook.com/video.php?v=10153820411888896",
			assertion: assert.False,
		},
		{
			url:       "https://www.youtube.com/watch?v=BaW_jenozKc",
			assertion: assert.True,
			title:     `youtube-dl test video "'/\ä↭𝕐`,
			author:    "Philipp Hagemeister",
			duration:  time.Second * 10,
			published: newDate(2012, 10, 2),
		},
		{
			url:       "https://www.youtube.com/watch?v=YQHsXMglC9A",
			assertion: assert.True,
			title:     "Adele - Hello",
			author:    "AdeleVEVO",
			duration:  time.Second * 367,
			published: newDate(2015, 10, 22),
		},
		{
			url:       "https://www.youtube.com/watch?v=H-30B0cqh88",
			assertion: assert.True,
			title:     "Kung Fu Panda 3 Official Trailer #3 (2016) - Jack Black, Angelina Jolie Animated Movie HD",
			author:    "Movieclips Trailers",
			duration:  time.Second * 145,
			published: newDate(2015, 12, 16),
		},
		{
			url:       "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			assertion: assert.True,
			title:     `Rick Astley - Never Gonna Give You Up (Video)`,
			author:    "RickAstleyVEVO",
			duration:  time.Second * 212,
			published: newDate(2009, 10, 24),
		},
		{
			url:       "https://www.youtube.com/watch?v=qHGTs1NSB1s",
			assertion: assert.True,
			title:     "Why Linus Torvalds doesn't use Ubuntu or Debian",
			author:    "TFiR: Open Source and Emerging Tech",
			duration:  time.Second * 162,
			published: newDate(2014, 9, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			info, err := GetVideoInfo(tt.url)

			tt.assertion(t, err == nil)

			if err == nil {
				assert.Equal(t, tt.duration, info.Duration)
				assert.Equal(t, tt.title, info.Title)
				assert.Equal(t, tt.published, info.DatePublished)
				assert.Equal(t, tt.author, info.Author)
			}
		})
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

func newDate(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}
