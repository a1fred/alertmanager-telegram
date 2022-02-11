package telegramBot

import (
	"log"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/client_golang/prometheus"
	tele "gopkg.in/telebot.v3"
)

func RunBot(
	token string,
	alertmanagerMessages <-chan webhook.Message,
	recipients []Recipient,
	logger *log.Logger,
	messagesSentCounter, messagesSendingErrorCounter prometheus.Counter,
) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		logger.Fatalf("Bot initialize error: %s", err.Error())
		return
	}

	go func() {
		for message := range alertmanagerMessages {
			botMessage, err := FormatAlertHtml(message)
			if err != nil {
				logger.Printf("Execute message template failed failed: %s\n", err)
				messagesSendingErrorCounter.Inc()
				continue
			}

			for _, r := range recipients {
				_, err = bot.Send(&r, botMessage, tele.ModeHTML)
				if err != nil {
					logger.Printf("Send message to %s failed: %s\n", r.Recipient(), err)
					messagesSendingErrorCounter.Inc()
					continue
				}
				messagesSentCounter.Inc()
			}
		}
	}()

	logger.Printf("Running telegram bot with %d recipients\n", len(recipients))

	bot.Start()
}
