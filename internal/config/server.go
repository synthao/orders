package config

import (
	"errors"
	"fmt"
	"os"
)

var ErrMissingEnvVariable = errors.New("missing environment variable")

type Server struct {
	Port string
}

func NewServerConfig() (*Server, error) {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return nil, fmt.Errorf("%w, PORT", ErrMissingEnvVariable)
	}

	c := Server{
		Port: port,
	}

	return &c, nil
}
