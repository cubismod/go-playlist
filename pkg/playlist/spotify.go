package playlist

import (
	"context"
	"time"

	"github.com/zmb3/spotify/v2"
)

const TimeoutTime = 1 * time.Hour

func getItems(ctx context.Context, client *spotify.Client, config SpotifyConfig, playlistID string) []spotify.PlaylistItem {
	ctx, cancel := context.WithTimeout(ctx, TimeoutTime)
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
