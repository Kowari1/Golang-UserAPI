package main

import (
	"userapi/internal/config"
	connect "userapi/internal/db"
	"userapi/internal/handler"
	"userapi/internal/kafka"
	"userapi/internal/logger"
	"userapi/internal/middleware"
	"userapi/internal/redisdb"
	"userapi/internal/repository"
	"userapi/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	logger.Log.Info("Loading environment variables")
	if err := config.LoadEnv(); err != nil {
		logger.Log.Fatal("Failed to load .env file", zap.Error(err))
	}

	jwtKey, err := config.GetJwtKey()
	if err != nil {
		logger.Log.Fatal("JWT_KEY error", zap.Error(err))
	}

	dsn, err := config.GetDBDsn()
	if err != nil {
		logger.Log.Fatal("DB_DSN error", zap.Error(err))
	}

	port := config.GetPort()

	logger.Log.Info("Connecting to database")
	db := connect.InitDB(dsn)

	redisClient := redisdb.InitRedisClient()

	broker, err := config.GetKafkaBroker()
	if err != nil {
		logger.Log.Fatal("broker error", zap.Error(err))
	}

	topic, err := config.GetKafkaTopic()
	if err != nil {
		logger.Log.Fatal("topic error", zap.Error(err))
	}

	kafkaProducer := kafka.NewProducer(
		broker,
		topic,
	)

	defer kafkaProducer.Close()

	repo := repository.NewUserRepository(db)
	validator := service.NewValidator(repo)
	redisService := service.NewRedisClient(redisClient)
	userService := service.NewUserService(repo, redisService, jwtKey)

	if err := userService.EnsureDefaultAdmin(); err != nil {
		logger.Log.Fatal("Failed to create default admin", zap.Error(err))
	}

	handler := handler.NewUserHandler(userService, validator, redisService, kafkaProducer)

	r := gin.New()
	r.Use(gin.Logger(), middleware.ErrorRecovery())

	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.Login)
	r.POST("/logout", handler.Logout)

	authUser := r.Group("/")
	authUser.Use(middleware.JWTMiddleware(jwtKey, redisService))
	{
		authUser.PUT("/users/:login", handler.UpdateProfile)
	}

	authAdmin := r.Group("/admin")
	authAdmin.Use(
		middleware.JWTMiddleware(jwtKey, redisService),
		middleware.RequireAdmin(),
	)
	{
		authAdmin.POST("/register", handler.RegisterAdmin)
		authAdmin.GET("/users", handler.GetAll)
		authAdmin.GET("/users/:login", handler.GetByLogin)
		authAdmin.PUT("/users/:login", handler.Update)
		authAdmin.DELETE("/users/:id", handler.Delete)
	}

	logger.Log.Info("Starting HTTP server", zap.String("port", port))

	if err := r.Run(port); err != nil {
		logger.Log.Fatal("Failed to run server", zap.Error(err))
	}
}
