package main

import (
	"net/http"

	"github.com/syumai/workers"
	_ "github.com/syumai/workers/cloudflare/d1"
	_ "github.com/syumai/workers/cloudflare/kv"

	"github.com/otakakot/otakakot.workers.dev/internal/api"
	"github.com/otakakot/otakakot.workers.dev/internal/handler"
	"github.com/otakakot/otakakot.workers.dev/internal/middleware"
)

func main() {
	server := handler.NewServer()
	handler := api.NewStrictHandler(server, nil)
	http.Handle("/", middleware.CORS()(api.Handler(handler)))
	workers.Serve(nil) // use http.DefaultServeMux
}
