package handler

import (
	"context"
	"net/http"

	"github.com/syumai/workers-go/cloudflare/kv"

	"github.com/otakakot/otakakot.workers.dev/internal/api"
)

// POST /signout
func (hdl *Handler) Signout(ctx context.Context, request api.SignoutRequestObject) (api.SignoutResponseObject, error) {
	ns, err := kv.NewNamespace(sessionKVBinding)
	if err == nil && request.Params.Session != nil && *request.Params.Session != "" {
		_ = ns.Delete("session:" + *request.Params.Session)
	}

	cleared := (&http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
	}).String()

	return api.Signout204Response{
		Headers: api.Signout204ResponseHeaders{SetCookie: &cleared},
	}, nil
}
