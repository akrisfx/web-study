package main

import (
	"backend/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к базе данных
	pool := ConnectToDB()
	defer pool.Close()
	router := gin.Default()
	
	// Configure CORS properly
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false, // Set to false for '*' origins
		MaxAge:           12 * time.Hour,
	}))
	
	h := &services.Handler{DB: pool}
	
	router.GET("/api/reports", h.NetGetRecyclingReports)
	router.GET("/api/reports/:id", h.NetGetRecyclingReportsByID)
	router.POST("/api/reports", h.NetAddRecyclingReport)
	router.PUT("/api/reports/:id", h.NetUpdateRecyclingReport)
	router.DELETE("/api/reports/:id", h.NetDeleteRecyclingReport)

	// Waste Types
	router.GET("/api/waste-types", h.NetGetWasteTypes)
	router.GET("/api/waste-types/:id", h.NetGetWasteTypeByID)
	router.POST("/api/waste-types", h.NetAddWasteType)
	router.PUT("/api/waste-types/:id", h.NetUpdateWasteType)
	router.DELETE("/api/waste-types/:id", h.NetDelWasteType)

	// Collection Points
	router.GET("/api/collection-points", h.NetGetCollectionPoints)
	router.GET("/api/collection-points/:id", h.NetGetCollectionPointByID)
	router.POST("/api/collection-points", h.NetAddCollectionPoint)
	router.PUT("/api/collection-points/:id", h.NetUpdateCollectionPoint)
	router.DELETE("/api/collection-points/:id", h.NetDelCollectionPoint)

	// Users
	router.GET("/api/users", h.NetGetUsers)
	router.GET("/api/users/:id", h.NetGetUserByID)
	router.POST("/api/users", h.NetAddUser)
	router.PUT("/api/users/:id", h.NetUpdateUser)
	router.DELETE("/api/users/:id", h.NetDelUser)

	router.Run("localhost:1440")
}
