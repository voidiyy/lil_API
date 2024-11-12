package types

import (
	"github.com/google/uuid"
	"time"
)

type Worker struct {
	ID        uuid.UUID `psql:"id" json:"id"`
	Name      string    `psql:"name" json:"name"`
	Password  string    `psql:"password" json:"password"`
	Email     string    `psql:"email" json:"email"`
	Role      string    `psql:"role" json:"role"`
	IsActive  bool      `psql:"isactive" json:"isActive"`
	CreatedAt time.Time `json:"-"`
}
