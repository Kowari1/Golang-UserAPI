package connect

import (
	"userapi/internal/model"

	"userapi/internal/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Log.Error("failed to connect to DB", zap.Error(err))
		panic(err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Log.Error("Migration error: %v", zap.Error(err))
	}

	logger.Log.Info("connected to DB", zap.String("dsn", dsn))

	return db
}
