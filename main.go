package main

import (
	"fmt"
	"log"
	"os"

	"github.com/a1fred/alertmanager-telegram/src/cmd/daemon"
	"github.com/a1fred/alertmanager-telegram/src/cmd/health"
	"github.com/jessevdk/go-flags"
)

var revision = "unknown"

func main() {
	_, err := os.Stderr.WriteString(fmt.Sprintf("alertmanager-telegram@%s\n", revision))
	if err != nil {
		log.Fatalln(err)
	}

	parser := flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)

	// server
	err = daemon.NewDaemonCmd(parser, revision)
	if err != nil {
		log.Fatalln(err)
	}

	// health
	err = health.NewHealthCmd(parser, revision)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
