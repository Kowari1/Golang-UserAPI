package dto

type LoginRequest struct {
	Login    string `json:"login" validate:"required,alphanum,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}
