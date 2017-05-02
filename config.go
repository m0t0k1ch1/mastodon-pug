package main

import (
	"time"

	mastodon "github.com/mattn/go-mastodon"
)

type Config struct {
	Bot      *BotConfig       `json:"bot"`
	Ai       *AiConfig        `json:"ai"`
	Mastodon *mastodon.Config `json:"mastodon"`
}

type BotConfig struct {
	Name     string             `json:"name"`
	Response *BotResponseConfig `json:"response"`
}

type BotResponseConfig struct {
	Visibility string        `json:"visibility"`
	Delay      time.Duration `json:"delay"` // sec
}

type AiConfig struct {
	Key string `json:"key"`
}
