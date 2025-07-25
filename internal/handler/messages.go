package handler

const (
	ErrInvalidJSON        = "invalid JSON"
	ErrConversionFailed   = "conversion failed"
	ErrRegistrationFailed = "registration failed"
	ErrLoginFailed        = "invalid login or password"
	ErrDelete             = "failed to delete user"
	ErrUpdateFailed       = "update failed"
	MsgUserRegistered     = "user registered"
	MsgAdminRegistered    = "admin registered"
	MsgUserDeleted        = "user deleted"
	MsgUserUpdated        = "user updated"
	MsgUserAuthorize      = "authorization successful"
	ErrUUID               = "invalid UUID"

	LogRegisterFail  = "register: service failed"
	LogValidationErr = "validation failed"
	LogUpdateFail    = "update: service failed"
	LogGetAddFail    = "failed to get all users"
	LogGetByLogin    = "failed to get user by login"
)
