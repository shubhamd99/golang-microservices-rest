package main

import (
	"GoMicroservices/handlers"
	"log"
	"net/http"
	"os"
)

func main() {

	// Logger
	l := log.New(os.Stdout, "product-api ", log.LstdFlags) // Ex. product-api 2021/05/26 22:44:43 Hello World
	// Handler
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// HTTP request router (or multiplexer)
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	http.ListenAndServe(":9000", sm) // address, server handler
}
