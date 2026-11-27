package handler

import (
	"context"

	"github.com/syumai/workers-go/cloudflare/kv"

	"github.com/otakakot/otakakot.workers.dev/internal/api"
)

// GET /me
func (hdl *Handler) Me(ctx context.Context, request api.MeRequestObject) (api.MeResponseObject, error) {
	if request.Params.Session == nil || *request.Params.Session == "" {
		return api.Me401JSONResponse{Message: "unauthorized"}, nil
	}

	ns, err := kv.NewNamespace(sessionKVBinding)
	if err != nil {
		return api.Me401JSONResponse{Message: "unauthorized"}, nil
	}

	userID, err := ns.GetString("session:"+*request.Params.Session, nil)
	if err != nil || userID == "" || userID == "<null>" {
		return api.Me401JSONResponse{Message: "unauthorized"}, nil
	}

	return api.Me200JSONResponse{UserId: userID}, nil
}
