package main

import (
	"context"
	"log"
	"os"

	ulai "github.com/m0t0k1ch1/go-ulai"
	mastodon "github.com/mattn/go-mastodon"
)

const (
	LogPrefix = "[mastodon-pug] "
)

type Bot struct {
	config Config
	ai     *ulai.Client
	mstdn  *mastodon.Client
	logger *log.Logger
}

func NewBot(config Config) *Bot {
	ai := ulai.NewClient()
	ai.SetKey(config.Ai.Key)

	return &Bot{
		config: config,
		ai:     ai,
		mstdn:  mastodon.NewClient(config.Mastodon),
		logger: log.New(
			os.Stdout,
			LogPrefix,
			log.Ldate|log.Ltime,
		),
	}
}

func (bot *Bot) log(v ...interface{}) {
	bot.logger.Println(v)
}

func (bot *Bot) Run() error {
	evChan, err := bot.mstdn.StreamingUser(context.Background())
	if err != nil {
		return err
	}

	bot.log("start streaming")

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
