package routes

import (
	"noams/handlers"
	"noams/middleware"
	"noams/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, pinger *services.Pinger) {
	authService := services.NewAuthService(db)
	credentialService := services.NewCredentialService(db)
	deviceService := services.NewDeviceService(db)

	authHandler := handlers.NewAuthHandler(authService)
	credentialHandler := handlers.NewCredentialHandler(credentialService)
	deviceHandler := handlers.NewDeviceHandler(deviceService, pinger)
	inspectionHandler := handlers.NewInspectionHandler(db)
	configHandler := handlers.NewConfigHandler(db)
	alertHandler := handlers.NewAlertHandler(db)
	taskHandler := handlers.NewTaskHandler(db)
	monitorHandler := handlers.NewMonitorHandler(db)

	api := r.Group("/api/v1")

	// Public routes
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)

	// Protected routes
	authed := api.Group("")
	authed.Use(middleware.AuthRequired())
	{
		// User
		authed.GET("/auth/profile", authHandler.Profile)

		// Devices
		authed.GET("/devices", deviceHandler.List)
		authed.GET("/devices/all", deviceHandler.ListAll)
		authed.GET("/devices/stats", deviceHandler.Stats)
		authed.GET("/devices/:id", deviceHandler.Get)
		authed.POST("/devices", deviceHandler.Create)
		authed.GET("/devices/:id/ping", deviceHandler.Ping)
		authed.PUT("/devices/:id", deviceHandler.Update)
		authed.DELETE("/devices/:id", middleware.AdminRequired(), deviceHandler.Delete)

		// Credentials
		authed.GET("/credentials", credentialHandler.List)
		authed.GET("/credentials/:id", credentialHandler.Get)
		authed.POST("/credentials", credentialHandler.Create)
		authed.PUT("/credentials/:id", credentialHandler.Update)
		authed.DELETE("/credentials/:id", credentialHandler.Delete)

		// Inspection
		authed.POST("/devices/:id/inspect", inspectionHandler.InspectDevice)
		authed.POST("/inspections/batch", inspectionHandler.BatchInspect)
		authed.GET("/inspections/report", inspectionHandler.Report)
		authed.GET("/inspections/latest", inspectionHandler.LatestReport)
		authed.GET("/inspections/history/:deviceId", inspectionHandler.History)

		// Config backup
		authed.POST("/configs/backup", configHandler.Backup)
		authed.GET("/configs/history/all", configHandler.ListAll)
		authed.GET("/configs/history/:deviceId", configHandler.History)
		authed.POST("/configs/rollback", configHandler.Rollback)
		authed.GET("/configs/diff", configHandler.Diff)

		// Alerts
		authed.GET("/alerts", alertHandler.List)
		authed.GET("/alerts/stats", alertHandler.Stats)
		authed.PUT("/alerts/:id/confirm", alertHandler.Confirm)
		authed.PUT("/alerts/:id/resolve", alertHandler.Resolve)

		// Tasks
		authed.GET("/tasks", taskHandler.List)
		authed.POST("/tasks", taskHandler.Create)
		authed.PUT("/tasks/:id", taskHandler.Update)
		authed.DELETE("/tasks/:id", taskHandler.Delete)
		authed.PUT("/tasks/:id/toggle", taskHandler.Toggle)
		authed.GET("/tasks/logs", taskHandler.Logs)

		// Monitor / Dashboard
		authed.GET("/monitor/dashboard", monitorHandler.Dashboard)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
