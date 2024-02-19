package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Version string
)

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if Version == "" && ok && buildInfo != nil && buildInfo.Main.Version != "" {
		Version = buildInfo.Main.Version
	}
	if Version == "" {
		Version = "devel"
	}
}

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address to listen on for web interface and telemetry.",
		).Default(":9176").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		xmlrpcPath = kingpin.Flag(
			"openvpnas.xmlrpc-path",
			"Path to the XML-RPC unix domain socket file.",
		).Default("/usr/local/openvpn_as/etc/sock/sagent.localroot").String()
	)
	kingpin.Version(fmt.Sprintf("openvpnas_exporter version %v", Version))
	kingpin.Parse()

	log.Printf("Starting openvpnas_exporter version %v\n", Version)
	log.Printf("Listen address: %v\n", *listenAddress)
	log.Printf("Metrics path: %v\n", *metricsPath)
	log.Printf("XML-RPC path: %v\n", *xmlrpcPath)

	// TODO(holmsten): Once https://github.com/prometheus/client_golang/pull/1427 is in a release (> 1.8)
	// we can update and use the following below instead:
	// buildInfo := version.NewCollector("openvpnas")
	buildInfo := prometheus.NewBuildInfoCollector()
	prometheus.MustRegister(buildInfo)

	exporter, err := New(*xmlrpcPath)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
