package daemon

import (
	"log"
	"net/http"
	"os"

	"github.com/a1fred/alertmanager-telegram/alertmanager-telegram/httpServer"
	"github.com/a1fred/alertmanager-telegram/alertmanager-telegram/telegramBot"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/alertmanager/notify/webhook"
)

type TelegramOptions struct {
	Token  string   `long:"token" description:"Telegram bot token" env:"TOKEN" required:"true"`
	ChatId []string `short:"r" long:"recipient" description:"Telegram chat ids" env:"CHAT_ID"`
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

	httpLogger := log.New(os.Stdout, "   [http]    ", log.LstdFlags)
	teleLogger := log.New(os.Stdout, " [telegram]  ", log.LstdFlags)
	alertmanagerMessages := make(chan webhook.Message)

	go telegramBot.RunBot(
		s.telegramOptions.Token,
		alertmanagerMessages,
		recipients,
		teleLogger,
	)

	serveMux := httpServer.NewHttpServeMux(httpLogger, alertmanagerMessages)
	httpLogger.Printf("Http listen: http://%s\n", s.Listen)
	return http.ListenAndServe(s.Listen, serveMux)
}
