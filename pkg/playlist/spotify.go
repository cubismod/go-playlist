package playlist

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/zmb3/spotify/v2"
)

const TimeoutTime = 30 * time.Second

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

func addToPlaylist(ctx context.Context, client *spotify.Client, playlistID string, ids []spotify.ID) {
	ctx, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	_, err := client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), ids...)

	if err != nil {
		log.WithError(err).Error("Unable to add tracks to aggregator playlist")
	}
}

func removeFromPlaylist(ctx context.Context, client *spotify.Client, playlistID string, ids []spotify.ID) error {
	ctx, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	_, err := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), ids...)

	return err
}
