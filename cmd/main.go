package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"headless-todo-file-service/internal/adapters/endpoints"
	"log"
	"net/http"
)

func main() {
	initConfig()

	client, closeConnection := ConnectMongo()
	nc, sc, closeNats := ConnectNats()
	defer func() {
		closeConnection()
		closeNats()
	}()

	c := Init(client, nc, sc)

	http.Handle("/create-file", endpoints.CreateFileHandler(c))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("HTTP_PORT")), nil))
}
