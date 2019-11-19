package ytdl

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataRows(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	file, err := os.Open("fixtures/metadata_rows.json")
	require.NoError(err)
	defer file.Close()

	data := initialData{}
	require.NoError(json.NewDecoder(file).Decode(&data))

	contents := data.Contents.TwoColumnWatchNextResults.Results.Results.Contents
	require.Len(contents, 3)

	rows := contents[1].VideoSecondaryInfoRenderer.MetadataRowContainer.MetadataRowContainerRenderer.Rows
	require.Len(rows, 9)

	assert.Equal("Notice", rows[0].MetadataRowRenderer.Title.String())
	assert.Equal("Justin Timberlake", rows.Get("Artist"))
	assert.Equal("Tunnel Vision", rows.Get("Song"))
}
