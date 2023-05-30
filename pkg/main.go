package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/cubismod/go-playlist/pkg/playlist"
	"github.com/go-co-op/gocron"
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
	// basic client setup
	client, err := playlist.RunAuthServer()

	if err != nil {
		log.WithError(err).Fatal("unable to login to spotify")
	}

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
			{
				Name:  "clean",
				Usage: "clean the aggregator playlist",
				Action: func(cCtx *cli.Context) error {
					playlist.CleanupTask(spotifyConfig.Aggregator.ID, spotifyConfig, client)
					return nil
				},
			},
			{
				Name:  "batch",
				Usage: "run all commands in a batch",
				Action: func(cCtx *cli.Context) error {
					for _, spotifyPlaylist := range spotifyConfig.Playlists {
						playlist.ScanAndAdd(spotifyPlaylist.ID, spotifyConfig, client)
					}
					playlist.CleanupTask(spotifyConfig.Aggregator.ID, spotifyConfig, client)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("unable to start app")
	}
}
