package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const route = "/api/v1/user"

func RegisterUserRoutes(mux *http.ServeMux, handler UserHandler) {
	mux.HandleFunc(fmt.Sprintf("POST %s", route), handler.CreateUser)
}

type UserHandler struct {
	UserService IUserService
}

func constructErrorResponse(w http.ResponseWriter, e ErrorResponse) {
	w.WriteHeader(e.Status)
	res, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fullname := r.URL.Query().Get("fullname")

	if len(fullname) < 1 {
		constructErrorResponse(
			w,
			ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "please enter your fullname",
			},
		)
		return
	}

	if user, err := h.UserService.CreateUser(r.Context(), fullname); err != nil {
		constructErrorResponse(
			w,
			ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "please enter your fullname",
			},
		)
	} else {
		w.WriteHeader(http.StatusCreated)
		response, _ := json.Marshal(UserDto{Fullname: user.Fullname})

		_, _ = w.Write(response)
	}
}
