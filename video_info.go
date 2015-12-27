package ytdl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

const youtubeBaseURL = "https://www.youtube.com/watch"
const youtubeEmbededBaseURL = "https://www.youtube.com/embed/"
const youtubeVideoEURL = "https://youtube.googleapis.com/v/"
const youtubeVideoInfoURL = "https://www.youtube.com/get_video_info"
const youtubeDateFormat = "2006-01-02"

// VideoInfo contains the info a youtube video
type VideoInfo struct {
	// The video ID
	ID string `json:"id"`
	// The video title
	Title string `json:"title"`
	// The video description
	Description string `json:"description"`
	// The date the video was published
	DatePublished time.Time `json:"datePublished"`
	// Formats the video is available in
	Formats []Format `json:"formats"`
	// List of keywords associated with the video
	Keywords []string `json:"keywords"`
	// Author of the video
	Author string `json:"author"`
	// Duration of the video
	Duration time.Duration

	htmlPlayerFile string
}

// GetVideoInfo fetches info from a url string, url object, or a url string
func GetVideoInfo(value interface{}) (*VideoInfo, error) {
	switch value.(type) {
	case *url.URL:
		return GetVideoInfoFromURL(value.(*url.URL))
	case string:
		u, err := url.ParseRequestURI(value.(string))
		if err != nil {
			return GetVideoInfoFromID(value.(string))
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

// GetVideoInfoFromID fetches video info from a youtube video id
func GetVideoInfoFromID(id string) (*VideoInfo, error) {
	u, _ := url.ParseRequestURI(youtubeBaseURL)
	values := u.Query()
	values.Set("v", id)
	u.RawQuery = values.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
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

func getVideoInfoFromHTML(id string, html []byte) (*VideoInfo, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	info := &VideoInfo{}

	// extract description and title
	info.Description = strings.TrimSpace(doc.Find("#eow-description").Text())
	info.Title = strings.TrimSpace(doc.Find("#eow-title").Text())
	info.ID = id
	dateStr, ok := doc.Find("meta[itemprop=\"datePublished\"]").Attr("content")
	if !ok {
		log.Debug("Unable to extract date published")
	} else {
		date, err := time.Parse(youtubeDateFormat, dateStr)
		if err == nil {
			info.DatePublished = date
		} else {
			log.Debug("Unable to parse date published", err.Error())
		}
	}

	// match json in javascript
	re := regexp.MustCompile("ytplayer.config = (.*?);ytplayer.load")
	matches := re.FindSubmatch(html)
	var jsonConfig map[string]interface{}
	if len(matches) > 1 {
		err = json.Unmarshal(matches[1], &jsonConfig)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debug("Unable to extract json from default url, trying embedded url")
		var resp *http.Response
		resp, err = http.Get(youtubeEmbededBaseURL + id)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Embeded url request returned status code %d	", resp.StatusCode)
		}
		html, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		//	re = regexp.MustCompile("\"sts\"\\s*:\\s*(\\d+)")
		re = regexp.MustCompile("yt.setConfig\\('PLAYER_CONFIG', (.*?)\\);</script>")

		matches := re.FindSubmatch(html)
		if len(matches) < 2 {
			return nil, fmt.Errorf("Error extracting sts from embedded url response")
		}
		dec := json.NewDecoder(bytes.NewBuffer(matches[1]))
		err = dec.Decode(&jsonConfig)
		if err != nil {
			return nil, fmt.Errorf("Unable to extract json from embedded url: %s", err.Error())
		}
		query := url.Values{
			"sts":      []string{strconv.Itoa(int(jsonConfig["sts"].(float64)))},
			"video_id": []string{id},
			"eurl":     []string{youtubeVideoEURL + id},
		}

		resp, err = http.Get(youtubeVideoInfoURL + "?" + query.Encode())
		if err != nil {
			return nil, fmt.Errorf("Error fetching video info: %s", err.Error)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Video info response invalid status code")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Unable to read video info response body: %s", err.Error())
		}
		query, err = url.ParseQuery(string(body))
		if err != nil {
			return nil, fmt.Errorf("Unable to parse video info data: %s", err.Error())
		}
		args := make(map[string]interface{})
		for k, v := range query {
			if len(v) > 0 {
				args[k] = v[0]
			}
		}
		jsonConfig["args"] = args
	}

	inf := jsonConfig["args"].(map[string]interface{})
	if status, ok := inf["status"].(string); ok && status == "fail" {
		return nil, fmt.Errorf("Error %d:%s", inf["errorcode"], inf["reason"])
	}
	if a, ok := inf["author"].(string); ok {
		info.Author = a
	} else {
		log.Debug("Unable to extract author")
	}

	if length, ok := inf["length_seconds"].(string); ok {
		if duration, err := strconv.ParseInt(length, 10, 64); err == nil {
			info.Duration = time.Second * time.Duration(duration)
		} else {
			log.Debug("Unable to parse duration string: ", length)
		}
	} else {
		log.Debug("Unable to extract duration")
	}

	// For the future maybe
	parseKey := func(key string) []string {
		val, ok := inf[key].(string)
		if !ok {
			return nil
		}
		vals := []string{}
		split := strings.Split(val, ",")
		for _, v := range split {
			if v != "" {
				vals = append(vals, v)
			}
		}
		return vals
	}
	info.Keywords = parseKey("keywords")

	/*
		fmtList := parseKey("fmt_list")
		fexp := parseKey("fexp")
		watermark := parseKey("watermark")

		if len(fmtList) != 0 {
			vals := []string{}
			for _, v := range fmtList {
				vals = append(vals, strings.Split(v, "/")...)
		} else {
			info["fmt_list"] = []string{}
		}

		videoVerticals := []string{}
		if videoVertsStr, ok := inf["video_verticals"].(string); ok {
			videoVertsStr = string([]byte(videoVertsStr)[1 : len(videoVertsStr)-2])
			videoVertsSplit := strings.Split(videoVertsStr, ", ")
			for _, v := range videoVertsSplit {
				if v != "" {
					videoVerticals = append(videoVerticals, v)
				}
			}
		}
	*/
	var formatStrings []string
	if fmtStreamMap, ok := inf["url_encoded_fmt_stream_map"].(string); ok {
		formatStrings = append(formatStrings, strings.Split(fmtStreamMap, ",")...)
	}

	if adaptiveFormats, ok := inf["adaptive_fmts"].(string); ok {
		formatStrings = append(formatStrings, strings.Split(adaptiveFormats, ",")...)
	}
	var formats []Format
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
				log.Debug("No metadata found for itag: ", itag, ", skipping...")
			}
		} else {
			log.Debug("Unable to format string", err.Error())
		}
	}
	info.Formats = formats
	info.htmlPlayerFile = jsonConfig["assets"].(map[string]interface{})["js"].(string)
	return info, nil
}
