package main

import (
	"github.com/dghubble/go-twitter/twitter"

	"github.com/dghubble/oauth1"
)

type TwitterClient struct {
	cfg    TwitterConfig
	client *twitter.Client
}

// TwitterPost wraps the tweet body text to be posted to Twitter.
type TwitterPost struct {
	Status string `json:"status"` // Tweet body.
}

func (c *TwitterClient) doAuth() *twitter.Client {
	config := oauth1.NewConfig(c.cfg.ConsumerKey, c.cfg.ConsumerSecret)
	token := oauth1.NewToken(c.cfg.AccessToken, c.cfg.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

// PostTweet takes a TwitterPost object and posts to Twitter.
func (c *TwitterClient) PostTweet(post string) error {
	_, _, err := c.client.Statuses.Update(post, nil)
	return err
}

func NewTwitterClient(config TwitterConfig) *TwitterClient {
	client := &TwitterClient{cfg: config}
	c := client.doAuth()
	client.client = c
	return client
}
