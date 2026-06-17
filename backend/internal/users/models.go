package users

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleSuperAdmin      Role = "super_admin"
	RoleModerator       Role = "moderator"
	RoleGraphicDesigner Role = "graphic_designer"
	RolePublisher       Role = "publisher"
	RoleContributor     Role = "contributor"
)

func (r Role) Valid() bool {
	switch r {
	case RoleSuperAdmin, RoleModerator, RoleGraphicDesigner, RolePublisher, RoleContributor:
		return true
	}
	return false
}

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type User struct {
	ID            uuid.UUID
	Name          string
	Email         string
	PasswordHash  string
	Role          Role
	WalletAddress *string
	Bio           *string
	Status        Status
	InvitedBy     *uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
