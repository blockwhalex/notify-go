package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	DB map[string]DBConfig `toml:"db"`
}

type DBConfig struct {
	DriverName string `toml:"driver_name"`
	Dsn        string `toml:"dsn"`
}

func LoadConfig(path string) AppConfig {
	conf := new(AppConfig)
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		log.Fatal("读取配置失败", err)
	}
	return *conf
}
