package ytdl

import "strings"

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
	URL              string  `json:"url"`
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

type initialData struct {
	Contents struct {
		TwoColumnWatchNextResults struct {
			Results struct {
				Results struct {
					Contents []struct {
						VideoSecondaryInfoRenderer struct {
							Owner struct {
								VideoOwnerRenderer struct {
									Thumbnail struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"thumbnail"`
									Title               Content `json:"title"`
									SubscriberCountText Content `json:"subscriberCountText"`
									TrackingParams      string  `json:"trackingParams"`
								} `json:"videoOwnerRenderer"`
							} `json:"owner"`
							Description          Content `json:"description"`
							MetadataRowContainer struct {
								MetadataRowContainerRenderer struct {
									Rows MetadataRows `json:"rows"`
								} `json:"metadataRowContainerRenderer"`
							} `json:"metadataRowContainer"`
						} `json:"videoSecondaryInfoRenderer,omitempty"`
					} `json:"contents"`
				} `json:"results"`
			} `json:"results"`
		} `json:"twoColumnWatchNextResults"`
	} `json:"contents"`
}

type Content struct {
	SimpleText *string `json:"simpleText,omitempty"`
	Lines      []struct {
		Text string `json:"text,omitempty"`
	} `json:"runs"`
}

func (c *Content) String() string {
	if c.SimpleText != nil {
		return *c.SimpleText
	}

	var sb strings.Builder
	for i := range c.Lines {
		sb.WriteString(c.Lines[i].Text)
	}
	return sb.String()
}

type MetadataRows []struct {
	MetadataRowRenderer struct {
		Title    Content   `json:"title"`
		Contents []Content `json:"contents"`
	} `json:"metadataRowRenderer,omitempty"`
}

func (rows MetadataRows) Get(title string) string {
	for i := range rows {
		row := &rows[i]

		if row.MetadataRowRenderer.Title.String() == title {
			if contents := row.MetadataRowRenderer.Contents; len(contents) > 0 {
				return contents[0].String()
			}
		}
	}

	return ""
}
