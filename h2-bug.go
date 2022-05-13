package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Some code inspired by
// https://medium.com/@thrawn01/http-2-cleartext-h2c-client-example-in-go-8167c7a4181e

func client() {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	urlStr := "http://localhost:8001/debug"
	data := url.Values{"test": {"value1"}}
	resp, err := client.PostForm(urlStr, data)
	if err != nil {
		log.Fatalf("client err: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("client body read err: %v", err)
	}
	fmt.Printf("Response %d proto %s, body:\n%s\n", resp.StatusCode, resp.Proto, string(body))
}

func server() {
	h2s := &http2.Server{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got %v %v request from %v\n", r.Method, r.Proto, r.RemoteAddr)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("server body read err: %v", err)
		}
		fmt.Printf("read %d from body\n", len(body))
		//w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Hello, %v; HTTP Version: %v\n\nbody:\n%s\n", r.URL.Path, r.Proto, string(body))
	})

	server := &http.Server{
		Addr:    ":8001",
		Handler: h2c.NewHandler(handler, h2s),
	}
	fmt.Printf("H2c Server starting on %s\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func main() {
	server()
	//	client()
}
