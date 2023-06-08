package playlist

import (
	"context"
	"math/rand"
	"time"

	"github.com/apex/log"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/zmb3/spotify/v2"
)

func ScanAndAdd(ctx context.Context, playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.WithFields(log.Fields{
		"action":   "scanning",
		"playlist": playlistID,
	}).Info("Scan and add")

	items := getItems(ctx, client, config, playlistID)
	var trackIDs []spotify.ID

	// now add to aggregator playlist
	for i, item := range items {
		if i%70 == 0 && len(trackIDs) != 0 {
			ctx, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()

			_, err := client.AddTracksToPlaylist(ctx, spotify.ID(config.Aggregator.ID), trackIDs...)
			trackIDs = nil

			if err != nil {
				log.WithError(err).Error("Unable to add tracks to playlist")
			}
			return
		} else {
			trackIDs = append(trackIDs, item.Track.Track.ID)
			log.WithFields(log.Fields{
				"playlistID": playlistID,
				"trackName":  item.Track.Track.Name,
				"artists":    item.Track.Track.Artists,
				"album":      item.Track.Track.Album.Name,
			}).Info("Adding to aggregator")
		}
	}

	_, err := client.AddTracksToPlaylist(ctx, spotify.ID(config.Aggregator.ID), trackIDs...)

	if err != nil {
		log.WithError(err).Error("Unable to add tracks to playlist")
	}
}

// removes duplicate songs from spotify then re adds them back
// the reason is because you can't delete individual tracks via positions in the web api seemingly
// https://developer.spotify.com/documentation/web-api/reference/remove-tracks-playlist
func removeAndAdd(ctx context.Context, playlistID string, idsToRemove []spotify.ID, client *spotify.Client) {
	ctx, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	_, err := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), idsToRemove...)
	if err != nil {
		log.WithError(err).Error("Error removing tracks!")
	}
	_, err = client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), idsToRemove...)
	if err != nil {
		log.WithError(err).Error("Error adding track back!")
	}
}

// delete a few random items then find duplicates and remove
func CleanupTask(ctx context.Context, playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	items := getItems(ctx, client, config, playlistID)
	var deleteItems []spotify.ID

	if len(items) > 0 {
		for i := 0; i < 25; i++ {
			randIndex := rand.Intn(len(items))
			randomItem := items[randIndex]
			if randomItem.Track.Track != nil {
				deleteItems = append(deleteItems, randomItem.Track.Track.ID)
				log.WithFields(log.Fields{
					"name":   randomItem.Track.Track.Name,
					"artist": randomItem.Track.Track.Artists,
					"album":  randomItem.Track.Track.Album.Name,
				}).Info("Remove from playlist")
			}
		}
	}

	_, err := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), deleteItems...)
	if err != nil {
		log.WithError(err).Error("Error removing tracks!")
		return
	}

	duplicates := mapset.NewSet[string]()

	// find duplicates
	for i, i1 := range items {
		for j, i2 := range items {
			if i != j && i1.Track.Track != nil && i2.Track.Track != nil &&
				i1.Track.Track.Name == i2.Track.Track.Name {
				duplicates.Add(i1.Track.Track.ID.String())

				log.WithFields(log.Fields{
					"name":   i1.Track.Track.Name,
					"artist": i1.Track.Track.Artists,
					"album":  i1.Track.Track.Album.Name,
					"pos_1":  i,
					"pos_2":  j,
				}).Info("Remove duplicate")
			}
		}
	}

	var idsToRemove []spotify.ID

	for k, d := range duplicates.ToSlice() {
		idsToRemove = append(idsToRemove, spotify.ID(d))

		if k != 0 && k%70 == 0 {
			removeAndAdd(ctx, playlistID, idsToRemove, client)
			idsToRemove = nil
		}
	}

	if idsToRemove != nil {
		removeAndAdd(ctx, playlistID, idsToRemove, client)
	}
}
