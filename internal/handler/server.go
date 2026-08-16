package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/syumai/workers/cloudflare/kv"

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

const sessionKVBinding = "KV"

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Health(ctx context.Context, request api.HealthRequestObject) (api.HealthResponseObject, error) {
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

func (s *Server) Signin(ctx context.Context, request api.SigninRequestObject) (api.SigninResponseObject, error) {
	body := request.Body
	if body == nil || body.Email == "" || body.Password == "" {
		return api.Signin400JSONResponse{Message: "email and password are required"}, nil
	}

	mock, ok := mockUsers[body.Email]
	if !ok || mock.password != body.Password {
		return api.Signin401JSONResponse{Message: "invalid email or password"}, nil
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return api.SignindefaultJSONResponse{Body: api.ErrorResponse{Message: "failed to create session"}, StatusCode: http.StatusInternalServerError}, nil
	}

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

func (s *Server) Me(ctx context.Context, request api.MeRequestObject) (api.MeResponseObject, error) {
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

func (s *Server) Signout(ctx context.Context, request api.SignoutRequestObject) (api.SignoutResponseObject, error) {
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

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
