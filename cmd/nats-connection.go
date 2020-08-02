package main

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

func ConnectNats() (*nats.Conn, stan.Conn, func()) {
	nc, err := nats.Connect(viper.GetString("NATS_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	clientID, err := nc.GetClientID()
	if err != nil {
		log.Fatal(err)
	}
	sc, err := stan.Connect(viper.GetString("NATS_CLUSTER_ID"), strconv.FormatUint(clientID, 10))
	return nc, sc, func() {
		_ = sc.Close()
		nc.Close()
	}
}
