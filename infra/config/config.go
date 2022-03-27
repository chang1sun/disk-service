package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var config *Config

type Config struct {
	Version             string `yaml:"version"`
	InitUserStorageSize int64  `yaml:"init_user_storage_size"`
	AuthKey             string `yaml:"auth_secret_key"`
	PwSalt              string `yaml:"pw_salt"`
	TLS                 struct {
		Key string `yaml:"key"`
		Crt string `yaml:"crt"`
	} `yaml:"tls"`
	Mysql struct {
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
	var confPath string
	if os.Getenv("RUN_MODE") == "prod" {
		confPath = "conf/prod.conf.yaml"
	} else {
		confPath = "conf/dev.conf.yaml"
	}
	if err := readYaml(confPath, config); err != nil {
		return err
	}
	return nil
}

func readYaml(file string, config *Config) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot find conf file, err: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return err
	}
	return nil
}
