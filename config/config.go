package config

import "github.com/alecthomas/kong"

type Config struct {
	Insecure bool `short:"i" help:"Disable ssl verification"`
}

func ParseConfig() (Config, error) {
	var conf Config
	kong.Parse(&conf)
	return conf, nil
}
