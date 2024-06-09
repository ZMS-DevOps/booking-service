package main

import (
	"github.com/ZMS-DevOps/booking-service/startup"
	cfg "github.com/ZMS-DevOps/booking-service/startup/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	config := cfg.NewConfig()
	producer, _ := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"security.protocol": "sasl_plaintext",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     "user1",
		"sasl.password":     config.KafkaAuthPassword,
	})
	defer producer.Close()
	server := startup.NewServer(config)
	server.Start(producer)
}
