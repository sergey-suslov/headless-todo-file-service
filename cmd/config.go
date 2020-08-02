package main

import (
	"github.com/spf13/viper"
	"log"
)

func initConfig() {
	viper.SetDefault("HTTP_PORT", "8081")

	viper.SetDefault("DB_NAME", "db")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_USER", "user")

	viper.SetDefault("METRICS_NAMESPACE", "namespace")
	viper.SetDefault("METRICS_SUBSYSTEM", "service")

	viper.SetDefault("JWT_SECRET", "XXX")

	viper.SetDefault("NATS_CONNECTION_STRING", "localhost:4222")
	viper.SetDefault("NATS_CLUSTER_ID", "todo-cluster")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
