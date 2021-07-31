package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
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

func ReadSeenFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// CheckAndAddToSeen checks to see if we have recorded seeing this URL
// in a feed before, if it is a new one it is appended to the slice
// of seen urls and returns false. Otherwise if it is a seen one returns true.
func CheckSeen(url string, seen []string) bool {
	for _, s := range seen {
		if s == url {
			return true
		}
	}
	return false
}

// WriteSeen overwrites the seen file with latest list of seen urls.
func WriteSeen(lines []string, path string) error {
	log.Infof("Writing seen urls to %s", path)
	file, err := os.Create(path)
	if err != nil {
		log.Error(err)
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Error(err)
		return err
	}

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func main() {
	cfg, err := LoadConfig("config-prod.yml")
	if err != nil {
		log.Fatalf("Error loading config. Shutting Down. (%s)", err.Error())
	}

	seen, err := ReadSeenFile("seen.txt")
	if err != nil {
		log.Fatalf("Error loading seen urls. Shutting Down. (%s)", err.Error())
	}
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
			if CheckSeen(item.Link, seen) {
				// This is a known URL that has already been posted. Take no further action.
				log.Infof("Already posted %s\n", item.Link)
				continue
			}
			// New URL, go ahead and post it.
			seen = append(seen, item.Link)
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

	WriteSeen(seen, "seen.txt")

	client := NewTwitterClient(cfg.Twitter)

	for _, tweet := range pendingTweets {
		log.Infof("Posting tweet (%s)")
		if err := client.PostTweet(tweet); err != nil {
			log.Errorf("Error posting tweet (%s). %s", tweet, err.Error())
		}
	}

}
