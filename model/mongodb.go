package model

type MongoConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"db"`
	Username string `yaml:"user"`
	Password string `yaml:"pwd"`
}
