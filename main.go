package main

import (
	"database/sql"
	"net/http"

	"github.com/syumai/workers"
	_ "github.com/syumai/workers/cloudflare/d1"
)

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
		db, err := sql.Open("d1", "DB")
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer db.Close()

		var one int
		if err := db.QueryRowContext(req.Context(), "SELECT 1").Scan(&one); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Write([]byte("health"))
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
