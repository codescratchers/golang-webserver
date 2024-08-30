package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndGetUser(t *testing.T) {
	t.Parallel()

	repo := NewUserRepository(dbInstance)
	userHandler := UserHandler{UserService: NewUserService(repo)}

	t.Run("should not save user as fullname and email are missing", func(t *testing.T) {
		t.Parallel()

		url := fmt.Sprintf("%s?fullname=%s&email=%s", route, "", "")

		req, _ := http.NewRequest("POST", url, nil)
		rec := httptest.NewRecorder()
		http.HandlerFunc(userHandler.createUser).ServeHTTP(rec, req)

		if s := rec.Code; s != http.StatusBadRequest {
			t.Errorf("given %v, expect %v", s, http.StatusBadRequest)
		}
	})

	t.Run("should save user", func(t *testing.T) {
		t.Parallel()

		url := fmt.Sprintf("%s?fullname=%s&email=%s", route, "frank", "frank@email.com")

		req, err := http.NewRequestWithContext(context.Background(), "POST", url, nil)

		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()
		http.HandlerFunc(userHandler.createUser).ServeHTTP(rec, req)

		if s := rec.Code; s != http.StatusCreated {
			t.Errorf("given %v, expect %v", s, http.StatusCreated)
			return
		}

		var dto UserDto

		if b, err := io.ReadAll(rec.Body); err != nil {
			t.Errorf("error reading rec body")
			return
		} else if err = json.Unmarshal(b, &dto); err != nil {
			t.Errorf("error transforming body to UserDto")
			return
		}

		if len(dto.Email) < 1 || len(dto.Fullname) < 1 {
			t.Errorf("dto properties should not be empty")
		}
	})

	t.Run("should retrieve user by email", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		// pre-save user
		if _, err := repo.Save(ctx, User{Fullname: "John", Email: "john@email.com"}); err != nil {
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
