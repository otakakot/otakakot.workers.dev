package handler

const sessionKVBinding = "KV"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}
