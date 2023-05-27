package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/cubismod/go-playlist/pkg/playlist"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/zmb3/spotify/v2"
)

func serve(spotifyConfig playlist.SpotifyConfig, client *spotify.Client) {
	// now we setup the cron loop
	scheduler := gocron.NewScheduler(time.Local)

	for _, configPlaylist := range spotifyConfig.Playlists {
		_, err := scheduler.Cron(configPlaylist.ScanCron).Do(playlist.ScanAndAdd, configPlaylist.ID, spotifyConfig, client)
		if err != nil {
			log.WithError(err).Error("Could not schedule job")
		} else {
			log.WithFields(log.Fields{
				"action":   "schedule_job",
				"job_name": "scan_job",
				"playlist": configPlaylist.Name,
				"cron":     configPlaylist.ScanCron,
			}).Info("Scheduled job")
		}
	}

	_, err := scheduler.Cron(spotifyConfig.Aggregator.CleanupCron).Do(playlist.CleanupTask, spotifyConfig.Aggregator.ID, spotifyConfig, client)
	if err != nil {
		log.WithError(err).Fatal("Could not schedule cleanup job")
	}
	scheduler.StartBlocking()
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.WithError(err).Fatal("Failed to load env variables")
	}

	log.Info(os.Getenv("SPOTIFY_ID"))

	// basic client setup
	client := playlist.RunAuthServer()

	log.Info("Spotify connected")

	spotifyConfig := playlist.Load()

	app := &cli.App{
		Name:  "go-playlist",
		Usage: "Automate Spotify with Go",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "run a persistent cron server",
				Action: func(cCtx *cli.Context) error {
					serve(spotifyConfig, client)
					return nil
				},
			},
			{
				Name:  "scan",
				Usage: "scan and add an individual playlist",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						return cli.Exit("You must specify a playlist ID", -1)
					}

					playlistID := cCtx.Args().Get(0)
					playlist.ScanAndAdd(playlistID, spotifyConfig, client)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("unable to start app")
	}
}
