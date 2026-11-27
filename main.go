package main

import (
	"net/http"

	"github.com/syumai/workers-go"
	_ "github.com/syumai/workers-go/cloudflare/d1"
	_ "github.com/syumai/workers-go/cloudflare/kv"

	"github.com/otakakot/otakakot.workers.dev/internal/api"
	"github.com/otakakot/otakakot.workers.dev/internal/handler"
	"github.com/otakakot/otakakot.workers.dev/internal/middleware"
)

func main() {
	hdl := handler.NewHandler()
	srv := api.NewStrictHandler(hdl, nil)
	http.Handle("/", middleware.CORS()(api.Handler(srv)))
	workers.Serve(nil) // use http.DefaultServeMux
}
