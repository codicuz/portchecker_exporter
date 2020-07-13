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
	hostInstance         = flag.String("host-instance", getEnv("HOST_INSTANCE", "HOST_id"), "HOST identificator")
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
            Name: "check80port_" + *hostInstance,
            Help: "Check 80 port",
        })
	prometheus.MustRegister(check80port)
    
    check8081port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check8081port_" + *hostInstance,
            Help: "Check 8081 port",
        })
	prometheus.MustRegister(check8081port)
    
    check8082port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check8082port_" + *hostInstance,
            Help: "Check 8082 port",
        })
	prometheus.MustRegister(check8082port)
    
    check10080port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check10080port_" + *hostInstance,
            Help: "Check 10080 port",
        })
	prometheus.MustRegister(check10080port)
    
    check8080port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check8080port_" + *hostInstance,
            Help: "Check 8080 port",
        })
    prometheus.MustRegister(check8080port)
    
    check22port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check22port_" + *hostInstance,
            Help: "Check 22 port",
        })
	prometheus.MustRegister(check22port)
	
	go func () {
		for {
            raw_connect("tcp", host, "80", check80port)
            raw_connect("tcp", host, "8081", check8081port)
            raw_connect("tcp", host, "8082", check8082port)
            raw_connect("tcp", host, "10080", check10080port)
            raw_connect("tcp", host, "8080", check8080port)
			time.Sleep(globalsleep)
		}
	} ()

	http.Handle(*metricPath, promhttp.Handler())
    log.Infoln("Listening on", *listenAddress)
    log.Fatal(http.ListenAndServe(*listenAddress, nil))
}