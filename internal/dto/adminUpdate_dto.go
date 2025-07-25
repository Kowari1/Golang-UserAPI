package dto

import (
	"userapi/internal/model"

	"github.com/google/uuid"
)

type AdminUpdateRequest struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey"`
	Login    string    `gorm:"unique;not null" validate:"required,alphanum,min=3,max=20"`
	Password string    `json:"password" validate:"required,min=8,max=20"`
	Name     string    `json:"name" validate:"required"`
	Gender   int       `json:"gender" validate:"oneof=0 1 2"`
	Admin    bool
}

func (r AdminUpdateRequest) ToUserModel() (model.User, error) {
	return model.User{
		ID:       r.ID,
		Login:    r.Login,
		Password: r.Password,
		Name:     r.Name,
		Gender:   r.Gender,
		Admin:    r.Admin,
	}, nil
}

func (r AdminUpdateRequest) GetLogin() string {
	return r.Login
}
