package ytdl

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const youtubeBaseURL = "https://www.youtube.com/watch"
const youtubeEmbededBaseURL = "https://www.youtube.com/embed/"
const youtubeVideoEURL = "https://youtube.googleapis.com/v/"
const youtubeVideoInfoURL = "https://www.youtube.com/get_video_info"
const youtubeDateFormat = "2006-01-02"

// VideoInfo contains the info a youtube video
type VideoInfo struct {
	ID             string     // The video ID
	Title          string     // The video title
	Description    string     // The video description
	DatePublished  time.Time  // The date the video was published
	Formats        FormatList // Formats the video is available in
	Keywords       []string   // List of keywords associated with the video
	Uploader       string     // Author of the video
	Song           string
	Artist         string
	Writers        string
	Duration       time.Duration // Duration of the video
	htmlPlayerFile string
}

// GetVideoInfo fetches info from a url string, url object, or a url string
func GetVideoInfo(value interface{}) (*VideoInfo, error) {
	switch t := value.(type) {
	case *url.URL:
		return GetVideoInfoFromURL(t)
	case string:
		u, err := url.ParseRequestURI(t)
		if err != nil {
			return GetVideoInfoFromID(t)
		}
		if u.Host == "youtu.be" {
			return GetVideoInfoFromShortURL(u)
		}
		return GetVideoInfoFromURL(u)
	default:
		return nil, fmt.Errorf("Identifier type must be a string, *url.URL, or []byte")
	}
}

// GetVideoInfoFromURL fetches video info from a youtube url
func GetVideoInfoFromURL(u *url.URL) (*VideoInfo, error) {
	videoID := u.Query().Get("v")
	if len(videoID) == 0 {
		return nil, fmt.Errorf("Invalid youtube url, no video id")
	}
	return GetVideoInfoFromID(videoID)
}

// GetVideoInfoFromShortURL fetches video info from a short youtube url
func GetVideoInfoFromShortURL(u *url.URL) (*VideoInfo, error) {
	if len(u.Path) >= 1 {
		if path := u.Path[1:]; path != "" {
			return GetVideoInfoFromID(path)
		}
	}
	return nil, errors.New("Could not parse short URL")
}

// GetVideoInfoFromID fetches video info from a youtube video id
func GetVideoInfoFromID(id string) (*VideoInfo, error) {
	u, _ := url.ParseRequestURI(youtubeBaseURL)
	values := u.Query()
	values.Set("v", id)
	u.RawQuery = values.Encode()

	body, err := httpGetAndCheckResponseReadBody(u.String())

	if err != nil {
		return nil, err
	}
	return getVideoInfoFromHTML(id, body)
}

// GetDownloadURL gets the download url for a format
func (info *VideoInfo) GetDownloadURL(format Format) (*url.URL, error) {
	return getDownloadURL(format, info.htmlPlayerFile)
}

// GetThumbnailURL returns a url for the thumbnail image
// with the given quality
func (info *VideoInfo) GetThumbnailURL(quality ThumbnailQuality) *url.URL {
	u, _ := url.Parse(fmt.Sprintf("http://img.youtube.com/vi/%s/%s.jpg",
		info.ID, quality))
	return u
}

// Download is a convenience method to download a format to an io.Writer
func (info *VideoInfo) Download(format Format, dest io.Writer) error {
	u, err := info.GetDownloadURL(format)
	if err != nil {
		return err
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}
	_, err = io.Copy(dest, resp.Body)
	return err
}

var (
	regexpPlayerConfig          = regexp.MustCompile("ytplayer.config = (.*?);ytplayer.load")
	regexpInitialData           = regexp.MustCompile(`\["ytInitialData"\] = (.+);`)
	regexpInitialPlayerResponse = regexp.MustCompile(`\["ytInitialPlayerResponse"\] = (.+);`)
)

