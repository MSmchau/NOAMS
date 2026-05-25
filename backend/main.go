package main

import (
	"context"
	"log"
	"os"
	"time"

	"noams/config"
	"noams/middleware"
	"noams/models"
	"noams/routes"
	"noams/services"
	"noams/utils"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	config.InitConfig()

	initLogger()

	db := initDB()

	initJWT()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(config.AppConfig.Server.CorsOrigins))

	pinger := startPinger(db)
	routes.Setup(r, db, pinger)

	port := config.AppConfig.Server.Port
	zap.L().Info("NOAMS server starting", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}

func initLogger() {
	level := zap.InfoLevel
	if config.AppConfig.Log.Level == "debug" {
		level = zap.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	if config.AppConfig.Log.Filename != "" {
		logFile, err := os.OpenFile(config.AppConfig.Log.Filename,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			fileWriter := zapcore.AddSync(logFile)
			fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:      "time",
				LevelKey:     "level",
				MessageKey:   "msg",
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeLevel:  zapcore.LowercaseLevelEncoder,
			})
			core = zapcore.NewCore(fileEncoder, fileWriter, level)
		}
	}

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

func initDB() *gorm.DB {
	dsn := config.AppConfig.DSN()
	driver := config.AppConfig.DBDriver()
	zap.L().Info("connecting to database",
		zap.String("driver", driver),
		zap.String("dsn", dsn))

	var dialector gorm.Dialector
	if driver == "sqlite" {
		// Ensure data directory exists
		os.MkdirAll("data", 0755)
		dialector = sqlite.Open(dsn)
	} else {
		dialector = mysql.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.DeviceGroup{},
		&models.Credential{},
		&models.Device{},
		&models.InspectionResult{},
		&models.ConfigBackup{},
		&models.Alert{},
		&models.ScheduledTask{},
	); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	zap.L().Info("database migration completed")

	// Initialize default admin user if not exists
	initDefaultAdmin(db)

	return db
}

func initDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hash, err := utils.HashPassword("admin123")
		if err != nil {
			zap.L().Error("failed to hash default password", zap.Error(err))
			return
		}
		admin := models.User{
			Username: "admin",
			Password: hash,
			Nickname: "系统管理员",
			Role:     "admin",
			Status:   1,
		}
		if err := db.Create(&admin).Error; err != nil {
			zap.L().Error("failed to create default admin", zap.Error(err))
		} else {
			zap.L().Info("default admin user created (username: admin, password: admin123)")
		}
	}
}

func initJWT() {
	secret := config.AppConfig.JWT.Secret
	if secret == "" {
		secret = "noams-default-jwt-secret-change-in-production"
		zap.L().Warn("using default JWT secret, please set NOAMS_JWT_SECRET in production")
	}
	utils.JWTSecret = []byte(secret)
}

func startPinger(db *gorm.DB) *services.Pinger {
	pinger := services.NewPinger(db, services.PingerConfig{
		Interval:    time.Duration(config.AppConfig.Monitor.PingInterval) * time.Second,
		Timeout:     time.Duration(config.AppConfig.Monitor.PingTimeout) * time.Second,
		Retry:       config.AppConfig.Monitor.PingRetry,
		Method:      config.AppConfig.Monitor.PingMethod,
		Concurrency: config.AppConfig.Monitor.Concurrency,
	})
	ctx := context.Background()
	go pinger.Start(ctx)
	return pinger
}

func init() {
	os.MkdirAll("logs", 0755)
}
