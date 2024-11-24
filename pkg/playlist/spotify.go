package playlist

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/zmb3/spotify/v2"
)

const TimeoutTime = 30 * time.Second

func getItems(ctx context.Context, client *spotify.Client, config SpotifyConfig, playlistId string) []spotify.PlaylistItem {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var tracks []spotify.PlaylistItem
	trackPage, err := client.GetPlaylistItems(ctx, spotify.ID(playlistId))
	for {
		if err == nil {
			tracks = append(tracks, trackPage.Items...)
			err = client.NextPage(ctx, trackPage)
		}
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.WithError(err).Error("Error fetching playlist items")
			return tracks
		}
	}
	log.WithFields(log.Fields{
		"playlistId": playlistId,
		"tracks":     len(tracks),
	}).Info("Loaded songs from playlist")
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

func clearPlaylist(ctx context.Context, playlistID string, config SpotifyConfig, client *spotify.Client) error {
	ctx, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	log.WithFields(log.Fields{
		"action":   "clear playlist",
		"playlist": playlistID,
	}).Info("Clearing")

	playlistItems := getItems(ctx, client, config, playlistID)
	trackIDs := getPlaylistIDs(playlistItems)

	err := removeFromPlaylist(ctx, client, playlistID, trackIDs)

	return err

}

// retrieves the playlist IDs from []spotify.PlaylistItem
func getPlaylistIDs(items []spotify.PlaylistItem) []spotify.ID {
	var ids []spotify.ID
	for _, item := range items {
		ids = append(ids, item.Track.Track.ID)
	}
	return ids
}
