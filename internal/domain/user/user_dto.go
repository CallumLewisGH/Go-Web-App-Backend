package userModel

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	// Authentication
	Username      string
	Email         string
	EmailVerified bool
	LastLogin     *time.Time
	AuthId        string

	// Profile
	ProfilePicture *string
	Bio            string

	// Preferences
	Locale   string
	Timezone string

	// Status
	IsActive      bool
	IsBanned      bool
	DeactivatedAt *time.Time

	// Relationships
	// Roles         []Role     `gorm:"many2many:user_roles;"`
	// Settings 	 []Setting  `gorm:"many2many:user_settings;"`
}
