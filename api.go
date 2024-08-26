package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type Storage struct {
	statuses map[string]int
	users    map[string]*User
	nextID   uint64
	mu       *sync.RWMutex
}

type ProfileParams struct {
	Login string `apivalidator:"required"`
}

type CreateParams struct {
	Login  string `apivalidator:"required,min=10"`
	Name   string `apivalidator:"paramname=full_name"`
	Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	Age    int    `apivalidator:"min=0,max=128"`
}

type User struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Status   int    `json:"status"`
}

type NewUser struct {
	ID uint64 `json:"id"`
}

// apigen:api {"url": "/user/profile", "auth": false}
func (storage *Storage) Profile(ctx context.Context, in ProfileParams) (*User, error) {
	if in.Login == "bad_user" {
		return nil, fmt.Errorf("bad user")
	}

	storage.mu.RLock()
	user, exist := storage.users[in.Login]
	storage.mu.RUnlock()
	if !exist {
		return nil, ApiError{http.StatusNotFound, fmt.Errorf("user not exist")}
	}

	return user, nil
}

// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
func (storage *Storage) Create(ctx context.Context, in CreateParams) (*NewUser, error) {

	if in.Login == "bad_username" {
		return nil, fmt.Errorf("bad user")
	}

	storage.mu.Lock()
	defer storage.mu.Unlock()

	_, exist := storage.users[in.Login]
	if exist {
		return nil, ApiError{http.StatusConflict, fmt.Errorf("user %s exist", in.Login)}
	}

	id := storage.nextID
	storage.nextID++
	storage.users[in.Login] = &User{
		ID:       id,
		Login:    in.Login,
		FullName: in.Name,
		Status:   storage.statuses[in.Status],
	}

	return &NewUser{id}, nil
}
