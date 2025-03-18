package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID              uint            `gorm:"primaryKey;not null" json:"-"`
	Uuid            uuid.UUID       `gorm:"index:users_by_uuid;not null" json:"uuid"`
	FirstName       string          `gorm:"not null;type:varchar(191)" json:"first_name"`
	LastName        string          `gorm:"not null;type:varchar(191)" json:"last_name"`
	Email           string          `gorm:"unique;not null" json:"email"`
	EmailVerifiedAt *time.Time      `json:"email_verified_at"`
	AuthProvider    uint            `gorm:"not null" json:"-"`
	Password        string          `gorm:"not null;type:varchar(191)" json:"-"`
	PasswordReset   *datatypes.JSON `json:"-"`
	CreatedAt       time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"not null" json:"-"`
	DeletedAt       gorm.DeletedAt  `json:"-"`
}

func (user *User) BeforeCreate(transaction *gorm.DB) (err error) {
	user.Uuid = uuid.New()
	return
}
