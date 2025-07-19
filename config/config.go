package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/alecthomas/kong"
)

var confDirs = []string{
	"./config",
	"~/.config/wh",
	"~/.config/wooordhunt-cli",
	"/etc/wh",
	"/etc/wooordhunt-cli",
}

var errFileNotExists = errors.New("File not exists")

type Config struct {
	Insecure bool `json:"disable_ssl"`
}

type Options struct {
	Insecure bool   `short:"i" help:"Disable ssl verification"`
	Config   string `short:"c" help:"Config directory" type:"existingdir"`
}

func ParseConfig() (Config, error) {
	var opts Options
	var conf Config
  var err error

	kong.Parse(&opts)

	if len(opts.Config) == 0 {
    for i := range confDirs {
      conf, err = parseConfig(confDirs[i])
      if err == nil {
        break
      }
    }
	} else {
    conf, err = parseConfig(opts.Config)
    if err != nil {
      return Config{}, err
    }
  }

	if opts.Insecure {
		conf.Insecure = true
	}

	return conf, nil
}

func parseConfig(confDir string) (Config, error) {
  path := substituteHomeDir(confDir + "/config.json")
  data, err := os.ReadFile(path)

  if err != nil {
    return Config{}, err
  }

  var res Config
  err = json.Unmarshal([]byte(data), &res)

  if err != nil {
    return Config{}, err
  }

  return res, nil
}

func substituteHomeDir(path string) string {
  if (len(path) < 1 || '~' != path[0]) {
    return path
  }

  homedir, err := os.UserHomeDir()

  if err != nil {
    return path
  }

  return homedir + path[1:]
}
