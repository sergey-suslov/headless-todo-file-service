package main

import (
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
)

func ConnectNats() (*nats.Conn, func()) {
	nc, err := nats.Connect(viper.GetString("NATS_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	return nc, nc.Close
}
