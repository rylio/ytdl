package ytdl

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataRows(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	jsondata, err := ioutil.ReadFile("fixtures/metadata_rows.json")
	require.NoError(err)

	info := VideoInfo{}
	require.NoError(info.addMetadata(jsondata))

	assert.Equal("Justin Timberlake", info.Artist)
	assert.Equal("Tunnel Vision", info.Song)
}
