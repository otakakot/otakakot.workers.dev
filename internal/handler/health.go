package handler

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syumai/workers-go/cloudflare/kv"

	"github.com/otakakot/otakakot.workers.dev/internal/api"
)

// GET /
func (hdl *Handler) Health(ctx context.Context, request api.HealthRequestObject) (api.HealthResponseObject, error) {
	db, err := sql.Open("d1", "DB")
	if err != nil {
		return api.Health503TextResponse(err.Error()), nil
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, "SELECT 1"); err != nil {
		return api.Health503TextResponse(err.Error()), nil
	}

	ns, err := kv.NewNamespace(sessionKVBinding)
	if err != nil {
		return api.Health503TextResponse(err.Error()), nil
	}

	const (
		key   = "health"
		value = "health"
	)

	if err := ns.PutString(key, value, &kv.PutOptions{ExpirationTTL: 60}); err != nil {
		return api.Health503TextResponse(err.Error()), nil
	}

	got, err := ns.GetString(key, nil)
	if err != nil {
		return api.Health503TextResponse(err.Error()), nil
	}

	if got != value {
		return api.Health503TextResponse(fmt.Sprintf("kv value mismatch: got %q, want %q", got, value)), nil
	}

	if err := ns.Delete(key); err != nil {
		return api.Health503TextResponse(err.Error()), nil
	}

	return api.Health200Response{}, nil
}
