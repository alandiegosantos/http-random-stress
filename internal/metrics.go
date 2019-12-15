package internal

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
)

var (
	httpReqCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_rq",
			Help: "The total number of processed requests",
		},
		[]string{
			"http_response_class",
			"http_url",
		},
	)
)

func GetHttpReqCounter(httpResponseClass int, url string) float64 {

	metric := &dto.Metric{}

	httpReqCounter.With(
		prometheus.Labels{
			"http_response_class": strconv.Itoa(httpResponseClass),
			"http_url":            url,
		},
	).Write(metric)

	return metric.GetCounter().GetValue()

}

func IncrementHttpReqCounter(httpResponseClass int, url string) {

	httpReqCounter.With(
		prometheus.Labels{
			"http_response_class": strconv.Itoa(httpResponseClass),
			"http_url":            url,
		},
	).Inc()

}
