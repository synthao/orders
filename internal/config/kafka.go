package config

import (
	"os"
)

type Kafka struct {
	Address string
	Topic   string
}

func NewKafkaConfig() (*Kafka, error) {
	var c Kafka

	c.Topic = os.Getenv("KAFKA_TOPIC")
	c.Address = os.Getenv("KAFKA_ADDRESS")

	return &c, nil
}
