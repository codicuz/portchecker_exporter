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
    namespace = "was7"
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
		
	check1522port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check1522port_" + *wasInstance,
            Help: "Check 1522 Oracle listener port",
        })
    prometheus.MustRegister(check1522port)
	
	check9043port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check9043port_" + *wasInstance,
            Help: "Check 9043 WAS console port",
        })
    prometheus.MustRegister(check9043port)
    
    check9044port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check9044port_" + *wasInstance,
            Help: "Check 9044 WAS console port",
        })
    prometheus.MustRegister(check9044port)
    
    check9045port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check9045port_" + *wasInstance,
            Help: "Check 9045 WAS console port",
        })
	prometheus.MustRegister(check9045port)
		
	check10000port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check10000port_" + *wasInstance,
            Help: "Check 10000 WAS console port",
        })
	prometheus.MustRegister(check10000port)

	check10003port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check10003port_" + *wasInstance,
            Help: "Check 10003 WAS console port",
        })
	prometheus.MustRegister(check10003port)

	check10032port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check10032port_" + *wasInstance,
            Help: "Check 10032 WAS console port",
        })
    prometheus.MustRegister(check10032port)

	check10039port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check10039port_" + *wasInstance,
            Help: "Check 10039 Web Portal port",
        })
	prometheus.MustRegister(check10039port)
	
	check80port := prometheus.NewGauge (
        prometheus.GaugeOpts {
            Namespace: namespace,
            Subsystem: exporter,
            Name: "check80port_" + *wasInstance,
            Help: "Check 80 IHS port",
        })
	prometheus.MustRegister(check80port)
	
	go func () {
		for {
			raw_connect("tcp", host, "1522", check1522port)
            raw_connect("tcp", host, "9043", check9043port)
            raw_connect("tcp", host, "9044", check9044port)
            raw_connect("tcp", host, "9045", check9045port)
			raw_connect("tcp", host, "10000", check10000port)
			raw_connect("tcp", host, "10003", check10003port)
			raw_connect("tcp", host, "10032", check10032port)
			raw_connect("tcp", host, "10039", check10039port)
			raw_connect("tcp", host, "80", check80port)
			time.Sleep(globalsleep)
		}
	} ()

	http.Handle(*metricPath, promhttp.Handler())
    log.Infoln("Listening on", *listenAddress)
    log.Fatal(http.ListenAndServe(*listenAddress, nil))
}