package userModel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	//Standard
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Authentication
	Username      string `gorm:"size:50;uniqueIndex;not null"`
	Email         string `gorm:"size:255;uniqueIndex;not null"`
	EmailVerified bool   `gorm:"default:false"`
	PasswordHash  string `gorm:"size:255;not null" json:"-"`
	LastLogin     *time.Time
	AuthId        string `gorm:"size:255;uniqueIndex;not null" json:"-"`

	// Profile
	ProfilePicture *string `gorm:"text"`
	Bio            string  `gorm:"size:500"`

	// Preferences
	Timezone string `gorm:"size:50"`

	// Status
	IsActive      bool `gorm:"default:true"`
	IsBanned      bool `gorm:"default:false"`
	DeactivatedAt *time.Time

	// Relationships
	// Roles         []Role     `gorm:"many2many:user_roles;"`
	// Settings 	 []Setting  `gorm:"many2many:user_settings;"`
}
