package middlewares

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"headless-todo-file-service/internal/entities"
	"headless-todo-file-service/internal/services"
	"io"
	"time"
)

type PrometheusMetrics struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
}

func NewPrometheusMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram, countResult metrics.Histogram) *PrometheusMetrics {
	return &PrometheusMetrics{RequestCount: requestCount, RequestLatency: requestLatency, CountResult: countResult}
}

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           services.FilesService
}

func (mw *InstrumentingMiddleware) Create(ctx context.Context, name, userId, taskId string, file io.Reader) (output *entities.File, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Create", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.Create(ctx, name, userId, taskId, file)
}
