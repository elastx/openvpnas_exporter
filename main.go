package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	exporterVersion "github.com/elastx/openvpnas_exporter/version"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version string
)

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if version == "" && ok && buildInfo != nil && buildInfo.Main.Version != "" {
		version = buildInfo.Main.Version
	}
	if version == "" {
		version = "devel"
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
	kingpin.Version(fmt.Sprintf("openvpnas_exporter version %v", version))
	kingpin.Parse()

	log.Printf("Starting openvpnas_exporter version %v\n", version)
	log.Printf("Listen address: %v\n", *listenAddress)
	log.Printf("Metrics path: %v\n", *metricsPath)
	log.Printf("XML-RPC path: %v\n", *xmlrpcPath)

	buildInfo := exporterVersion.NewVersionCollector("openvpnas", version)
	prometheus.MustRegister(buildInfo)

	exporter, err := New(*xmlrpcPath)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
