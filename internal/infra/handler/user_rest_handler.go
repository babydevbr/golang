// Package handler manage flow
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/babydevbr/golang/internal/domain"
	"github.com/babydevbr/golang/internal/infra/gateway"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type servicer interface {
	SignUp(user *domain.User) (*domain.User, error)
}

// UserRestHandler http handler.
type UserRestHandler struct {
	usecase servicer
	guard   *gateway.Guardian
}

// NewUserRestHandler factory.
func NewUserRestHandler(uc servicer, guard *gateway.Guardian) *UserRestHandler {
	return &UserRestHandler{
		usecase: uc,
		guard:   guard,
	}
}

// SignUp endpoint.
func (uh *UserRestHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	req := &domain.User{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := uh.usecase.SignUp(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

// SignIn generate JWT Token.
func (uh *UserRestHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uh.guard.GenerateToken(w, r)
}

// SetUserRoutes mux configuration.
func (uh *UserRestHandler) SetUserRoutes(r *mux.Router, n negroni.Negroni) {
	r.Handle("/signup", n.With(
		negroni.WrapFunc(uh.SignUp),
	)).Methods(http.MethodPost, http.MethodOptions).Name("signup")

	r.Handle("/signin", n.With(
		negroni.WrapFunc(uh.SignIn),
	)).Methods(http.MethodPost, http.MethodOptions).Name("signin")
}
