package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"

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
  cli.program = *tea.NewProgram(initalModel)

	return cli, nil
}

func (cli CLI) Run() {
	cli.program.Run()
}

func (cli CLI) PrettifyedPrint(r client.TipsList) {
	maxlen := 0
	for i := range r.Tips {
		maxlen = max(maxlen, len(r.Tips[i].Word))
	}

	cyan := color.New(color.FgCyan).SprintFunc()

	for i := range r.Tips {
		word := r.Tips[i].Word
		tips := r.Tips[i].Tips
		spacesCount := maxlen - len(word) + 1
		spaces := strings.Repeat(" ", spacesCount)
		fmt.Printf("%s: %s%s\n", cyan(word), spaces, tips)
	}
}
