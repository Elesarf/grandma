package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	queryCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grandma_query_count",
		Help: "The total number of query",
	})

	photoQuery = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grandma_photo_query_count",
		Help: "The total number of photo query",
	})

	queryProcessingTime = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "query_processing_time_us",
			Help: "The time of query",
		})

	queryVideoProcessingTime = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "query_video_processing_time_us",
			Help: "The time of video query",
		})
)
