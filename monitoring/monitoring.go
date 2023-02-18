package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	Connections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "margay_gateway_connections",
		Help: "The number of open connections to the WebSocket API gateway.",
	}, []string{"route"})
	Messages = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "margay_gateway_messages",
		Help: "The number of WebSocket messages handled by the API gateway.",
	}, []string{"route", "direction"})
)

func Init() {
	prometheus.MustRegister(Connections)
	prometheus.MustRegister(Messages)
}

func StartMetricsServer(port string) {
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, promhttp.Handler()))
	}()
}

func IncrementConnectionCount(route string) {
	Connections.WithLabelValues(route).Inc()
}

func DecrementConnectionCount(route string) {
	Connections.WithLabelValues(route).Dec()
}

func IncrementMessageCount(route string, direction string) {
	Messages.WithLabelValues(route, direction).Inc()
}
