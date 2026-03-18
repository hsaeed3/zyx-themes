package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusPending  Status = "pending"
)

type User struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Status    Status            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Meta      map[string]string `json:"meta,omitempty"`
}

type UserService interface {
	GetUser(ctx context.Context, id int) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, name, email string) (*User, error)
}

type userService struct {
	store map[int]*User
	seq   int
}

func NewUserService() UserService {
	return &userService{store: make(map[int]*User)}
}

func (s *userService) GetUser(_ context.Context, id int) (*User, error) {
	u, ok := s.store[id]
	if !ok {
		return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)
	}
	return u, nil
}

func (s *userService) ListUsers(_ context.Context) ([]User, error) {
	users := make([]User, 0, len(s.store))
	for _, u := range s.store {
		users = append(users, *u)
	}
	return users, nil
}

func (s *userService) CreateUser(_ context.Context, name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}
	s.seq++
	u := &User{
		ID:        s.seq,
		Name:      name,
		Email:     email,
		Status:    StatusPending,
		CreatedAt: time.Now().UTC(),
	}
	s.store[u.ID] = u
	return u, nil
}

var ErrNotFound = errors.New("not found")

type Handler struct {
	svc UserService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		users, err := h.svc.ListUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var body struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		user, err := h.svc.CreateUser(r.Context(), body.Name, body.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	svc := NewUserService()
	h := &Handler{svc: svc}

	mux := http.NewServeMux()
	mux.Handle("/users", h)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
