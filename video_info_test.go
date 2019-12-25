package ytdl

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoInfo(t *testing.T) {

	tests := []struct {
		url         string
		assertion   assert.BoolAssertionFunc
		duration    time.Duration
		published   time.Time
		title       string
		uploader    string
		description string
		song        string
		artist      string
	}{
		{
			url:       "https://www.youtube.com/",
			assertion: assert.False,
		},
		{
			url:       "https://www.facebook.com/video.php?v=10153820411888896",
			assertion: assert.False,
		},
		//{
		//	url: "https://www.youtube.com/watch?v=TDgn8k9uyW4",
		//},
		{
			url:         "https://www.youtube.com/watch?v=BaW_jenozKc",
			assertion:   assert.True,
			title:       `youtube-dl test video "'/\ä↭𝕐`,
			uploader:    "Philipp Hagemeister",
			duration:    time.Second * 10,
			published:   newDate(2012, 10, 2),
			description: "test chars:  \"'/\\ä↭𝕐\ntest URL: https://github.com/rg3/youtube-dl/iss...\n\nThis is a test video for youtube-dl.\n\nFor more information, contact phihag@phihag.de .",
		},
		{
			url:         "https://www.youtube.com/watch?v=YQHsXMglC9A",
			assertion:   assert.True,
			title:       "Adele - Hello",
			uploader:    "AdeleVEVO",
			duration:    time.Second * 367,
			published:   newDate(2015, 10, 22),
			artist:      "Adele",
			song:        "Hello",
			description: "‘Hello' is taken from the new album, 25, out November 20. http://adele.com\nAvailable now from iTunes http://smarturl.it/itunes25 \nAvailable now from Amazon http://smarturl.it/25amazon \nAvailable now from Google Play http://smarturl.it/25gplay\nAvailable now at Target (US Only): http://smarturl.it/target25\n\nDirected by Xavier Dolan, @XDolan\n\nFollow Adele on:\n\nFacebook - https://www.facebook.com/Adele\nTwitter - https://twitter.com/Adele \nInstagram - http://instagram.com/Adele\n\nhttp://vevo.ly/jzAuJ1\n\nCommissioner: Phil Lee\nProduction Company: Believe Media/Sons of Manual/Metafilms\nDirector: Xavier Dolan\nExecutive Producer: Jannie McInnes\nProducer: Nancy Grant/Xavier Dolan\nCinematographer:  André Turpin\nProduction design : Colombe Raby\nEditor: Xavier Dolan\nAdele's lover : Tristan Wilds",
		},
		{
			url:       "https://www.youtube.com/watch?v=H-30B0cqh88",
			assertion: assert.True,
			title:     "Kung Fu Panda 3 Official Trailer #3 (2016) - Jack Black, Angelina Jolie Animated Movie HD",
			uploader:  "Movieclips Trailers",
			duration:  time.Second * 145,
			published: newDate(2015, 12, 16),

			description: "Subscribe to TRAILERS: http://bit.ly/sxaw6h\nSubscribe to COMING SOON: http://bit.ly/H2vZUn\nLike us on FACEBOOK: http://bit.ly/1QyRMsE\nFollow us on TWITTER: http://bit.ly/1ghOWmt\nKung Fu Panda 3 Official International Trailer #1 (2016) - Jack Black, Angelina Jolie Animation HD\n\nIn 2016, one of the most successful animated franchises in the world returns with its biggest comedy adventure yet, KUNG FU PANDA 3. When Po's long-lost panda father suddenly reappears, the reunited duo travels to a secret panda paradise to meet scores of hilarious new panda characters. But when the supernatural villain Kai begins to sweep across China defeating all the kung fu masters, Po must do the impossible - learn to train a village full of his fun-loving, clumsy brethren to become the ultimate band of Kung Fu Pandas!\n\nThe Fandango MOVIECLIPS Trailers channel is your destination for the hottest new trailers the second they drop. Whether it's the latest studio release, an indie horror flick, an evocative documentary, or that new RomCom you've been waiting for, the Fandango MOVIECLIPS team is here day and night to make sure all the best new movie trailers are here for you the moment they're released.\n\nIn addition to being the #1 Movie Trailers Channel on YouTube, we deliver amazing and engaging original videos each week. Watch our exclusive Ultimate Trailers, Showdowns, Instant Trailer Reviews, Monthly MashUps, Movie News, and so much more to keep you in the know.\n\nHere at Fandango MOVIECLIPS, we love movies as much as you!",
		},
		// Test VEVO video with age protection
		// https://github.com/ytdl-org/youtube-dl/issues/956
		{
			url:         "https://www.youtube.com/watch?v=07FYdnEawAQ",
			assertion:   assert.True,
			title:       `Justin Timberlake - Tunnel Vision (Official Music Video) (Explicit)`,
			uploader:    "justintimberlakeVEVO",
			duration:    time.Second * 419,
			published:   newDate(2013, 7, 3),
			song:        "Tunnel Vision",
			artist:      "Justin Timberlake",
			description: "Executive Producer: Jeff Nicholas \nProduced by Jonathan Craven and Nathan Scherrer \nDirected by Jonathan Craven, Simon McLoughlin and Jeff Nicholas for The Uprising Creative (http://theuprisingcreative.com) \nDirector Of Photography: Sing Howe Yam \nEditor: Jacqueline London\n\nOfficial music video by Justin Timberlake performing Tunnel Vision (Explicit). (C) 2013 RCA Records, a division of Sony Music Entertainment\n\n#JustinTimberlake #TunnelVision #Vevo #Pop #OfficialMuiscVideo",
		},
		{
			url:       "https://www.youtube.com/watch?v=qHGTs1NSB1s",
			assertion: assert.True,
			title:     "Why Linus Torvalds doesn't use Ubuntu or Debian",
			uploader:  "TFiR: Open Source and Emerging Tech",
			description: `Subscribe to our weekly newsletter: https://www.tfir.io/dnl
Become a patron of this channel: https://www.patreon.com/TFIR
Follow us on Twitter: https://twitter.com/tfir_io
Like us on Facebook: https://www.facebook.com/TFiRMedia/

Linus gives the practical reasons why he doesn't use Ubuntu or Debian.`,
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
				assert.Equal(t, tt.uploader, info.Uploader)
				assert.Equal(t, tt.song, info.Song)
				assert.Equal(t, tt.artist, info.Artist)
				assert.Equal(t, tt.description, info.Description)
			}
		})
	}
}

func TestGetDownloadURL(t *testing.T) {
	testCases := []string{
		"FrG4TEcSuRg",
		"jgVhBThJdXc",
		"MXgnIP4rMoI",
		"peBgUMT26jM",
		"aQZDbBGBJsM",
		"cRS4mS4gKwg",
		"0fllyJTBsRU",
	}
	for _, id := range testCases {
		t.Run(id, func(t *testing.T) {
			info, err := GetVideoInfo("https://www.youtube.com/watch?v=" + id)
			require.NoError(t, err)

			if len(info.Formats) == 0 {
				t.Fatal("empty format list")
			}

			format := info.Formats.Worst(FormatResolutionKey)[0]
			_, err = info.GetDownloadURL(format)
			assert.NoError(t, err)
		})
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
