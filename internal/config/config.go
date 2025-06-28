package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env   string       `yaml:"env"`
	GRPC  *GRPCConfig  `yaml:"grpc"`
	Kafka *KafkaConfig `yaml:"kafka"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type KafkaConfig struct {
	Port             int      `yaml:"port"`
	Topics           []string `yaml:"topics"`
	WorkersCount     int      `yaml:"workers_count"`
	BootstrapServers string   `yaml:"bootstrap_servers"`
	GroupID          string   `yaml:"group_id"`
	AutoOffsetReset  string   `yaml:"auto_offset_reset"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)

	if err != nil {
		panic(err)
	}

	return &cfg
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
