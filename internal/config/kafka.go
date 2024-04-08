package config

import (
	"go.uber.org/config"
)

type Kafka struct {
	Address string `yaml:"address"`
	Topic   string `yaml:"topic"`
}

func NewKafkaConfig() (*Kafka, error) {
	provider, err := config.NewYAML(config.File(filename))
	if err != nil {
		return nil, err
	}

	var c Kafka

	err = provider.Get("kafka").Populate(&c)
	if err != nil {
		panic(err)
	}

	return &c, nil
}

// GetAddress fetches the Kafka broker address from the environment variable KAFKA_ADDRESS,
// falling back to the address specified in the YAML configuration if the environment variable is not set.
func (k *Kafka) GetAddress() string {
	return getFromEnv("KAFKA_ADDRESS", k.Address)
}

// GetTopic retrieves the Kafka topic from the environment variable KAFKA_TOPIC,
// with a fallback to the topic specified in the YAML configuration if the environment variable is not present.
func (k *Kafka) GetTopic() string {
	return getFromEnv("KAFKA_TOPIC", k.Topic)
}
