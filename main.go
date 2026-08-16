package main

import (
	"net/http"

	"github.com/syumai/workers"
	_ "github.com/syumai/workers/cloudflare/d1"

	"github.com/otakakot/otakakot.workers.dev/internal/handler"
)

func main() {
	http.HandleFunc("GET /", handler.Health)
	workers.Serve(nil) // use http.DefaultServeMux
}
