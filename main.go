
//
// nutanix-exporter
//
// Prometheus Exportewr for Nutanix API
//
// Author: Martin Weber <martin.weber@de.clara.net>
// Company: Claranet GmbH
//

package main

import (
	"./nutanix"
	"./collector"
	"flag"
	"net/http"
//	"time"
//	"regexp"
//	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

var (
	namespace		= "nutanix"
	nutanixUrl		= flag.String("nutanix.url", "", "Nutanix URL to connect to API https://nutanix.local.host:9440")
	nutanixUser		= flag.String("nutanix.username", "", "Nutanix API User")
	nutanixPassword		= flag.String("nutanix.password", "", "Nutanix API User Password")
	listenAddress		= flag.String("listen-address", ":9405", "The address to lisiten on for HTTP requests.")
)

var (
	// Nutanix API
	nutanixApi		*nutanix.Nutanix
)

func main() {
	flag.Parse()

	log.Debug("Create Nutanix instance")
	nutanixApi = nutanix.NewNutanix(*nutanixUrl, *nutanixUser, *nutanixPassword)

	prometheus.MustRegister( collector.NewStorageExporter(nutanixApi) )
	prometheus.MustRegister( collector.NewClusterExporter(nutanixApi) )

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Nutanix Exporter</title></head>
		<body>
		<h1>Nutanix Exporter</h1>
		<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>`))
	})

	log.Printf("Starting Server: %s", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
