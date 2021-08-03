package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	var updateSeen bool
	var testRun bool
	var configPath string
	var seenPath string
	flag.BoolVar(&updateSeen, "u", false, "Parses feeds and updates seen.txt with all posts available up to current time without posting any Tweets.")
	flag.BoolVar(&testRun, "test", false, "Sends a test tweet and immediatley exits the application.")
	flag.StringVar(&configPath, "config", "config.yml", "Path to config.yml file.")
	flag.StringVar(&seenPath, "seen", "seen.txt", "Path to text file containing seen URLs.")
	flag.Parse()

	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config. Shutting Down. (%s)", err.Error())
	}
	cfg.UpdateSeenOnly = updateSeen
	cfg.TestRun = testRun
	cfg.SeenDataPath = seenPath

	bot := NewBot(cfg)
	if err := bot.Run(); err != nil {
		log.Error(err)
	}
}
