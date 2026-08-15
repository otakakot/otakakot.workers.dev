package main

import (
	"net/http"

	"github.com/syumai/workers"
)

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
		msg := "Hello!"
		w.Write([]byte(msg))
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
