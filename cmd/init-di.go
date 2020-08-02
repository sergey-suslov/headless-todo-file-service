package main

import (
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/nats-io/nats.go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	file_manager "headless-todo-file-service/internal/adapters/file-manager"
	"headless-todo-file-service/internal/adapters/middlewares"
	"headless-todo-file-service/internal/adapters/repositories"
	"headless-todo-file-service/internal/services"
	"log"
	"os"
	"path/filepath"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init(client *mongo.Client, nc *nats.Conn) *dig.Container {
	c := dig.New()

	err := c.Provide(func() *mongo.Client {
		return client
	})
	handleError(err)

	err = c.Provide(func() *nats.Conn {
		return nc
	})
	handleError(err)

	err = c.Provide(func() *mongo.Database {
		return client.Database(viper.GetString("DB_NAME"))
	})
	handleError(err)

	err = c.Provide(func() repositories.FileManager {
		return file_manager.NewFileManagerLocal(filepath.Join(filepath.Dir(".."), viper.GetString("LOCAL_STORAGE_PATH")))
	})
	handleError(err)

	err = c.Provide(repositories.NewFilesRepositoryMongo)
	handleError(err)

	err = c.Provide(repositories.NewTasksRepositoryNats)
	handleError(err)

	err = c.Provide(func() *middlewares.PrometheusMetrics {
		fieldKeys := []string{"method", "error"}
		requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: viper.GetString("METRICS_NAMESPACE"),
			Subsystem: viper.GetString("METRICS_SUBSYSTEM"),
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
		requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: viper.GetString("METRICS_NAMESPACE"),
			Subsystem: viper.GetString("METRICS_SUBSYSTEM"),
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys)
		countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: viper.GetString("METRICS_NAMESPACE"),
			Subsystem: viper.GetString("METRICS_SUBSYSTEM"),
			Name:      "count_result",
			Help:      "The result of each count method.",
		}, []string{}) // no fields here
		return middlewares.NewPrometheusMetrics(requestCount, requestLatency, countResult)
	})

	err = c.Provide(services.NewFilesServiceService)
	handleError(err)

	err = c.Provide(services.NewTasksService)
	handleError(err)

	err = c.Provide(func() kitlog.Logger {
		logger := kitlog.NewLogfmtLogger(os.Stderr)
		return logger
	})
	handleError(err)

	return c
}
