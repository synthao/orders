package config

import "go.uber.org/config"

type Server struct {
	Port int `yaml:"port"`
}

func NewServerConfig() (*Server, error) {
	provider, err := config.NewYAML(config.File(name))
	if err != nil {
		return nil, err
	}

	var c Server

	err = provider.Get("server").Populate(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
