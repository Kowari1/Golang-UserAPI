package contract

import "userapi/internal/model"

type IUserModelConvert interface {
	ToUserModel() (model.User, error)
	GetLogin() string
}
