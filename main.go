package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	ulai "github.com/m0t0k1ch1/go-ulai"
	mastodon "github.com/mattn/go-mastodon"
)

const (
	DefaultConfigPath = "config.json"
)

type Config struct {
	Bot      *BotConfig       `json:"bot"`
	Ai       *AiConfig        `json:"ai"`
	Mastodon *mastodon.Config `json:"mastodon"`
}

type BotConfig struct {
	User string `json:"user"`
}

type AiConfig struct {
	Key string `json:"key"`
}

type Bot struct {
	config Config
	ai     *ulai.Client
	mstdn  *mastodon.Client
}

func NewBot(config Config) *Bot {
	ai := ulai.NewClient()
	ai.SetKey(config.Ai.Key)

	return &Bot{
		config: config,
		ai:     ai,
		mstdn:  mastodon.NewClient(config.Mastodon),
	}
}

func (bot *Bot) Run() error {
	evChan, err := bot.mstdn.StreamingUser(context.Background())
	if err != nil {
		return err
	}

	for ev := range evChan {
		switch evType := ev.(type) {
		case *mastodon.NotificationEvent:
			notification := evType.Notification
			if notification.Type != "mention" {
				continue
			}

			// TODO
			log.Println(notification.Status.Content)
		}
	}

	return nil
}

func main() {
	var configPath = flag.String("conf", DefaultConfigPath, "path to your config file")
	flag.Parse()

	configFile, err := os.Open(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Fatal(err)
	}

	if err := NewBot(config).Run(); err != nil {
		log.Fatal(err)
	}
}
