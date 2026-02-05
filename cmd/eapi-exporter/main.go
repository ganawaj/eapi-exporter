package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aristanetworks/goeapi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	
	// Enable RSA key exchange for older Arista EOS devices
	// os.Setenv("GODEBUG", "tlsrsakex=1")

	mux := http.NewServeMux()

	protocol := os.Getenv("EAPI_PROTOCOL")
	if protocol != "http" && protocol != "https" {
		if protocol != "" {
			fmt.Printf("WARN: '%q' is invalid, setting to 'https'\n", protocol)
		}
		protocol = "https"
	}

	host := os.Getenv("EAPI_HOST")
	if host == "" {
		fmt.Println("EAPI_HOST needs to be specified")
		return
	}

	username := os.Getenv("EAPI_USERNAME")
	if username == "" {
		fmt.Println("EAPI_USERNAME needs to be specified")
		return
	}

	password := os.Getenv("EAPI_PASSWORD")
	if password == "" {
		fmt.Println("EAPI_PASSWORD needs to be specified")
		return
	}

	port := 443
	if portStr := os.Getenv("EAPI_PORT"); portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Printf("EAPI_PORT %q is not a valid number\n", portStr)
			return
		}
		port = p
	}

	// Create node connection
	node, err := goeapi.Connect(protocol, host, username, password, port)
	if err != nil {
		log.Fatal(err)
	}

	// Health endpoints
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, _ *http.Request) {
		_, err := node.RunCommands([]string{"show hostname"}, "json")
		if err != nil {
			http.Error(w, "unable to communicate with node", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	prometheus.MustRegister(
		NewInterfaceCollector(node),
		NewSystemCollector(node),
	)

	mux.Handle("GET /metrics", promhttp.Handler())
	log.Println("listening on :9120")
	log.Fatal(http.ListenAndServe(":9120", mux))
}
