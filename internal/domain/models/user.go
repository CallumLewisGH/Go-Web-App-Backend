package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	// Authentication
	Username      string `gorm:"size:50;uniqueIndex;not null"`
	Email         string `gorm:"size:255;uniqueIndex;not null"`
	EmailVerified bool   `gorm:"default:false"`
	PasswordHash  string `gorm:"size:255;not null" json:"-"`
	LastLogin     *time.Time

	// Profile
	ProfilePicture *string `gorm:"text"`
	Bio            string  `gorm:"size:500"`

	// Preferences
	Locale   string `gorm:"size:10;default:'en'"`
	Timezone string `gorm:"size:50"`

	// Status
	IsActive      bool `gorm:"default:true"`
	IsBanned      bool `gorm:"default:false"`
	DeactivatedAt *time.Time

	// Relationships
	// Roles         []Role     `gorm:"many2many:user_roles;"`
	// Settings 	 []Setting  `gorm:"many2many:user_settings;"`
}
