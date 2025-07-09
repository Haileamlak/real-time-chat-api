package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/haileamlak/chat-system/controllers"
)

func SetupRouter(userController controllers.UserController, messageController controllers.MessageController, webSocketController controllers.WebSocketController, authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	router.POST("/signup", userController.SignUp)
	router.POST("/login", userController.Login)

	// Protected routes
	auth := router.Group("/")
	auth.Use(authMiddleware)

	dm := auth.Group("/dm")
	{
		dm.POST("/send", messageController.SendDM)
		dm.GET("/:user", messageController.GetDMHistory)
	}

	group := auth.Group("/group")
	{
		group.POST("/create", messageController.CreateGroup)
		group.POST("/join", messageController.JoinGroup)
		group.POST("/send", messageController.SendGroupMessage)
		group.GET("/:name/history", messageController.GetGroupHistory)
	}
	broadcast := auth.Group("/broadcast")
	{

		broadcast.POST("/send", messageController.SendBroadcast)
		broadcast.GET("/history", messageController.GetBroadcastHistory)

	}

	router.GET("/ws", webSocketController.WebSocketHandler)

	return router
}
