package main

import mastodon "github.com/mattn/go-mastodon"

type Config struct {
	Ai       *AiConfig        `json:"ai"`
	Mastodon *mastodon.Config `json:"mastodon"`
}

type AiConfig struct {
	Key string `json:"key"`
}
