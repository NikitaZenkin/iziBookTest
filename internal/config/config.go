package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB DB `yaml:"db"`
}

type DB struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DocumentsDbName string `yaml:"documents_db_name"`
	UsersDbName     string `yaml:"users_db_name"`
}

func GetConfig(filePath string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
