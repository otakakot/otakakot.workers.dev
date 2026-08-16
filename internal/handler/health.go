package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/syumai/workers/cloudflare/kv"
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

	ns, err := kv.NewNamespace("KV")
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	key := "health"
	value := "health"

	if err := ns.PutString(key, value, &kv.PutOptions{ExpirationTTL: 60}); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	got, err := ns.GetString(key, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if got != value {
		http.Error(w, fmt.Sprintf("kv value mismatch: got %q, want %q", got, value), http.StatusServiceUnavailable)
		return
	}

	if err := ns.Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Write([]byte("health"))
}
