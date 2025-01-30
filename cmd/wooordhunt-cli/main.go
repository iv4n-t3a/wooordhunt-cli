package main

import (
	"log"

	"github.com/iv4n-t3a/wooordhunt-cli/config"
	"github.com/iv4n-t3a/wooordhunt-cli/internal/cli"
)


func main() {
  conf, err := config.ParseConfig()

  if err != nil {
    log.Fatal(err)
  }

  CLI, err := cli.NewCLI(conf)

  if err != nil {
    log.Fatal(err)
  }

  CLI.Run()
}
