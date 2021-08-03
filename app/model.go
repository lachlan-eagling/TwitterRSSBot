package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type TwitterConfig struct {
	PostURL        string `yaml:"post_url"`
	ConsumerKey    string `yaml:"consumer_key"`
	ConsumerSecret string `yaml:"consumer_secret"`
	AccessToken    string `yaml:"access_token"`
	AccessSecret   string `yaml:"access_secret"`
}

type Config struct {
	Twitter        TwitterConfig `yaml:"twitter"`
	Sources        []Source      `yaml:"sources"`
	UpdateSeenOnly bool
	TestRun bool
	SeenDataPath string
}

type Source struct {
	PublicationName string   `yaml:"publication_name"`
	AuthorTwitter   string   `yaml:"author_twitter"`
	URL             string   `yaml:"url"`
	HashTags        []string `yaml:"hash_tags"`
}

type PostData struct {
	AuthorTwitterHandle string
	URL                 string
	HashTags            string
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(f, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
