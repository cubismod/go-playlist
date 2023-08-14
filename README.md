# go-playlist
Spotify playlist automation with go

Designed to run in a Docker container.

Go playlist works around a central "aggregator playlist" which is used as a large source to multiple different playlists that you can specify.
The application is run as a persistent service and runs cronjobs within go to perform playlist automation tasks.
These include:
- scanning and adding songs to aggregator
    - done on a cronjob
    - duplicates are not added
- cleanup of the aggregator playlist
    - random removal of songs from the playlist to get the size of the aggregator to around 2500 songs
    - duplicates removal (not really working too great heh)

Configuration of go-playlist is done with a YAML file, example:

```yaml
---
aggregator:
  name: Untamed Vibes
  id: 5W9dfQRHyfjUa3W3xTRM9P
  cleanupCron: 00 16 * * *
playlists:
  - name: Release Radar
    id: 37i9dQZEVXbp8mKzlhCRsm
    scanCron: 15 4 * * 5
  - name: Discover Weekly
    id: 37i9dQZEVXcSC5XajOgcu6
    scanCron: 00 5 * * 1
  - name: On Repeat
    id: 37i9dQZF1EphySG32fl9Th
    scanCron: 00 10 1 * *
  - name: Daily Mix 1
    id: 37i9dQZF1E39QMjkhWNgVi
    scanCron: 14 22 * * 2
  - name: Daily Mix 2
    id: 37i9dQZF1E39c63pUlbgc8
    scanCron: 7 11 * * 4
  - name: Daily Mix 5
    id: 37i9dQZF1E36a2Q3RIDXKr
    scanCron: 15 14 20 * *
  - name: Gay Rats Blend
    id: 37i9dQZF1EJsXYbdr4o8jX
    scanCron: 36 20 * * 6
  - name: Homisexual
    id: 37i9dQZF1EJNcLQrIKKFUG
    scanCron: 21 12 * * 3
  - name: Daily Mix 6
    id: 37i9dQZF1E37nLuFiSy1SM
    scanCron: 0 0 1,15 * 3
  - name: Daily Mix 4
    id: 37i9dQZF1E34VBygDvZz0a
    scanCron: 5 4 3 * 0
  - name: Daily Mix 3
    id: 37i9dQZF1E387X5zO3whNF
    scanCron: 45 20 14 * *
  - name: Lorem
    id: 37i9dQZF1DXdwmD5Q7Gxah
    scanCron: 0 2 17 * *
  - name: LGBTQ Discord
    id: 37i9dQZF1EJz7mWdfQBeHl
    scanCron: 46 7 * * 1
  - name: Pollen
    id: 37i9dQZF1DWWBHeXOYZf74
    scanCron: 22 11 19 * *
```
Cron syntax is used to schedule when actions happen.

# Environment Variables
```ini
# retrieve these from the spotify developer portal
SPOTIFY_ID=<id>
SPOTIFY_SECRET=<secret>
 # used for local webserver to get auth details
GO_PLAY_HOSTNAME=http://localhost
GO_PLAY_LISTEN_ADDR=0.0.0.0
# or 127.0.0.1 for localhost
GO_PLAY_PORT=8080
# used to send a notification to your phone when login is requested
NTFY_URL=https://ntfy.com/my_topic_here
NTFY_PW=<secure_password>
```
