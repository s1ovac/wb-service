package config

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	BindAddress string `toml:"bind_addr"`
	Protocol    string `toml:"protocol"`
}

type StorageConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	DataBase string `toml:"database"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

var (
	configStoragePath string
)

func init() {
	flag.StringVar(&configStoragePath, "config-storage-path", "/home/s1ovac/github.com/wb-service/order-subscribe/config.toml", "path to storage config file")
}

func NewStorageConfig() (s *StorageConfig) {
	flag.Parse()
	_, err := toml.DecodeFile(configStoragePath, &s)
	if err != nil {
		return nil
	}
	return
}

func NewServerConfig() (s *ServerConfig) {
	flag.Parse()
	_, err := toml.DecodeFile(configStoragePath, &s)
	if err != nil {
		return nil
	}
	return
}
