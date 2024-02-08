package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type chaiRouter struct{}

var chiDispatcher = chi.NewRouter()

// GET implements Router.
func (*chaiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(uri, f)
}

// POST implements Router.
func (*chaiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(uri, f)
}

// SERVE implements Router.
func (*chaiRouter) SERVE(port string) {
	fmt.Printf("Chi HTTP server running on port %v", port)
	log.Fatalln(http.ListenAndServe(port, chiDispatcher))
}

func NewChiRouter() Router {
	return &chaiRouter{}
}
