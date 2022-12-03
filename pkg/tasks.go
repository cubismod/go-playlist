package playlist

import (
	"github.com/apex/log"
	"github.com/dgraph-io/badger/v3"
	"github.com/zmb3/spotify/v2"
)

func cleanupTask(spotifyConfig SpotifyConfig, client *spotify.Client) {
	log.Info("Performing AggregatorPlaylist cleanup")
}

func cacheTask(spotifyConfig SpotifyConfig, db *badger.DB, client *spotify.Client) {
	log.Info("Caching AggregatorPlaylist tracks to disk")
}

func scanPlaylistTask(playlist string, spotifyConfig SpotifyConfig, db *badger.DB, client *spotify.Client) {
	log.WithFields(log.Fields{
		"action":   "scanning",
		"playlist": playlist,
	})
}
