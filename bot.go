package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	ulai "github.com/m0t0k1ch1/go-ulai"
	mastodon "github.com/mattn/go-mastodon"
)

const (
	LogPrefix = "mastodon-pug: "
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
			time.AfterFunc(bot.config.Bot.Response.Delay*time.Second, func() {
				notification := evType.Notification
				if notification.Type != "mention" {
					return
				}

				// extract message
				fromMessage, err := ExtractMessage(notification.Status.Content)
				if err != nil {
					bot.log(err)
					return
				}
				bot.log(notification.Account.Acct+":", fromMessage)
				fromMessage = strings.TrimPrefix(fromMessage, "@"+bot.config.Bot.Name)

				// build response message
				toMessage, err := bot.ai.Chat(context.Background(), fromMessage)
				if err != nil {
					bot.log(err)
					return
				}
				toMessage = fmt.Sprintf(
					"@%s %s",
					notification.Status.Account.Acct, toMessage,
				)
				bot.log(bot.config.Bot.Name+":", toMessage)

				// post message
				if _, err := bot.mstdn.PostStatus(context.Background(), &mastodon.Toot{
					Status:      toMessage,
					InReplyToID: notification.Status.ID,
					Visibility:  bot.config.Bot.Response.Visibility,
				}); err != nil {
					bot.log(err)
				}
			})
		}
	}

	bot.log("stop streaming")

	return nil
}
