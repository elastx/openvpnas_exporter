package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"alexejk.io/go-xmlrpc"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	xmlrpcPath           string
	version              *prometheus.Desc
	up                   *prometheus.Desc
	ccCurrent            *prometheus.Desc
	ccLimit              *prometheus.Desc
	ccMax                *prometheus.Desc
	ccTotal              *prometheus.Desc
	lastSuccessfulUpdate *prometheus.Desc
	nextUpdate           *prometheus.Desc
}

func New(xmlrpcPath string) (*Exporter, error) {
	version := prometheus.NewDesc(
		"openvpnas_server_version_info",
		"Contains OpenVPN AS server version info",
		[]string{"version"}, nil,
	)
	up := prometheus.NewDesc(
		"openvpnas_up",
		"Whether scaping OpenVPN AS metrics was successful.",
		nil, nil,
	)
	ccCurrent := prometheus.NewDesc(
		"openvpnas_server_connected_clients",
		"Number of currently connected clients to the server.",
		nil, nil,
	)
	ccLimit := prometheus.NewDesc(
		"openvpnas_server_connected_clients_limit",
		"Server concurrent client connection limit.",
		nil, nil,
	)
	ccMax := prometheus.NewDesc(
		"openvpnas_subscription_connected_clients_limit",
		"Maximum number of concurrent client connections allowed by the OpenVPN AS subscription.",
		nil, nil,
	)
	ccTotal := prometheus.NewDesc(
		"openvpnas_subscription_connected_clients",
		"Total number of client connections currently in use across the OpenVPN AS subscription.",
		nil, nil,
	)
	lastSuccessfulUpdate := prometheus.NewDesc(
		"openvpnas_subscription_status_last_update_time_seconds",
		"UNIX timestamp when the OpenVPN AS subscription was last synced.",
		nil, nil,
	)
	nextUpdate := prometheus.NewDesc(
		"openvpnas_subscription_status_next_update_time_seconds",
		"UNIX timestamp of the next planned OpenVPN AS subscription sync.",
		nil, nil,
	)

	return &Exporter{
		xmlrpcPath:           xmlrpcPath,
		version:              version,
		up:                   up,
		ccCurrent:            ccCurrent,
		ccLimit:              ccLimit,
		ccMax:                ccMax,
		ccTotal:              ccTotal,
		lastSuccessfulUpdate: lastSuccessfulUpdate,
		nextUpdate:           nextUpdate,
	}, nil
}

func (exporter *Exporter) Collect(ch chan<- prometheus.Metric) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", exporter.xmlrpcPath)
			},
		},
	}

	opts := []xmlrpc.Option{xmlrpc.HttpClient(&httpc)}
	client, _ := xmlrpc.NewClient("http://localhost/", opts...)
	defer client.Close()

	err := exporter.GetVersion(*client, ch)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(exporter.up, prometheus.GaugeValue, 0)
		return
	}

	err = exporter.GetSubscriptionStatus(*client, ch)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(exporter.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(exporter.up, prometheus.GaugeValue, 1)
}

func (exporter *Exporter) Describe(ch chan<- *prometheus.Desc) {
}

func (exporter *Exporter) GetVersion(client xmlrpc.Client, ch chan<- prometheus.Metric) error {
	result := &struct {
		Version string
	}{}

	err := client.Call("GetASLongVersion", nil, result)
	if err != nil {
		log.Printf("Unable to fetch version metrics: %v", err)
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.version, prometheus.GaugeValue, 1, result.Version)
	return nil
}

func (exporter *Exporter) GetSubscriptionStatus(client xmlrpc.Client, ch chan<- prometheus.Metric) error {
	result := &struct {
		SubscriptionStatus struct {
			AgentDisabled           bool     `xmlrpc:"agent_disabled"`
			AgentId                 string   `xmlrpc:"agent_id"`
			CcLimit                 int      `xmlrpc:"cc_limit"`
			CurrentCc               int      `xmlrpc:"current_cc"`
			Error                   string   `xmlrpc:"error"`
			FallbackCc              int      `xmlrpc:"fallback_cc"`
			GracePeriod             int      `xmlrpc:"grace_period"`
			LastSuccessfulUpdate    int      `xmlrpc:"last_successful_update"`
			LastSuccessfulUpdateAge int      `xmlrpc:"last_successful_update_age"`
			MaxCc                   int      `xmlrpc:"max_cc"`
			Name                    string   `xmlrpc:"name"`
			NextUpdate              int      `xmlrpc:"next_update"`
			NextUpdateIn            int      `xmlrpc:"next_update_in"`
			Notes                   []string `xmlrpc:"notes"`
			Overdraft               bool     `xmlrpc:"overdraft"`
			Server                  string   `xmlrpc:"server"`
			State                   string   `xmlrpc:"state"`
			Subkey                  string   `xmlrpc:"subkey"`
			TotalCc                 int      `xmlrpc:"total_cc"`
			Type                    string   `xmlrpc:"type"`
			UpdatesFailed           int      `xmlrpc:"updates_failed"`
		}
	}{}

	err := client.Call("GetSubscriptionStatus", nil, result)
	if err != nil {
		log.Printf("Unable to fetch subscription metrics: %v", err)
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.ccCurrent, prometheus.GaugeValue, float64(result.SubscriptionStatus.CurrentCc))
	ch <- prometheus.MustNewConstMetric(exporter.ccLimit, prometheus.GaugeValue, float64(result.SubscriptionStatus.CcLimit))
	ch <- prometheus.MustNewConstMetric(exporter.ccMax, prometheus.GaugeValue, float64(result.SubscriptionStatus.MaxCc))
	ch <- prometheus.MustNewConstMetric(exporter.ccTotal, prometheus.GaugeValue, float64(result.SubscriptionStatus.TotalCc))
	ch <- prometheus.MustNewConstMetric(exporter.lastSuccessfulUpdate, prometheus.GaugeValue, float64(result.SubscriptionStatus.LastSuccessfulUpdate))
	ch <- prometheus.MustNewConstMetric(exporter.nextUpdate, prometheus.GaugeValue, float64(result.SubscriptionStatus.NextUpdate))
	return nil
}
