package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	certFile, keyFile, runtimeClass, port string
)

func main() {
	flag.StringVar(&certFile, "certFile", "/certs/server.pem", "Path to x509 https certificate")
	flag.StringVar(&keyFile, "keyFile", "/certs/server-key.pem", "Path to x509 https private key for SSL")
	flag.StringVar(&runtimeClass, "runtimeClass", "gvisor", "The runtimeclass of the environment")
	flag.StringVar(&port, "port", "8443", "Default port to listen to")

	flag.Parse()

	certs, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Printf("Error loading cert/key ", err.Error())
		os.Exit(1)
	}
	
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		TLSConfig:         &tls.Config{
			Certificates:                []tls.Certificate{certs},
		},
	}

	handler := AdmissionHandler{
		RuntimeClass: runtimeClass,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/admission", handler.admissionhandler)

	server.Handler = mux

	go func() {
		log.Printf("Listening on port %v", port)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Printf("Failed to listen and serve webhook server: %v", err)
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutting down webserver")
	server.Shutdown(context.Background())
}
