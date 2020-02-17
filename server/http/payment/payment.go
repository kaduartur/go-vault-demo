package payment

import (
	"net/http"

	"github.com/kaduartur/go-vault-demo/server/vault"
)

type Handler struct {
}

func NewHandler(v *vault.Manager) *Handler {
	return &Handler{}
}

func (h *Handler) HandlePay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			processPost(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func processPost(w http.ResponseWriter, r *http.Request) {

}
