package config

import (
	"fmt"
	"go.uber.org/config"
)

type DB struct {
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Charset  string `yaml:"charset"`
	Params   string `yaml:"params"`
}

func NewDBConfig() (*DB, error) {
	provider, err := config.NewYAML(config.File(name))
	if err != nil {
		return nil, err
	}

	var c DB

	err = provider.Get("db").Populate(&c)
	if err != nil {
		panic(err)
	}

	return &c, nil
}

func (c *DB) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
	)
}
