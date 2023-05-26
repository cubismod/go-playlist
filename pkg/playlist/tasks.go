package playlist

import (
	"context"
	"math/rand"
	"time"

	"github.com/apex/log"
	"github.com/zmb3/spotify/v2"
)

func ScanAndAdd(playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	log.WithFields(log.Fields{
		"action":   "scanning",
		"playlist": playlistID,
	})

	items := getItems(client, config, playlistID)
	var trackIDs []spotify.ID

	// now add to aggregator playlist
	for i, item := range items {
		if i%70 == 0 && len(trackIDs) != 0 {
			_, err := client.AddTracksToPlaylist(ctx, spotify.ID(config.Aggregator.ID), trackIDs...)
			trackIDs = nil

			if err != nil {
				log.WithError(err)
			}
		} else {
			trackIDs = append(trackIDs, item.Track.Track.ID)
			log.WithFields(log.Fields{
				"playlistID": playlistID,
				"action":     "add to aggregator",
				"trackName":  item.Track.Track.Name,
				"artists":    item.Track.Track.Artists,
				"album":      item.Track.Track.Album.Name,
			})
		}
	}
}

// delete a few random items then find duplicates and remove
func CleanupTask(playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	items := getItems(client, config, playlistID)
	var deleteItems []spotify.ID

	for i := 0; i < 20; i++ {
		randIndex := rand.Intn(len(items))
		randomItem := items[randIndex]
		deleteItems = append(deleteItems, randomItem.Track.Track.ID)
		log.WithFields(log.Fields{
			"action": "remove from playlist",
			"name":   randomItem.Track.Track.Name,
			"artist": randomItem.Track.Track.Artists,
			"album":  randomItem.Track.Track.Album,
		})
	}

	_, err := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), deleteItems...)
	if err != nil {
		log.WithError(err)
		return
	}
}
