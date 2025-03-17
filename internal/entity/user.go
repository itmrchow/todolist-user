package entity

import (
	"crypto/sha512"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `gorm:"embedded"`
	ID         uuid.UUID `gorm:"primaryKey"`
	Name       string    `gorm:"size:20;not null"`
	Email      string    `gorm:"uniqueIndex;size:255;not null"`
	Password   string    `gorm:"size:255;not null"`
}

func (u *User) HashPassword() {

	passwordBytes := sha512.Sum512([]byte(u.Password))
	password := fmt.Sprintf("%x", passwordBytes)

	u.Password = password
}
