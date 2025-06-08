package userModel

import (
	"time"
)

type UserDTO struct {
	ID uint

	// Authentication
	Username      string
	Email         string
	EmailVerified bool
	LastLogin     *time.Time

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
