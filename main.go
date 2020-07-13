package main

import (
	"net"
	"net/http"
	"flag"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"
) 

var (
	version string = "1.0-dev"
    listenAddress       = flag.String("listen-address", getEnv("LISTEN_ADDRESS", ":9800"), "Address to listen on for web interface and telemetry. (env: LISTEN_ADDRESS)")
    metricPath          = flag.String("telemetry-path", getEnv("TELEMETRY_PATH", "/metrics"), "Path under which to expose metrics. (env: TELEMETRY_PATH)")
	wasInstance         = flag.String("was-instance", getEnv("WAS_INSTANCE", "WAS_id"), "WAS identificator")
)

const (
    namespace = "infra"
    exporter = "exporter"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func raw_connect(network, host, port string, metrics prometheus.Gauge) {
	l, err := net.Dial(network, host + ":" + port)		
	if err != nil {
		metrics.Set(0)
		log.Warn("Port " + port + " is Not Ok")
	} else {
		metrics.Set(1)
		log.Infoln("Port " + port + " is Ok")
		defer l.Close()
	}
}

func main() {
	host := os.Getenv("HOST_IP")
	globalsleep := 15000 * time.Millisecond
	log.Infoln("Hello, portChecker_exporter!")

    check80port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check80port_" + *wasInstance,
            Help: "Check 80 port",
        })
	prometheus.MustRegister(check80port)
	
	go func () {
		for {
			raw_connect("tcp", host, "80", check80port)
			time.Sleep(globalsleep)
		}
	} ()

	http.Handle(*metricPath, promhttp.Handler())
    log.Infoln("Listening on", *listenAddress)
    log.Fatal(http.ListenAndServe(*listenAddress, nil))
}