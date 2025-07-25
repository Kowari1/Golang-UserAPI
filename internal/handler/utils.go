package handler

import (
	"net/http"
	"time"
	"userapi/internal/contract"
	"userapi/internal/kafka"
	"userapi/internal/logger"
	"userapi/internal/model"
	"userapi/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BindValidateConvert(
	c *gin.Context,
	dtoObj contract.IUserModelConvert,
	validator *service.UserValidator,
) (model.User, bool) {
	if err := c.ShouldBindJSON(dtoObj); err != nil {
		logger.WarnError(c, ErrInvalidJSON, err)

		JSONErrorMsg(c, http.StatusBadRequest, ErrInvalidJSON)

		return model.User{}, false
	}

	if errs := validator.ValidateStruct(dtoObj); len(errs) > 0 {
		logger.WarnFields(c, LogValidationErr, zap.Any("errors", errs))

		JSONError(c, http.StatusBadRequest, errs)

		return model.User{}, false
	}

	user, err := dtoObj.ToUserModel()

	if err != nil {
		logger.WarnError(c, ErrConversionFailed, err)
		JSONErrorMsg(c, http.StatusBadRequest, ErrConversionFailed)
		return model.User{}, false
	}

	return user, true
}

func HandleRegister(
	c *gin.Context,
	dtoObj contract.IUserModelConvert,
	createdBy string,
	validator *service.UserValidator,
	registerFunc func(model.User) error,
	kafkaProducer *kafka.KafkaProducer,
) {
	user, ok := BindValidateConvert(c, dtoObj, validator)

	if !ok {
		return
	}

	user.CreatedBy = createdBy

	if err := registerFunc(user); err != nil {
		logger.WarnError(c, LogValidationErr, err)

		JSONErrorMsg(c, http.StatusInternalServerError, ErrRegistrationFailed)

		return
	}

	event := kafka.UserRegisteredEvent{
		UserID: user.ID.String(),
		Login:  user.Login,
		Time:   time.Now().Format(time.RFC3339),
	}

	go kafkaProducer.SendMessage(c.Request.Context(), event)

	JSONCreated(c, gin.H{"message": MsgUserRegistered})
}

func HandleUpdate(
	c *gin.Context,
	dtoObj contract.IUserModelConvert,
	validator *service.UserValidator,
	updateFunc func(model.User) error,
) {
	user, ok := BindValidateConvert(c, dtoObj, validator)

	if !ok {
		return
	}

	user.ModifiedBy = c.GetString("login")

	if err := updateFunc(user); err != nil {
		logger.WarnError(c, LogUpdateFail, err)

		JSONErrorMsg(c, http.StatusInternalServerError, ErrUpdateFailed)

		return
	}

	JSONOK(c, gin.H{"message": MsgUserUpdated})
}
