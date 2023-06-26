package playlist

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type PlaylistConfig struct {
	Name     string `yaml:"name"`
	ID       string `yaml:"id"`
	ScanCron string `yaml:"scanCron"`
}

type AggregatorConfig struct {
	Name        string `yaml:"name"`
	ID          string `yaml:"id"`
	CleanupCron string `yaml:"cleanupCron"`
}

type SpotifyConfig struct {
	Aggregator AggregatorConfig `yaml:"aggregator"`
	Playlists  []PlaylistConfig `yaml:"playlists"`
}

func LoadConfig() SpotifyConfig {
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
