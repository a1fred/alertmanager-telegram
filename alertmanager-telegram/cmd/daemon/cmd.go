package daemon

import (
	"log"
	"net/http"
	"os"

	"github.com/a1fred/alertmanager-telegram/alertmanager-telegram/httpServer"
	"github.com/a1fred/alertmanager-telegram/alertmanager-telegram/telegramBot"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TelegramOptions struct {
	Token  string   `long:"token" description:"Telegram bot token" env:"TOKEN" required:"true"`
	ChatId []string `short:"r" long:"recipient" description:"Telegram chat ids, repeat -r to set multiple, for environment set comma separated ids" env-delim:"," env:"CHAT_ID"`
}

type Cmd struct {
	revision string

	Listen string `long:"listen" description:"Webhook listen" env:"LISTEN" default:"127.0.0.1:8080"`

	telegramOptions *TelegramOptions
}

func NewDaemonCmd(parser *flags.Parser, revision string) error {
	telegramOptions := &TelegramOptions{}

	command, err := parser.AddCommand(
		"daemon",
		"Daemon",
		"Run daemon",
		&Cmd{
			revision:        revision,
			telegramOptions: telegramOptions,
		},
	)
	if err != nil {
		return err
	}

	g, err := command.AddGroup("Telegram daemon", "Telegram options", telegramOptions)
	if err != nil {
		return err
	}
	g.Namespace = "telegram"
	g.EnvNamespace = "TELEGRAM"

	return nil
}

func (s *Cmd) Execute(args []string) error {
	if len(s.telegramOptions.ChatId) == 0 {
		log.Fatal("Please specify some recipient ids")
	}

	recipients := make([]telegramBot.Recipient, 0)
	for _, r := range s.telegramOptions.ChatId {
		recipients = append(recipients, *telegramBot.NewRecipient(r))
	}

	promRegistry := prometheus.NewRegistry()
	alertsReceivedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_telegram_alerts_received",
		Help: "Number of alerts received",
	})
	messagesSentCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_telegram_messages_sent",
		Help: "Number of messages sent to telegram recipients",
	})
	messagesSendingErrorCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_telegram_messages_sending_error",
		Help: "Number of errors message sending to telegram recipients",
	})
	promRegistry.MustRegister(alertsReceivedCounter, messagesSentCounter, messagesSendingErrorCounter)

	httpLogger := log.New(os.Stdout, "   [http]    ", log.LstdFlags)
	teleLogger := log.New(os.Stdout, " [telegram]  ", log.LstdFlags)
	alertmanagerMessages := make(chan webhook.Message)

	go telegramBot.RunBot(
		s.telegramOptions.Token,
		alertmanagerMessages,
		recipients,
		teleLogger,
		messagesSentCounter,
		messagesSendingErrorCounter,
	)

	serveMux := httpServer.NewHttpServeMux(httpLogger, alertmanagerMessages, alertsReceivedCounter)
	serveMux.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))
	httpLogger.Printf("Http listen: http://%s\n", s.Listen)
	return http.ListenAndServe(s.Listen, serveMux)
}
