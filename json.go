package ytdl

import (
	"github.com/antchfx/jsonquery"
)

type playerConfig struct {
	Assets struct {
		JS string `json:"js"`
	} `json:"assets"`
	Args struct {
		Status                 string `json:"status"`
		Errorcode              string `json:"errorcode"`
		Reason                 string `json:"reason"`
		PlayerResponse         string `json:"player_response"`
		URLEncodedFmtStreamMap string `json:"url_encoded_fmt_stream_map"`
		AdaptiveFmts           string `json:"adaptive_fmts"`
		Dashmpd                string `json:"dashmpd"`
	} `json:"args"`
}

type formatInfo struct {
	Itag             int     `json:"itag"`
	MimeType         string  `json:"mimeType"`
	Bitrate          int     `json:"bitrate"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	LastModified     string  `json:"lastModified"`
	ContentLength    string  `json:"contentLength"`
	Quality          string  `json:"quality"`
	QualityLabel     string  `json:"qualityLabel"`
	ProjectionType   string  `json:"projectionType"`
	AverageBitrate   int     `json:"averageBitrate"`
	AudioQuality     string  `json:"audioQuality"`
	ApproxDurationMs string  `json:"approxDurationMs"`
	AudioSampleRate  string  `json:"audioSampleRate"`
	AudioChannels    int     `json:"audioChannels"`
	Cipher           *string `json:"cipher"`
	SignatureCipher  *string `json:"signatureCipher"`
	URL              string  `json:"url"`
	Index            *Range  `json:"indexRange,omitempty"`
	Init             *Range  `json:"initRange,omitempty"`
}

type playerResponse struct {
	PlayabilityStatus struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	} `json:"playabilityStatus"`

	StreamingData struct {
		ExpiresInSeconds string       `json:"expiresInSeconds"`
		Formats          []formatInfo `json:"formats"`
		AdaptiveFormats  []formatInfo `json:"adaptiveFormats"`
		DashManifestUrl  string       `json:"dashManifestUrl"`
		HlsManifestUrl   string       `json:"hlsManifestUrl"`
	} `json:"streamingData"`

	VideoDetails struct {
		Title         string   `json:"title"`
		Author        string   `json:"author"`
		LengthSeconds string   `json:"lengthSeconds"`
		Keywords      []string `json:"keywords"`
		ViewCount     string   `json:"viewCount"`
	} `json:"videoDetails"`

	Microformat struct {
		Renderer struct {
			ViewCount   string `json:"viewCount"`
			PublishDate string `json:"publishDate"`
			UploadDate  string `json:"uploadDate"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
}

type representation struct {
	Itag   int    `xml:"id,attr"`
	Height int    `xml:"height,attr"`
	URL    string `xml:"BaseURL"`
}

func getMetaDataRow(row *jsonquery.Node) (string, string) {
	title, _ := jsonquery.Query(row, "title")
	text, _ := jsonquery.Query(row, "contents//simpleText")

	if text == nil {
		text, _ = jsonquery.Query(row, "contents//text")
	}

	if title == nil || text == nil {
		return "", ""
	}

	return title.InnerText(), text.InnerText()
}
