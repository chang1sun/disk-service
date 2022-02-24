package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var config *Config

type Config struct {
	Version             string `yaml:"version"`
	InitUserStorageSize int64  `yaml:"init_user_storage_size"`
	AuthKey             string `yaml:"auth_secret_key"`
	PwSalt              string `yaml:"pw_salt"`
	RequestOrigin string `yaml:"request_origin"`
	Mysql               struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"db_name"`
		Addr     string `yaml:"addr"`
	} `yaml:"mysql"`
	MongoDB struct {
		Database string `yaml:"db_name"`
		Addr     string `yaml:"addr"`
	} `yaml:"mongo"`
	Redis struct {
		Addr     string `yaml:"addr"`
		DBShare  int    `yaml:"db_share"`
		DBUpload int    `yaml:"db_upload"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
}

func GetConfig() *Config {
	return config
}

func InitConfig() error {
	config = &Config{}
	yamlFile, err := ioutil.ReadFile("application.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return err
	}
	return nil
}