func getVideoInfoFromHTML(id string, html []byte) (*VideoInfo, error) {

	info := &VideoInfo{}

	if matches := regexpInitialData.FindSubmatch(html); len(matches) > 0 {
		data := initialData{}

		if err := json.Unmarshal(matches[1], &data); err != nil {
			return nil, err
		}

		contents := data.Contents.TwoColumnWatchNextResults.Results.Results.Contents

		if len(contents) >= 2 {
			info.Description = contents[1].VideoSecondaryInfoRenderer.Description.String()
			rows := contents[1].VideoSecondaryInfoRenderer.MetadataRowContainer.MetadataRowContainerRenderer.Rows

			info.Artist = rows.Get("Artist")
			info.Song = rows.Get("Song")
			info.Writers = rows.Get("Writers")
		}
	}

	info.ID = id

	var jsonConfig playerConfig

	// match json in javascript
	if matches := regexpPlayerConfig.FindSubmatch(html); len(matches) > 1 {

		err := json.Unmarshal(matches[1], &jsonConfig)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debug().Msg("Unable to extract json from default url, trying embedded url")

		info, err := getVideoInfoFromEmbedded(id)
		if err != nil {
			return nil, err
		}
		query := url.Values{
			"video_id": []string{id},
			"eurl":     []string{youtubeVideoEURL + id},
		}

		if sts, ok := info["sts"].(float64); ok {
			query.Add("sts", strconv.Itoa(int(sts)))
		}

		body, err := httpGetAndCheckResponseReadBody(youtubeVideoInfoURL + "?" + query.Encode())

		if err != nil {
			return nil, fmt.Errorf("Unable to read video info: %w", err)
		}
		query, err = url.ParseQuery(string(body))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse video info data: %w", err)
		}

		for k, v := range query {
			switch k {
			case "errorcode":
				jsonConfig.Args.Errorcode = v[0]
			case "reason":
				jsonConfig.Args.Reason = v[0]
			case "status":
				jsonConfig.Args.Status = v[0]
			case "player_response":
				jsonConfig.Args.PlayerResponse = v[0]
			case "url_encoded_fmt_stream_map":
				jsonConfig.Args.URLEncodedFmtStreamMap = v[0]
			case "adaptive_fmts":
				jsonConfig.Args.AdaptiveFmts = v[0]
			case "dashmpd":
				jsonConfig.Args.Dashmpd = v[0]
			default:
				// log.Debug().Msgf("unknown query param: %v", k)
			}
		}
	}

	inf := jsonConfig.Args
	if inf.Status == "fail" {
		return nil, fmt.Errorf("Error %s:%s", inf.Errorcode, inf.Reason)
	}

	if inf.PlayerResponse != "" {
		response := &playerResponse{}

		if err := json.Unmarshal([]byte(inf.PlayerResponse), &response); err != nil {
			return nil, fmt.Errorf("Couldn't parse player response: %w", err)
		}

		if response.PlayabilityStatus.Status != "OK" {
			return nil, fmt.Errorf("Unavailable because: %s", response.PlayabilityStatus.Reason)
		}

		if seconds := response.VideoDetails.LengthSeconds; seconds != "" {
			val, err := strconv.Atoi(seconds)
			if err == nil {
				info.Duration = time.Duration(val) * time.Second
			}
		}

		if date, err := time.Parse(youtubeDateFormat, response.Microformat.Renderer.PublishDate); err == nil {
			info.DatePublished = date
		} else {
			log.Debug().Msgf("Unable to parse date published %v", err)
		}

		info.Title = response.VideoDetails.Title
		info.Uploader = response.VideoDetails.Author
	} else {
		log.Debug().Msg("Unable to extract player response JSON")
	}

	info.htmlPlayerFile = jsonConfig.Assets.JS

	var formatStrings []string
	if fmtStreamMap := inf.URLEncodedFmtStreamMap; fmtStreamMap != "" {
		formatStrings = append(formatStrings, strings.Split(fmtStreamMap, ",")...)
	}

	if adaptiveFormats := inf.AdaptiveFmts; adaptiveFormats != "" {
		formatStrings = append(formatStrings, strings.Split(adaptiveFormats, ",")...)
	}

	var formats FormatList
	for _, v := range formatStrings {
		query, err := url.ParseQuery(v)
		if err == nil {
			itag, _ := strconv.Atoi(query.Get("itag"))
			if format, ok := newFormat(itag); ok {
				if strings.HasPrefix(query.Get("conn"), "rtmp") {
					format.meta["rtmp"] = true
				}
				for k, v := range query {
					if len(v) == 1 {
						format.meta[k] = v[0]
					} else {
						format.meta[k] = v
					}
				}
				formats = append(formats, format)
			} else {
				log.Debug().Msgf("No metadata found for itag: %v, skipping...", itag)
			}
		} else {
			log.Debug().Msgf("Unable to format string %v", err)
		}
	}

	if dashManifestURL := inf.Dashmpd; dashManifestURL != "" {
		tokens, err := getSigTokens(info.htmlPlayerFile)
		if err != nil {
			return nil, fmt.Errorf("Unable to extract signature tokens: %w", err)
		}
		regex := regexp.MustCompile("\\/s\\/([a-fA-F0-9\\.]+)")
		regexSub := regexp.MustCompile("([a-fA-F0-9\\.]+)")
		dashManifestURL = regex.ReplaceAllStringFunc(dashManifestURL, func(str string) string {
			return "/signature/" + decipherTokens(tokens, regexSub.FindString(str))
		})
		dashFormats, err := getDashManifest(dashManifestURL)
		if err != nil {
			return nil, fmt.Errorf("Unable to extract dash manifest: %w", err)
		}

		for _, dashFormat := range dashFormats {
			added := false
			for j, format := range formats {
				if dashFormat.Itag == format.Itag {
					formats[j] = dashFormat
					added = true
					break
				}
			}
			if !added {
				formats = append(formats, dashFormat)
			}
		}
	}
	info.Formats = formats

	return info, nil
}

