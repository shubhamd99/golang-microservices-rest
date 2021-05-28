package main

import (
	"GoMicroservices/data"
	"GoMicroservices/handlers"
	"context"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
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
	v := data.NewValidation()

	// Handler
	ph := handlers.NewProducts(l, v)

	// ServeMux is an HTTP request multiplexer.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern that most closely matches the URL.
	// create a new server mux and register the handlers
	sm := mux.NewRouter()

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))
	// ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"})) Allow Every Origin

	// create a new server
	// A Server defines parameters for running an HTTP server. The zero value for Server is a valid configuration.
	// https://golang.org/pkg/net/http/#Server
	s := &http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Goroutines Function to Start the Server
	// This new goroutine will execute concurrently with the calling one.
	go func() {
		// start the server
		err := s.ListenAndServe()
		if err != nil {
			// Fatal is equivalent to Print() followed by a call to os.Exit(1)
			l.Fatal(err)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
