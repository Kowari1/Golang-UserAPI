package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	Login      string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	Name       string    `gorm:"not null"`
	Gender     int       `gorm:"not null"`
	Birthday   *time.Time
	Admin      bool      `gorm:"not null"`
	CreatedOn  time.Time `gorm:"autoCreateTime"`
	CreatedBy  string
	ModifiedOn time.Time `gorm:"autoUpdateTime"`
	ModifiedBy string
	RevokedOn  *time.Time
	RevokedBy  *string
}
