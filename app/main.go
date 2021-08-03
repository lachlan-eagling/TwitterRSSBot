package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mmcdole/gofeed"
)

func FlattenHashTags(h []string) string {
	tags := make([]string, len(h))
	for i, tag := range h {
		tags[i] = fmt.Sprintf("#%s", tag)
	}
	return strings.Join(tags[:], " ")
}

func FormatAuthor(a string) string {
	if a == "" {
		return ""
	}
	return fmt.Sprintf("@%s", a)
}

func main() {
	var updateSeen bool
	var test bool
	var configPath string
	var seenPath string
	flag.BoolVar(&updateSeen, "u", false, "Parses feeds and updates seen.txt with all posts available up to current time without posting any Tweets.")
	flag.BoolVar(&test, "test", false, "Sends a test tweet and immediatley exits the application.")
	flag.StringVar(&configPath, "config", "config.yml", "Path to config.yml file.")
	flag.StringVar(&seenPath, "seen", "seen.txt", "Path to text file containing seen URLs.")
	flag.Parse()

	log.Infof("Config Path: %s", configPath)
	log.Infof("Seen Path: %s", seenPath)

	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config. Shutting Down. (%s)", err.Error())
	}
	client := NewTwitterClient(cfg.Twitter)
	if test {
		log.Info("Sending test tweet")
		if err := client.PostTweet("Testing 1,2,3!"); err != nil {
			log.Errorf("Error posting test tweet. %s", err.Error())
		}
		return
	}

	store, err := NewStore(seenPath)
	if err != nil {
		log.Fatalf("Error setting up seen store. Shutting Down. (%s)", err.Error())
	}
	defer store.Close()

	parser := gofeed.NewParser()
	pendingTweets := make([]string, 0)

	for _, source := range cfg.Sources {
		feed, err := parser.ParseURL(source.URL)
		if err != nil {
			log.Error(err)
			return
		}
		hashTags := FlattenHashTags(source.HashTags)
		log.Infof("Parsing (%d) posts for %s", len(feed.Items), feed.Title)

		for _, item := range feed.Items {
			if store.Exists(item.Link) {
				log.Infof("Already posted %s\n", item.Link)
				continue
			}

			post := PostData{
				AuthorTwitterHandle: FormatAuthor(source.AuthorTwitter),
				URL:                 item.Link,
				HashTags:            hashTags,
			}
			var out bytes.Buffer
			t := template.Must(template.ParseFiles("tweet.tmpl"))
			t.Execute(&out, post)

			tweet := out.String()
			pendingTweets = append(pendingTweets, tweet)
		}

	}

	if updateSeen {
		log.Info("Seen sources updated, exiting application.")
		return
	}

	for _, tweet := range pendingTweets {
		log.Infof("Posting tweet (%s)", tweet)
		if err := client.PostTweet(tweet); err != nil {
			log.Errorf("Error posting tweet (%s). %s", tweet, err.Error())
		}
	}
}
