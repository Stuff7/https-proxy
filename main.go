package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	port := flag.String("port", "8080", "Local HTTP server port to forward requests to")
	httpsPort := flag.String("https-port", "8443", "HTTPS server port to listen on")
	flag.Parse()

	localAddr := fmt.Sprintf("http://localhost:%s", *port)

	target, err := url.Parse(localAddr)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	httpsServer := &http.Server{
		Addr: fmt.Sprintf(":%s", *httpsPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		}),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	crt := getOsEnv("LOCALHOST_SSL_CRT")
	key := getOsEnv("LOCALHOST_SSL_KEY")

	fmt.Printf("Serving HTTPS on https://localhost:%s\n", *httpsPort)
	err = httpsServer.ListenAndServeTLS(crt, key)
	if err != nil {
		log.Fatal("ListenAndServeTLS failed: ", err)
	}
}

func getOsEnv(name string) string {
	env := os.Getenv(name)

	if env == "" {
		log.Fatal("Missing env var ", name)
	}

	return env
}
