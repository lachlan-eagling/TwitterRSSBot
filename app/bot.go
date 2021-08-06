package main

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"bytes"
	"html/template"

	"github.com/mmcdole/gofeed"
)

type Bot struct {
	config        *Config
	store         *Store
	twitterClient *TwitterClient
	parser        *gofeed.Parser
}

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

func (b *Bot) Run() error {
	if b.config.TestRun {
		if err := b.twitterClient.PostTweet("Testing 1,2,3"); err != nil {
			log.Debugf("Error posting test tweet. (%s)", err.Error())
			return err
		}
		return nil
	}
	defer b.store.Close()

	pendingTweets := make([]string, 0)
	for _, source := range b.config.Sources {
		feed, err := b.parser.ParseURL(source.URL)
		if err != nil {
			return err
		}
		hashTags := FlattenHashTags(source.HashTags)
		log.Infof("Parsing (%d) posts for %s", len(feed.Items), feed.Title)

		for _, item := range feed.Items {
			if b.store.Exists(item.Link) {
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

	if b.config.UpdateSeenOnly {
		log.Info("Seen sources updated, exiting application.")
		return nil
	}

	if len(pendingTweets) > 0 {
		log.Infof("Sending %d tweets to Twitter API", len(pendingTweets))
	} else {
		log.Info("No tweets to send. Exiting.")
		return nil
	}

	for _, tweet := range pendingTweets {
		log.Infof("Posting tweet (%s)", tweet)
		if err := b.twitterClient.PostTweet(tweet); err != nil {
			log.Debugf("Error posting tweet (%s). %s", tweet, err.Error())
			return err
		}
	}
	log.Info("Done sending tweets")
	return nil
}

func NewBot(config *Config) *Bot {
	client := NewTwitterClient(config.Twitter)
	parser := gofeed.NewParser()
	store, err := NewStore(config.SeenDataPath)
	if err != nil {
		log.Fatalf("Error setting up seen store. Shutting Down. (%s)", err.Error())
	}

	bot := &Bot{
		config:        config,
		store:         store,
		twitterClient: client,
		parser:        parser,
	}

	return bot
}
