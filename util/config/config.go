package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	BotToken     string `yaml:"bot_token"`
	MongoURI     string `yaml:"mongo_uri"`
	GuildID      string `yaml:"guild_id"`
	Transactions string `yaml:"transactions_channel_id"`
	Suspensions  string `yaml:"suspensions_channel_id"`
}

func NewConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
