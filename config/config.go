package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	PubSub   PubSubConfig   `yaml:"pubsub"`
}

type PostgresConfig struct {
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password"`
	DBName      string        `yaml:"dbname"`
	SSLMode     string        `yaml:"sslmode"`
	Tables      []TableConfig `yaml:"tables"`
	Concurrency int           `yaml:"concurrency"`
}

type TableConfig struct {
	Name    string `yaml:"name"`
	TopicID string `yaml:"topic_id"`
}

type PubSubConfig struct {
	ProjectID   string `yaml:"project_id"`
	Credentials string `yaml:"credentials"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
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
