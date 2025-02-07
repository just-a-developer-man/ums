package models

import (
	"time"

	"github.com/google/uuid"
)

// Role is a custom type for user roles.
type Role string

// Predefined roles
const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents a user in the system.
type User struct {
	// ID is the unique identifier of the user. It uses UUID format.
	ID uuid.UUID

	// Name is the user's name.
	Name string

	// Email is the user's email address.
	Email string

	// HashedPassword is the hashed version of the user's password.
	HashedPassword string

	// Role defines the user's role in the system (e.g., "admin" or "user").
	Role Role

	// CreatedAt is the timestamp when the user was created.
	CreatedAt time.Time

	// UpdatedAt is the timestamp when the user was last updated.
	UpdatedAt time.Time
}
