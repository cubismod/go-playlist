package playlist

import (
	"context"
	"math/rand"
	"time"

	"github.com/apex/log"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/zmb3/spotify/v2"
)

func ScanAndAdd(playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	log.WithFields(log.Fields{
		"action":   "scanning",
		"playlist": playlistID,
	}).Info("Scan and add")

	items := getItems(client, config, playlistID)
	var trackIDs []spotify.ID

	// now add to aggregator playlist
	for i, item := range items {
		if i%70 == 0 && len(trackIDs) != 0 {
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

// delete a few random items then find duplicates and remove
func CleanupTask(playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	items := getItems(client, config, playlistID)
	var deleteItems []spotify.ID

	for i := 0; i < 5; i++ {
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

	_, err := client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), deleteItems...)
	if err != nil {
		log.WithError(err)
		return
	}

	duplicates := mapset.NewSet[*spotify.TrackToRemove]()

	// find duplicates
	for i, i1 := range items {
		for j, i2 := range items {
			if i != j && i1.Track.Track != nil && i2.Track.Track != nil &&
				i1.Track.Track.Name == i2.Track.Track.Name {

				trackToRemove := *&spotify.TrackToRemove{
					URI: i2.Track.Track.Endpoint,
					Positions: ,
				}
				duplicates.Add()
				log.WithFields(log.Fields{
					"name":   i1.Track.Track.Name,
					"artist": i1.Track.Track.Artists,
					"album":  i1.Track.Track.Album.Name,
				}).Info("Remove duplicate")
			}
		}
	}

	_, err = client.RemoveTracksFromPlaylist(ctx, spotify.ID(playlistID), duplicates.ToSlice()...)
	if err != nil {
		log.WithError(err)
		return
	}
}
