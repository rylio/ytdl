package ytdl

import "testing"

func TestPlaylist(t *testing.T) {
	id := "PLUl4u3cNGP61o86bXQPx7oCAHDrkt316Y"
	_, err := NewPlaylistFromID(id)
	if err != nil {
		t.Error(err)
	}
}

func TestPlaylistWithMore(t *testing.T) {
	id := "PLfv8hhuTqJGZTQ_gPG2t-pkJQLHbQoX2O"
	playlist, err := NewPlaylistFromID(id)
	if err != nil {
		t.Error(err)
	}
	if len(playlist.VideoIDs) <= 100 {
		t.Error("Did not fetch more videos")
	}
}
