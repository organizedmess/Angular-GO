package routes

import (
	"log"
	"net/http"
	"os"
	"time"

	"url-shortener/backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *handlers.URLHandler) *gin.Engine {
	router := gin.Default()

	origin := os.Getenv("FRONTEND_ORIGIN")
	if origin == "" {
		origin = "http://localhost:4200"
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Printf("[HTTP] %s %s -> %d (%s)", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	})

	router.GET("/api/health", func(c *gin.Context) {
		log.Println("[API] health check requested")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/api/shorten", handler.ShortenURL)
	router.GET("/api/stats/:code", handler.GetStats)
	router.GET("/:code", handler.Redirect)

	return router
}
