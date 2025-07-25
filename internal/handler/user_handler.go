package handler

import (
	"net/http"
	"userapi/internal/config"
	"userapi/internal/dto"
	"userapi/internal/logger"
	"userapi/internal/model"
	"userapi/internal/service"

	"time"

	"userapi/internal/kafka"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	service       *service.UserService
	validator     *service.UserValidator
	redisService  *service.RedisService
	kafkaProducer *kafka.KafkaProducer
}

func NewUserHandler(
	service *service.UserService,
	validator *service.UserValidator,
	redisService *service.RedisService,
	kafkaProducer *kafka.KafkaProducer,
) *UserHandler {
	return &UserHandler{
		service:       service,
		validator:     validator,
		redisService:  redisService,
		kafkaProducer: kafkaProducer,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnError(c, ErrInvalidJSON, err)
		JSONErrorMsg(c, http.StatusBadRequest, ErrInvalidJSON)
		return
	}

	if errs := h.validator.ValidateLoginRequest(&req); len(errs) > 0 {
		logger.WarnFields(c, LogValidationErr, zap.Any("errors", errs))
		JSONError(c, http.StatusBadRequest, errs)
		return
	}

	token, err := h.service.Login(req.Login, req.Password)

	if err != nil {
		logger.WarnError(c, ErrLoginFailed, err)
		JSONErrorMsg(c, http.StatusUnauthorized, ErrLoginFailed)
		return
	}

	c.Header("Authorization", "Bearer "+token)

	JSONOK(c, gin.H{"message": MsgUserAuthorize})
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	HandleRegister(
		c,
		&dto.RegisterRequest{},
		"self",
		h.validator,
		h.service.Register,
		h.kafkaProducer,
	)
}

func (h *UserHandler) RegisterAdmin(c *gin.Context) {
	HandleRegister(
		c,
		&dto.AdminRegisterRequest{},
		c.GetString("login"),
		h.validator,
		h.service.Register,
		h.kafkaProducer,
	)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		logger.WarnError(c, ErrUUID, err)

		JSONErrorMsg(c, http.StatusBadRequest, ErrUUID)

		return
	}

	if err := h.service.Delete(id); err != nil {
		logger.WarnError(c, ErrDelete, err)

		JSONErrorMsg(c, http.StatusInternalServerError, ErrDelete)

		return
	}

	JSONOK(c, gin.H{"message": MsgUserDeleted})
}

func (h *UserHandler) Update(c *gin.Context) {
	HandleUpdate(
		c,
		&dto.AdminUpdateRequest{},
		h.validator,
		h.service.Update,
	)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	var users []model.User
	ctx := c.Request.Context()

	users, err := h.service.GetAll(ctx)
	if err != nil {
		logger.WarnError(c, LogGetAddFail, err)

		JSONErrorMsg(c, http.StatusInternalServerError, ErrConversionFailed)

		return
	}

	JSONOK(c, users)
}

func (h *UserHandler) GetByLogin(c *gin.Context) {
	login := c.Param("login")

	user, err := h.service.GetByLogin(login)

	if err != nil {
		logger.WarnError(c, LogGetByLogin, err)

		JSONErrorMsg(c, http.StatusInternalServerError, err.Error())

		return
	}

	JSONOK(c, user)
}

func (h *UserHandler) Logout(c *gin.Context) {
	jti := c.GetString("jti")

	exp := time.Now().Add(config.GetJwtExpiration())
	ttl := time.Until(exp)

	err := h.redisService.SetToBlacklist(c.Request.Context(), jti, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	HandleUpdate(
		c,
		&dto.UpdateRequest{},
		h.validator,
		h.service.Update,
	)
}
