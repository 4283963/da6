package main

import (
	"fmt"
	"log"

	"aquarium-control/internal/config"
	"aquarium-control/internal/database"
	"aquarium-control/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func main() {
	if err := config.LoadConfig("./config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.InitDB(&config.AppConfig.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.Default()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		if ctx.Writer.Status() == 200 && ctx.Request.Method == "OPTIONS" {
			ctx.Abort()
			return
		}
		ctx.Next()
	})

	router.SetupRoutes(r)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
