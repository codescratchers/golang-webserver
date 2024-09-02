package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const route = "/api/v1/user"

func RegisterUserRoutes(mux *http.ServeMux, handler UserHandler) {
	mux.HandleFunc(fmt.Sprintf("POST %s", route), handler.createUser)
	mux.HandleFunc(fmt.Sprintf("GET %s", route), handler.userByEmail)
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

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		}
	}(r.Body)

	var dto UserDto

	if err = json.Unmarshal(body, &dto); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if len(dto.Fullname) < 1 || len(dto.Email) < 1 || len(dto.Role) < 1 {
		constructErrorResponse(
			w,
			ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
			},
		)
		return
	}

	// if email exist duplicate else if an error occurs save else successfully saved
	if _, err := h.UserService.UserByEmail(r.Context(), dto.Email); err == nil {
		constructErrorResponse(
			w,
			ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "duplicate email",
			},
		)
		return
	}

	if _, err := h.UserService.CreateUser(r.Context(), dto); err != nil {
		constructErrorResponse(
			w,
			ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("%s", err),
			},
		)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *UserHandler) userByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if len(email) < 1 {
		constructErrorResponse(
			w,
			ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "please an email",
			},
		)
		return
	}

	if user, err := h.UserService.UserByEmail(r.Context(), email); err != nil {
		constructErrorResponse(
			w,
			ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "user not found",
			},
		)
	} else {
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(UserDto{Fullname: user.Fullname, Email: user.Email})
		_, _ = w.Write(response)
	}
}
