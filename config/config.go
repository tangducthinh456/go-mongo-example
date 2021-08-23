package config

import (
	"io/ioutil"

	"thinhtd4/customlog"
	"thinhtd4/model"

	"gopkg.in/yaml.v2"
)

type Config struct {
	model.MongoConfig  `yaml:"mongodb"`
	model.RedisConfig  `yaml:"redis"`
	model.ServerConfig `yaml:"server"`
}

var (
	serverConfig *model.ServerConfig = new(model.ServerConfig)
	mongoConfig  *model.MongoConfig  = new(model.MongoConfig)
	redisConfig  *model.RedisConfig  = new(model.RedisConfig)
	config       *Config
)

func Init() {
	customlog.Info("Start initialize config")
	initConfig, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		customlog.Err(err)
		panic(err)
	}

	err = yaml.Unmarshal(initConfig, &config)
	if err != nil {
		customlog.Err(err)
	}
	
	*serverConfig = config.ServerConfig
	*mongoConfig = config.MongoConfig
	*redisConfig = config.RedisConfig
}

func ServerConfig() *model.ServerConfig {
	return serverConfig
}

func MongoConfig() *model.MongoConfig {
	return mongoConfig
}

func RedisConfig() *model.RedisConfig {
	return redisConfig
}
