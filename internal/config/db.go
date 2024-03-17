package config

import (
	"fmt"
	"go.uber.org/config"
	"os"
)

type DB struct {
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Charset  string `yaml:"charset"`
	Params   string `yaml:"params"`
}

func NewDBConfig() (*DB, error) {
	provider, err := config.NewYAML(config.File(filename))
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
	host := getFromEnv("DB_HOST", c.Host)
	port := getFromEnv("DB_PORT", c.Port)
	user := getFromEnv("DB_USER", c.User)
	password := getFromEnv("DB_PASSWORD", c.Password)
	dbName := getFromEnv("DB_NAME", c.Name)

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbName,
	)
}

func getFromEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return def
}
