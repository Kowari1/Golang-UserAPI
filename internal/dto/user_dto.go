package dto

import (
	"time"
	"userapi/internal/model"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Login    string     `gorm:"unique;not null" validate:"required,alphanum,min=3,max=20"`
	Password string     `json:"password" validate:"required,min=8,max=20"`
	Name     string     `json:"name" validate:"required"`
	Gender   int        `json:"gender" validate:"oneof=0 1 2"`
	Birthday *time.Time `json:"birthday,omitempty"`
}

func (r RegisterRequest) ToUserModel() (model.User, error) {
	return model.User{
		ID:       uuid.New(),
		Login:    r.Login,
		Password: r.Password,
		Name:     r.Name,
		Gender:   r.Gender,
		Birthday: r.Birthday,
	}, nil
}

func (r RegisterRequest) GetLogin() string {
	return r.Login
}
