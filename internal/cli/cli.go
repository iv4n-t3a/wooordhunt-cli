package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/iv4n-t3a/wooordhunt-cli/config"
	"github.com/iv4n-t3a/wooordhunt-cli/internal/client"
)

type CLI struct {
	conf    config.Config
	client  client.Client
	program tea.Program
}

func NewCLI(conf config.Config) (CLI, error) {
	cli := CLI{
		conf: conf,
	}

	clnt, err := client.NewClient(conf)
	if err != nil {
		return CLI{}, err
	}
	cli.client = clnt

	initalModel := newInput(&cli)
	cli.program = *tea.NewProgram(&initalModel)

	return cli, nil
}

func (cli CLI) Run() {
	cli.program.Run()
}
