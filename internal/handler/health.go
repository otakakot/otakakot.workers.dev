package handler

import (
	"database/sql"
	"net/http"
)

func Health(w http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("d1", "DB")
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	if _, err := db.ExecContext(req.Context(), "SELECT 1"); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Write([]byte("health"))
}
