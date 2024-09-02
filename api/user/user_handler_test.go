package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndGetUser(t *testing.T) {
	t.Parallel()

	userRepository := NewUserRepository(dbInstance)
	userHandler := UserHandler{UserService: NewUserService(userRepository, NewRoleRepository(dbInstance))}

	t.Run("should not save user as request body is missing", func(t *testing.T) {
		t.Parallel()

		req, _ := http.NewRequest("POST", route, nil)
		rec := httptest.NewRecorder()
		http.HandlerFunc(userHandler.createUser).ServeHTTP(rec, req)

		if s := rec.Code; s != http.StatusBadRequest {
			t.Errorf("given %v, expect %v", s, http.StatusBadRequest)
		}
	})

	t.Run("should save user", func(t *testing.T) {
		t.Parallel()

		user := UserDto{
			Fullname: "Franko",
			Email:    "Franko@example.com",
			Role:     USER,
		}

		body, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("failed to marshal user: %v", err)
		}

		req, err := http.NewRequestWithContext(context.Background(), "POST", route, bytes.NewBuffer(body))

		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()
		http.HandlerFunc(userHandler.createUser).ServeHTTP(rec, req)

		if s := rec.Code; s != http.StatusCreated {
			t.Errorf("given %v, expect %v", s, http.StatusCreated)
			return
		}
	})

	t.Run("should test db rollback behaviour when a bottom level query returns error", func(t *testing.T) {
		t.Parallel()

		user := UserDto{
			Fullname: "Rollback",
			Email:    "rollback@demo.com",
			Role:     TEST,
		}

		body, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("failed to marshal user: %v", err)
		}

		req, err := http.NewRequestWithContext(context.Background(), "POST", route, bytes.NewBuffer(body))

		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()
		http.HandlerFunc(userHandler.createUser).ServeHTTP(rec, req)

		if s := rec.Code; s != http.StatusBadRequest {
			t.Errorf("given %v, expect %v", s, http.StatusBadRequest)
		}

		_, err = userRepository.UserByEmail(req.Context(), user.Email)
		if err == nil {
			t.Errorf("top level rollback committed though bottom level context encountered an error")
		}
	})

	t.Run("should retrieve user by email", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		// pre-save user
		if _, err := userRepository.Save(ctx, User{Fullname: "John", Email: "john@email.com"}); err != nil {
			t.Errorf("pre-save failed %s", err)
			return
		}

		url := fmt.Sprintf("%s?email=%s", route, "john@email.com")
		req, _ := http.NewRequest("GET", url, nil)
		rec := httptest.NewRecorder()
		http.HandlerFunc(userHandler.userByEmail).ServeHTTP(rec, req)

		if s := rec.Code; s != http.StatusOK {
			t.Errorf("given %v, expect %v", s, http.StatusOK)
		}
	})
}
