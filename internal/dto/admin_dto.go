package dto

import (
	"userapi/internal/model"

	"github.com/google/uuid"
)

type AdminRegisterRequest struct {
	RegisterRequest
	Admin bool
}

func (r AdminRegisterRequest) ToUserModel() (model.User, error) {
	return model.User{
		ID:       uuid.New(),
		Login:    r.Login,
		Password: r.Password,
		Name:     r.Name,
		Gender:   r.Gender,
		Birthday: r.Birthday,
		Admin:    r.Admin,
	}, nil
}

func (r AdminRegisterRequest) GetLogin() string {
	return r.Login
}
