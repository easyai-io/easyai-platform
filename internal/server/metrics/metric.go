package metrics

import (
	"net/http"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PromRegistry is the global prometheus registry.
var PromRegistry = prometheus.NewRegistry()

func init() {
	PromRegistry.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	))
}

// Register registers the provided collectors with the global prometheus
func Register(c ...prometheus.Collector) {
	PromRegistry.MustRegister(c...)
}

// Handler returns an http.Handler exposing the registered metrics.
func Handler() http.Handler {

	// Expose the registered metrics via HTTP.
	return promhttp.HandlerFor(
		PromRegistry,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	)
}
