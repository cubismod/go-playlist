package playlist

import (
	"context"
	"time"

	"github.com/zmb3/spotify/v2"
)

func getItems(client *spotify.Client, config SpotifyConfig, playlistID string) []spotify.PlaylistItem {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var tracks []spotify.PlaylistItem
	trackPage, err := client.GetPlaylistItems(ctx, spotify.ID(playlistID))
	for {
		if err == nil {
			tracks = append(tracks, trackPage.Items...)
			err = client.NextPage(ctx, trackPage)
		} else {
			break
		}
	}
	return tracks
}
