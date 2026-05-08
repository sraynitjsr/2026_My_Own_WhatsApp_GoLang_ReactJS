package controllers

import (
	"net/http"
	"time"
	"whatsapp-clone/config"
	"whatsapp-clone/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateConversation(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Type        models.ConversationType `json:"type" binding:"required"`
		Name        string                  `json:"name"`
		Description string                  `json:"description"`
		MemberIDs   []uint                  `json:"member_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create conversation
	conversation := models.Conversation{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := config.DB.Create(&conversation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	// Add members
	memberIDs := append(req.MemberIDs, userID) // Include current user
	for _, memberID := range memberIDs {
		member := models.ConversationMember{
			ConversationID: conversation.ID,
			UserID:         memberID,
			Role:           "member",
			JoinedAt:       time.Now(),
		}
		if memberID == userID {
			member.Role = "admin"
		}
		config.DB.Create(&member)
	}

	// Load conversation with members
	config.DB.Preload("Members.User").First(&conversation, conversation.ID)

	c.JSON(http.StatusCreated, conversation)
}

func GetConversations(c *gin.Context) {
	userID := c.GetUint("user_id")

	var conversations []models.Conversation
	if err := config.DB.
		Joins("JOIN conversation_members ON conversation_members.conversation_id = conversations.id").
		Where("conversation_members.user_id = ?", userID).
		Preload("Members.User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			// Get last message only
			return db.Order("created_at DESC").Limit(1)
		}).
		Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

func GetConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetUint("user_id")

	var conversation models.Conversation
	if err := config.DB.Preload("Members.User").First(&conversation, conversationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// Check if user is a member
	var member models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func UpdateConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetUint("user_id")

	// Check if user is admin
	var member models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ? AND role = ?", conversationID, userID, "admin").First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update conversation"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.Conversation{}).Where("id = ?", conversationID).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update conversation"})
		return
	}

	var conversation models.Conversation
	config.DB.Preload("Members.User").First(&conversation, conversationID)

	c.JSON(http.StatusOK, conversation)
}

func DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetUint("user_id")

	// Check if user is admin
	var member models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ? AND role = ?", conversationID, userID, "admin").First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete conversation"})
		return
	}

	if err := config.DB.Delete(&models.Conversation{}, conversationID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted"})
}

func AddMember(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetUint("user_id")

	// Check if requester is admin
	var requester models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ? AND role = ?", conversationID, userID, "admin").First(&requester).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can add members"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member := models.ConversationMember{
		ConversationID: uint(conversationID[0]),
		UserID:         req.UserID,
		Role:           "member",
		JoinedAt:       time.Now(),
	}

	if err := config.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added"})
}

func RemoveMember(c *gin.Context) {
	conversationID := c.Param("id")
	targetUserID := c.Param("userId")
	requestUserID := c.GetUint("user_id")

	// Check if requester is admin
	var requester models.ConversationMember
	if err := config.DB.Where("conversation_id = ? AND user_id = ? AND role = ?", conversationID, requestUserID, "admin").First(&requester).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can remove members"})
		return
	}

	if err := config.DB.Where("conversation_id = ? AND user_id = ?", conversationID, targetUserID).Delete(&models.ConversationMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}
