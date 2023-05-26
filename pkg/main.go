package main

import (
	"context"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/cubismod/go-playlist/pkg/playlist"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.WithError(err)
	}

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
	client := spotify.New(httpClient, spotify.WithRetry(true))

	log.Info("Spotify connected")

	spotifyConfig := playlist.Load()
	// now we setup the cron loop
	scheduler := gocron.NewScheduler(time.Local)

	for _, configPlaylist := range spotifyConfig.Playlists {
		_, err := scheduler.Cron(configPlaylist.ScanCron).Do(playlist.ScanAndAdd, configPlaylist.ID, spotifyConfig, client)
		if err != nil {
			log.WithError(err)
		}
	}

	_, err = scheduler.Cron(spotifyConfig.Aggregator.CleanupCron).Do(playlist.CleanupTask)
	if err != nil {
		log.WithError(err)
	}
	scheduler.StartBlocking()

}
