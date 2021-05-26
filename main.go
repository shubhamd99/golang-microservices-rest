package main

import (
	"GoMicroservices/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// Logger
	l := log.New(os.Stdout, "product-api ", log.LstdFlags) // Ex. product-api 2021/05/26 22:44:43 Hello World
	// Handler
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// ServeMux is an HTTP request multiplexer.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern that most closely matches the URL.
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	// A Server defines parameters for running an HTTP server. The zero value for Server is a valid configuration.
	// https://golang.org/pkg/net/http/#Server
	s := &http.Server{
		Addr: ":9090",
		Handler: sm,
		IdleTimeout: 120 *time.Second,
		ReadTimeout: 1 *time.Second,
		WriteTimeout: 1 *time.Second,
	}

	// Goroutines
	// This new goroutine will execute concurrently with the calling one.
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			// Fatal is equivalent to Print() followed by a call to os.Exit(1)
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	// Signal channel will broadcast the message on sigChannel whenever operating system kill or Interrupt command is received
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Receive from sigChan, and assign value to sig variable.
	sig := <- sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	// Package context defines the Context type, which carries deadlines, cancellation signals,
	// and other request-scoped values across API boundaries and between processes.
	tc, _ := context.WithTimeout(context.Background(), 30 *time.Second) // 30 seconds timeout

	// Shutdown gracefully shuts down the server without interrupting any active connections.
	// https://golang.org/pkg/net/http/#Server.Shutdown
	s.Shutdown(tc)
}
