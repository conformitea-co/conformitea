package user

import (
	"time"

	"conformitea/domain/organization"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID                   `json:"id"`
	Email         string                      `json:"email"`
	FirstName     string                      `json:"first_name"`
	LastName      string                      `json:"last_name"`
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
	Organizations []organization.Organization `json:"organizations,omitempty"`
}
