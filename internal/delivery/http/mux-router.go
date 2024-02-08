package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter struct{}

var muxDispatcher = mux.NewRouter()

// GET implements Router.
func (*muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("GET")
}

// POST implements Router.
func (*muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}

// SERVE implements Router.
func (*muxRouter) SERVE(port string) {
	fmt.Printf("Mux HTTP server running on port %v", port)
	log.Fatalln(http.ListenAndServe(port, muxDispatcher))
}

func NewMuxRouter() Router {
	return &muxRouter{}
}
