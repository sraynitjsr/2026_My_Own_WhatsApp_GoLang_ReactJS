package controllers

import (
	"net/http"
	"time"
	"whatsapp-clone/config"
	"whatsapp-clone/models"
	"whatsapp-clone/websocket"

	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ConversationID uint               `json:"conversation_id" binding:"required"`
		Content        string             `json:"content"`
		MessageType    models.MessageType `json:"message_type"`
		FileURL        string             `json:"file_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is member of conversation
	var member models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ?", req.ConversationID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Create message
	message := models.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		Content:        req.Content,
		MessageType:    req.MessageType,
		FileURL:        req.FileURL,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Load sender info
	config.DB.Preload("Sender").First(&message, message.ID)

	// Broadcast message via WebSocket
	websocket.BroadcastMessage(message)

	c.JSON(http.StatusCreated, message)
}

func GetMessages(c *gin.Context) {
	conversationID := c.Param("conversationId")
	userID := c.GetUint("user_id")

	// Check if user is member
	var member models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var messages []models.Message
	if err := config.DB.Where("conversation_id = ?", conversationID).
		Preload("Sender").
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func MarkAsRead(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetUint("user_id")

	var message models.Message
	if err := config.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Check if user is member of conversation
	var member models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ?", message.ConversationID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	now := time.Now()
	message.IsRead = true
	message.ReadAt = &now

	if err := config.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, message)
}

func DeleteMessage(c *gin.Context) {
	messageID := c.Param("id")
	userID := c.GetUint("user_id")

	var message models.Message
	if err := config.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Only sender can delete
	if message.SenderID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only sender can delete message"})
		return
	}

	if err := config.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
