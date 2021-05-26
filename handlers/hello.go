package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// A structure or struct in Golang is a user-defined type that allows to group/combine items of possibly different types into a single type.
type Hello struct {
	l *log.Logger
}

// First letter Capital to export
// * - A pointer is a variable whose value is the address of another variable, i.e., direct address of the memory location.
func NewHello(l *log.Logger) *Hello {
	// & - For getting the value of the memory address (pointer)
	return &Hello{l}
}

// Interfaces are named collections of method signatures.
// Implement this interface on Hello type.
func (h*Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		// rw.WriteHeader(http.StatusBadRequest)
		// rw.Write([]byte("Oops"))
		return
	}

	fmt.Fprintf(rw, "Success %s", d) // Return back to the user
}
