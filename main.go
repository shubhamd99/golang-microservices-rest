package main

import (
	"GoMicroservices/handlers"
	"context"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()

	// Logger
	l := log.New(os.Stdout, "product-api ", log.LstdFlags) // Ex. product-api 2021/05/26 22:44:43 Hello World

	// Handler
	ph := handlers.NewProducts(l)

	// ServeMux is an HTTP request multiplexer.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern that most closely matches the URL.
	// create a new server mux and register the handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	// create a new server
	// A Server defines parameters for running an HTTP server. The zero value for Server is a valid configuration.
	// https://golang.org/pkg/net/http/#Server
	s := &http.Server{
		Addr: *bindAddress,
		Handler: sm,
		IdleTimeout: 120 *time.Second,
		ReadTimeout: 1 *time.Second,
		WriteTimeout: 1 *time.Second,
	}

	// Goroutines
	// This new goroutine will execute concurrently with the calling one.
	go func() {
		// start the server
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
	// Block until a signal is received.
	sig := <- sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	// Package context defines the Context type, which carries deadlines, cancellation signals,
	// and other request-scoped values across API boundaries and between processes.
	ctx, _ := context.WithTimeout(context.Background(), 30 *time.Second) // 30 seconds timeout

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	// https://golang.org/pkg/net/http/#Server.Shutdown
	s.Shutdown(ctx)
}
