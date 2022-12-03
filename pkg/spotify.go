package playlist

import (
	"github.com/zmb3/spotify/v2"
)

func getSongs(client *spotify.Client, playlistID string, config SpotifyConfig) ([]spotify.SimpleTrack, error) {
	tracks, err := client.GetPlaylistItems()
}
