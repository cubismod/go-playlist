package playlist

import (
	"context"
	"math/rand"

	"github.com/apex/log"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/zmb3/spotify/v2"
)

// func isDuplicate(item spotify.PlaylistItem, playlistContents []spotify.PlaylistItem) bool {
// 	for _, pi := range playlistContents {
// 		if item.Track.Track != nil && pi.Track.Track != nil &&
// 			item.Track.Track.Name == pi.Track.Track.Name {
// 			return true
// 		}
// 	}
// 	return false
// }

func titlesToSet(playlistItems []spotify.PlaylistItem) mapset.Set[string] {
	set := mapset.NewSet[string]()
	for _, item := range playlistItems {
		if item.Track.Track != nil {
			set.Add(item.Track.Track.Name)
		}
	}
	return set
}

func ScanAndAdd(playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.WithFields(log.Fields{
		"action":   "scanning",
		"playlist": playlistID,
	}).Info("Scan and add")

	addPlaylistItems := getItems(ctx, client, config, playlistID)
	aggregatorItems := getItems(ctx, client, config, config.Aggregator.ID)
	var trackIDs []spotify.ID

	aggregatorTitleSet := titlesToSet(aggregatorItems)

	// now add to aggregator playlist
	for _, item := range addPlaylistItems {
		if len(trackIDs) >= 90 {
			addToPlaylist(ctx, client, config.Aggregator.ID, trackIDs)
			trackIDs = nil
		} else {
			title := []string{item.Track.Track.Name}
			if !aggregatorTitleSet.Contains(title...) {
				trackIDs = append(trackIDs, item.Track.Track.ID)
				log.WithFields(log.Fields{
					"playlistID": playlistID,
					"trackName":  item.Track.Track.Name,
					"artists":    item.Track.Track.Artists,
					"album":      item.Track.Track.Album.Name,
				}).Info("Adding to aggregator")
			}
		}
	}

	if len(trackIDs) != 0 {
		addToPlaylist(ctx, client, config.Aggregator.ID, trackIDs)
	}

	log.WithFields(log.Fields{
		"action":   "finished",
		"playlist": playlistID,
	}).Info("Scan and add")
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

// keep the aggregator playlist under 2500 tracks and remove duplicate songs
func CleanupTask(playlistID string, config SpotifyConfig, client *spotify.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.WithFields(log.Fields{
		"action":   "begin",
		"playlist": playlistID,
	}).Info("Cleanup")

	items := getItems(ctx, client, config, playlistID)
	var deleteItems []spotify.ID

	toDelete := len(items) - 2500

	if toDelete > 0 {
		for i := 0; i < toDelete; i++ {
			if len(deleteItems) >= 90 {
				err := removeFromPlaylist(ctx, client, playlistID, deleteItems)
				if err != nil {
					log.WithError(err).Error("Error removing tracks!")
					return
				}
				deleteItems = nil
			} else {
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

	log.WithFields(log.Fields{
		"action":   "finished",
		"playlist": playlistID,
	}).Info("Cleanup")
}
