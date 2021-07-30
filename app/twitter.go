package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	postURL = "https://api.twitter.com/1.1/statuses/update.json"

	// API Key/Secret
	consumerKey  = ""
	consumerSecret = ""

	// Access Token/Secret
	accessToken = ""
	accessSecret = ""
)

// TwitterPost wraps the tweet body text to be posted to Twitter.
type TwitterPost struct {
	Status string `json:"status"`  // Tweet body.
}

func doAuth() *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

// PostTweet takes a TwitterPost object and posts to Twitter.
func PostTweet(post string) error {
	// Authenticate and get Twitter client.
	client := doAuth()

	// Send the Tweet
	_, _, err := client.Statuses.Update(post, nil)
	return err
}
