package handler

import (
	"context"
	"net/http"
	"uuid"

	"github.com/syumai/workers-go/cloudflare/kv"

	"github.com/otakakot/otakakot.workers.dev/internal/api"
)

type mockUser struct {
	userID   string
	password string
}

var mockUsers = map[string]mockUser{
	"test@example.com": {
		userID:   "018f1234-5678-7abc-8def-0123456789ab",
		password: "11111111",
	},
}

// POST /signin
func (hdl *Handler) Signin(ctx context.Context, request api.SigninRequestObject) (api.SigninResponseObject, error) {
	body := request.Body
	if body == nil || body.Email == "" || body.Password == "" {
		return api.Signin400JSONResponse{Message: "email and password are required"}, nil
	}

	mock, ok := mockUsers[body.Email]
	if !ok || mock.password != body.Password {
		return api.Signin401JSONResponse{Message: "invalid email or password"}, nil
	}

	sessionID := uuid.New().String()

	ns, err := kv.NewNamespace(sessionKVBinding)
	if err != nil {
		return api.SignindefaultJSONResponse{Body: api.ErrorResponse{Message: "failed to access session store"}, StatusCode: http.StatusInternalServerError}, nil
	}

	if err := ns.PutString("session:"+sessionID, mock.userID, &kv.PutOptions{ExpirationTTL: 86400}); err != nil {
		return api.SignindefaultJSONResponse{Body: api.ErrorResponse{Message: "failed to create session"}, StatusCode: http.StatusInternalServerError}, nil
	}

	cookie := (&http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   86400,
	}).String()

	return api.Signin200JSONResponse{
		Body:    api.SigninResponse{UserId: mock.userID},
		Headers: api.Signin200ResponseHeaders{SetCookie: &cookie},
	}, nil
}
