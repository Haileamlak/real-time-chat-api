package main

import (
	"log"
	"os"

	"github.com/haileamlak/chat-system/infrastructure"
	"github.com/haileamlak/chat-system/repositories"
	"github.com/haileamlak/chat-system/usecases"
	"github.com/haileamlak/chat-system/controllers"
	"github.com/haileamlak/chat-system/routers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // fallback for local dev
	}
	redisService := infrastructure.NewRedisService(addr)
	defer redisService.Close()
	passwordService := infrastructure.NewPasswordService()
	tokenService := infrastructure.NewTokenService(redisService)

	authMiddleware := infrastructure.NewAuthMiddleware(tokenService)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(redisService)
	messageRepo := repositories.NewMessageRepository(redisService)

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo, passwordService, tokenService)
	messageUseCase := usecases.NewMessageUseCase(messageRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userUseCase)
	messageController := controllers.NewMessageController(messageUseCase)
	webSocketController := controllers.NewWebSocketController(messageUseCase, redisService)

	router := routers.SetupRouter(userController, messageController, webSocketController, authMiddleware.Authenticate())

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}

	log.Println("Server running on http://localhost:8080")
}
