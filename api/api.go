package api

import (
	"github.com/codescratchers/golang-webserver/api/user"
	"github.com/codescratchers/golang-webserver/database"
	"log"
	"net/http"
)

type ApiObj struct {
	addr    string
	Storage database.Storage
	mux     *http.ServeMux
}

func NewApiServer(addr string, store database.Storage, mux *http.ServeMux) *ApiObj {
	return &ApiObj{addr: addr, Storage: store, mux: mux}
}

func Serve(a *ApiObj) {
	// register user handler -> service -> repository
	user.RegisterUserRoutes(a.mux, user.UserHandler{
		UserService: user.NewUserService(
			a.Storage.DB,
			user.NewUserRepository(),
			user.NewRoleRepository(),
		),
	})

	// start the server
	log.Println("starting server on ", a.addr)

	if err := http.ListenAndServe(a.addr, a.mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
