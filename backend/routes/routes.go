package routes

import (
	"whatsapp-clone/controllers"
	"whatsapp-clone/middleware"
	"whatsapp-clone/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", controllers.GetCurrentUser)
				users.PUT("/me", controllers.UpdateCurrentUser)
				users.GET("/search", controllers.SearchUsers)
				users.PUT("/status", controllers.UpdateOnlineStatus)
			}

			// Conversation routes
			conversations := protected.Group("/conversations")
			{
				conversations.POST("", controllers.CreateConversation)
				conversations.GET("", controllers.GetConversations)
				conversations.GET("/:id", controllers.GetConversation)
				conversations.PUT("/:id", controllers.UpdateConversation)
				conversations.DELETE("/:id", controllers.DeleteConversation)
				conversations.POST("/:id/members", controllers.AddMember)
				conversations.DELETE("/:id/members/:userId", controllers.RemoveMember)
			}

			// Message routes
			messages := protected.Group("/messages")
			{
				messages.POST("", controllers.SendMessage)
				messages.GET("/conversation/:conversationId", controllers.GetMessages)
				messages.PUT("/:id/read", controllers.MarkAsRead)
				messages.DELETE("/:id", controllers.DeleteMessage)
			}

			// File upload
			protected.POST("/upload", controllers.UploadFile)
		}

		// WebSocket route
		v1.GET("/ws", middleware.AuthMiddleware(), websocket.HandleWebSocket)
	}
}
