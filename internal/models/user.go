package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id,omitempty" validate:"uuid_rfc4122"`
	Name     string    `json:"name,omitempty" validate:"required,alphanum"`
	Email    string    `json:"email,omitempty" validate:"required,email"`
	Password string    `json:"password,omitempty" validate:"required,min=6,max=64,printascii"`
}
