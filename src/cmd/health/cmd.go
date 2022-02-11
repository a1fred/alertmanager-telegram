package health

import (
	"fmt"
	"net/http"
	"path"

	"github.com/jessevdk/go-flags"
)

type Cmd struct {
	revision string

	DaemonHostPort string `long:"daemon_hostport" description:"Daemon host:port" env:"LISTEN" default:"127.0.0.1:8080"`
}

func NewHealthCmd(parser *flags.Parser, revision string) error {
	_, err := parser.AddCommand(
		"health",
		"health",
		"Check daemon health",
		&Cmd{
			revision: revision,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Cmd) Execute(args []string) error {
	healthUrl := "http://" + path.Join(s.DaemonHostPort, "/health")
	resp, err := http.Get(healthUrl)
	if err != nil {
		fmt.Printf("GET %s -> err:%s\n", healthUrl, err.Error())
		return err
	}

	fmt.Printf("GET %s -> %s\n", healthUrl, resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response status code %d", resp.StatusCode)
	}
	return nil
}
