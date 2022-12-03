package playlist

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type SpotifyConfig struct {
	AggregatorPlaylist string
	DiscoverWeekly     string
	ReleaseRadar       string
	OnRepeat           string
	RepeatRewind       string
	TimeCapsule        string
}

func load() SpotifyConfig {
	cfile, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Fatalf("config.yaml file missing\n%v", err)
	}

	config := SpotifyConfig{}

	err = yaml.Unmarshal(cfile, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