func getVideoInfoFromEmbedded(id string) (map[string]interface{}, error) {
	var jsonConfig map[string]interface{}

	html, err := httpGetAndCheckResponseReadBody(youtubeEmbededBaseURL + id)

	if err != nil {
		return nil, fmt.Errorf("Embeded url request returned %w", err)
	}

	//	re = regexp.MustCompile("\"sts\"\\s*:\\s*(\\d+)")
	re := regexp.MustCompile("yt.setConfig\\({'PLAYER_CONFIG': (.*?)}\\);")

	matches := re.FindSubmatch(html)
	if len(matches) < 2 {
		return nil, fmt.Errorf("Error extracting sts from embedded url response")
	}
	dec := json.NewDecoder(bytes.NewBuffer(matches[1]))
	err = dec.Decode(&jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to extract json from embedded url: %w", err)
	}

	return jsonConfig, nil
}

func getDashManifest(urlString string) (formats []Format, err error) {

	resp, err := httpGetAndCheckResponse(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	var token xml.Token
	for ; err == nil; token, err = dec.Token() {
		if el, ok := token.(xml.StartElement); ok && el.Name.Local == "Representation" {
			var rep representation
			err = dec.DecodeElement(&rep, &el)
			if err != nil {
				break
			}
			if format, ok := newFormat(rep.Itag); ok {
				format.meta["url"] = rep.URL
				if rep.Height != 0 {
					format.Resolution = strconv.Itoa(rep.Height) + "p"
				} else {
					format.Resolution = ""
				}
				formats = append(formats, format)
			} else {
				log.Debug().Msgf("No metadata found for itag: %v, skipping...", rep.Itag)
			}
		}
	}
	if err != io.EOF {
		return nil, err
	}
	return formats, nil
}
