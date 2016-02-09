package ytdl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

const youtubePlaylistURL = youtubeBaseURL + "/playlist"

// Playlist is a list of videos in a youtube playlist
type Playlist struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	VideoIDs    []string `json:"videoIDs"`
}

// NewPlaylistFromID gets a playlist from a playlist id
func NewPlaylistFromID(id string) (*Playlist, error) {
	urlStr := youtubePlaylistURL + "?list=" + id
	res, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	ids, err := extractVideoIDsFromPlaylistHTML(doc, "")
	if err != nil {
		return nil, err
	}
	playlist := &Playlist{
		ID:       id,
		VideoIDs: ids,
	}
	playlist.Title, _ = doc.Find("meta[name=\"title\"]").Attr("content")
	playlist.Description, _ = doc.Find("meta[name=\"description\"]").Attr("content")
	playlist.Author = doc.Find(".pl-header-details > li:first-child").Text()

	return playlist, nil
}

// NewPlaylistFromURLString gets a playlist from a url string for a playlist
func NewPlaylistFromURLString(urlStr string) (*Playlist, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}
	return NewPlaylistFromURL(u)
}

// NewPlaylistFromURL gets a playlist from a url struct
func NewPlaylistFromURL(u *url.URL) (*Playlist, error) {
	id := u.Query().Get("list")
	if id == "" {
		return nil, fmt.Errorf("Invalid playlist url, no id")
	}
	return NewPlaylistFromID(id)
}

func extractVideoIDsFromPlaylistHTML(doc *goquery.Document, more string) ([]string, error) {
	sel := doc.Find("[data-video-id]")
	var videoIDs = make([]string, 0, sel.Length())
	sel.Each(func(_ int, sel *goquery.Selection) {
		if id, ok := sel.Attr("data-video-id"); ok {
			videoIDs = append(videoIDs, id)
		}
	})
	if more == "" {
		more, _ = doc.Find("[data-uix-load-more-href]").First().Attr("data-uix-load-more-href")
	}
	if more != "" {
		res, err := http.Get(youtubeBaseURL + more)
		if err == nil {
			var contents []byte
			contents, err = ioutil.ReadAll(res.Body)
			res.Body.Close()
			jsonData := make(map[string]string)
			err = json.Unmarshal(contents, &jsonData)
			if err == nil {
				var moreVideoIds []string
				var moreHTML *goquery.Document
				// html <tr>s won't be parsed unless wrapped by a table
				contentHTML := "<table>" + jsonData["content_html"] + "</table>"
				moreHTML, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(contentHTML)))
				if err == nil {
					moreNextURL := ""
					if moreNextHTML, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(jsonData["load_more_widget_html"]))); err == nil {
						moreNextURL, _ = moreNextHTML.Find("[data-uix-load-more-href]").Attr("data-uix-load-more-href")
					}
					moreVideoIds, err = extractVideoIDsFromPlaylistHTML(moreHTML, moreNextURL)
					if err == nil {
						videoIDs = append(videoIDs, moreVideoIds...)
					}
				}
			}
		}
		if err != nil {
			log.Debug("Unable to get more playlist items", err.Error())
		}
	}
	return videoIDs, nil
}
