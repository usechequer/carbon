package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID              uint      `gorm:"primaryKey;not null"`
	Uuid            uuid.UUID `gorm:"index:users_by_uuid;not null"`
	FirstName       string    `gorm:"not null;type:varchar(191)"`
	LastName        string    `gorm:"not null;type:varchar(191)"`
	Email           string    `gorm:"unique;not null"`
	EmailVerifiedAt *time.Time
	AuthProvider    uint   `gorm:"not null"`
	Password        string `gorm:"not null;type:varchar(191)"`
	PasswordReset   *datatypes.JSON
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
	DeletedAt       gorm.DeletedAt
}

func (user *User) BeforeCreate(transaction *gorm.DB) (err error) {
	user.Uuid = uuid.New()
	return
}
