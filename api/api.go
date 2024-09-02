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

func Serve(s *ApiObj) {
	// register user handler -> service -> repository
	user.RegisterUserRoutes(s.mux, user.UserHandler{
		UserService: user.NewUserService(
			user.NewUserRepository(s.Storage.DB),
			user.NewRoleRepository(s.Storage.DB),
		),
	})

	// start the server
	log.Println("starting server on ", s.addr)

	if err := http.ListenAndServe(s.addr, s.mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
