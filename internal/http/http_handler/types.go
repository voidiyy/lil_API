package http_handler

import "time"

// RegisterWorkerRequest handle /worker/register
type RegisterWorkerRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

// LoginWorkerRequest handle /worker/login
type LoginWorkerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// ProfileWorkerResponse handle /worker/profile
type ProfileWorkerResponse struct {
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Role      string    `json:"role" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
}
