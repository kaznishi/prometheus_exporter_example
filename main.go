package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

const (
	namespace = "abcd"
)

type myCollector struct {
	exampleCount prometheus.Counter
	exampleGauge prometheus.Gauge
}

func newMyCollector() *myCollector {
	return &myCollector{
		exampleCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "example_count",
			Help:      "example counter help",
		}),
		exampleGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "example_gauge",
			Help:      "example gauge help",
		}),
	}
}

func (c *myCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.exampleCount.Desc()
	ch <- c.exampleGauge.Desc()
}

func (c *myCollector) Collect(ch chan<- prometheus.Metric) {
	dummyStaticNumber := float64(1234)

	ch <- prometheus.MustNewConstMetric(
		c.exampleCount.Desc(),
		prometheus.CounterValue,
		float64(dummyStaticNumber),
	)
	ch <- prometheus.MustNewConstMetric(
		c.exampleGauge.Desc(),
		prometheus.GaugeValue,
		float64(dummyStaticNumber),
	)
}

func main() {
	flag.Parse()

	myCollector := newMyCollector()
	prometheus.MustRegister(myCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
