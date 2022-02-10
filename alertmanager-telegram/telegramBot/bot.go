package telegramBot

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
	tele "gopkg.in/telebot.v3"
)

func NewRecipient(recipient string) *Recipient {
	return &Recipient{recipient: recipient}
}

type Recipient struct {
	recipient string
}

func (r *Recipient) Recipient() string {
	return r.recipient
}

func RunBot(
	token string,
	alertmanagerMessages <-chan webhook.Message,
	recipients []Recipient,
	logger *log.Logger,
) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		logger.Fatal(err)
		return
	}

	go func() {
		for message := range alertmanagerMessages {
			botMessage, err := FormatAlertHtml(message)
			if err != nil {
				logger.Printf("Execute message template failed failed: %s\n", err)
			}

			fmt.Printf("%+v\n", botMessage)

			for _, r := range recipients {
				_, err = bot.Send(&r, botMessage, tele.ModeHTML)
				if err != nil {
					logger.Printf("Send message to %s failed: %s\n", r.Recipient(), err)
				}
			}
		}
	}()

	logger.Printf("Running telegram bot with %d recipients\n", len(recipients))

	bot.Start()
}
