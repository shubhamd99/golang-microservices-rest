package handlers

import (
	"log"
	"net/http"
)

// A structure or struct in Golang is a user-defined type that allows to group/combine items of possibly different types into a single type.
type GoodBye struct {
	l *log.Logger
}

// First letter Capital to export
// * - A pointer is a variable whose value is the address of another variable, i.e., direct address of the memory location.
func NewGoodbye(l*log.Logger) *GoodBye {
	// & - For getting the value of the memory address (pointer)
	return &GoodBye{l}
}

// Interfaces are named collections of method signatures.
// Implement this interface on GoodBye type.
func (g*GoodBye) ServeHTTP(rw http.ResponseWriter, r *http.Request) ()  {
	rw.Write([]byte("Bye"))
}
