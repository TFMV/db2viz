package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	BigQuery BigQueryConfig `yaml:"bigquery"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type BigQueryConfig struct {
	ProjectID   string `yaml:"project_id"`
	DatasetID   string `yaml:"dataset_id"`
	TableID     string `yaml:"table_id"`
	Credentials string `yaml:"credentials"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return &cfg, err
}
