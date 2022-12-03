package playlist

import (
	"context"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/go-co-op/gocron"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func setupDB() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		log.WithError(err).Fatal("Unable to load database")
	}
	return db
}

func main() {

	// basic client setup
	ctx := context.Background()
	clientConfig := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := clientConfig.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	log.Info("Spotify connected")

	spotifyConfig := load()

	db := setupDB()

	cacheTask(spotifyConfig, db, client)

	// now we setup the cron loop
	scheduler := gocron.NewScheduler(time.Local)

	scheduler.Every(2).Day().At("03:30").Do(cleanupTask, spotifyConfig, client)
	scheduler.Every(6).Hours().Do(cacheTask, spotifyConfig, db, client)
	scheduler.Cron("0 5 * * 1").Do(scanPlaylistTask, "DiscoverWeekly", spotifyConfig, db, client)
	scheduler.Cron("0 5 * * 1").Do(scanPlaylistTask, "DiscoverWeekly", spotifyConfig, db, client)
	scheduler.Cron("0 5 * * 5").Do(scanPlaylistTask, "ReleaseRadar", spotifyConfig, db, client)
	scheduler.Cron("0 7 1/14 * *").Do(scanPlaylistTask, scanPlaylistTask, "OnRepeat", spotifyConfig, db, client)
	scheduler.Cron("0 9 1/12 * *").Do(scanPlaylistTask, scanPlaylistTask, "RepeatRewind", spotifyConfig, db, client)
	scheduler.Cron("0 11 1/10 * *").Do(scanPlaylistTask, scanPlaylistTask, "RepeatRewind", spotifyConfig, db, client)

	scheduler.StartBlocking()

}
