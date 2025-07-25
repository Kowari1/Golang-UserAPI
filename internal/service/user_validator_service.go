package service

import (
	"fmt"

	"userapi/internal/contract"
	"userapi/internal/dto"
	"userapi/internal/repository"

	"github.com/go-playground/validator/v10"
)

type UserValidator struct {
	repo repository.UserRepository
}

var validate = validator.New()

func NewValidator(repo repository.UserRepository) *UserValidator {
	return &UserValidator{repo: repo}
}

func (v *UserValidator) ValidateStruct(dto contract.IUserModelConvert) map[string]string {
	errors := make(map[string]string)

	if err := validate.Struct(dto); err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			errors[field] = buildMessage(field, tag, param)
		}
	}

	if exists, _ := v.repo.ExistsByLogin(dto.GetLogin()); exists {
		errors["login"] = "Login already taken"
	}

	return errors
}

func buildMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and digits", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func (v *UserValidator) ValidateLoginRequest(req *dto.LoginRequest) map[string]string {
	errors := map[string]string{}

	if req.Login == "" {
		errors["login"] = "Login is required"
	}

	if req.Password == "" {
		errors["password"] = "Password is required"
	}

	return errors
}
