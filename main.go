package main

import (
	"github.com/ZMS-DevOps/booking-service/startup"
	cfg "github.com/ZMS-DevOps/booking-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
