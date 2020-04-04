package ytdl

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestRegexPlayerConfig(t *testing.T) {
	var jsonConfig playerConfig
	html, err := ioutil.ReadFile("fixtures/1.html")
	if err != nil {
		t.Fatal(err)
	}
	// match json in javascript
	if matches := regexpPlayerConfig.FindSubmatch(html); len(matches) > 1 {
		data := append(matches[1], []byte("}")...)
		err := json.Unmarshal(data, &jsonConfig)
		if err != nil {
			t.Error(err)
		}
	}
}
