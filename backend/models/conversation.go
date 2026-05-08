package models

import (
	"time"

	"gorm.io/gorm"
)

type ConversationType string

const (
	DirectConversation ConversationType = "direct"
	GroupConversation  ConversationType = "group"
)

type Conversation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Type        ConversationType `gorm:"type:varchar(20);default:'direct'" json:"type"`
	Name        string           `json:"name,omitempty"`        // For group chats
	Avatar      string           `json:"avatar,omitempty"`      // For group chats
	Description string           `json:"description,omitempty"` // For group chats

	// Relationships
	Messages []Message            `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
	Members  []ConversationMember `gorm:"foreignKey:ConversationID" json:"members,omitempty"`
}

type ConversationMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ConversationID uint      `gorm:"not null;index" json:"conversation_id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	Role           string    `gorm:"default:'member'" json:"role"` // admin, member
	JoinedAt       time.Time `json:"joined_at"`

	// Relationships
	User         User         `gorm:"foreignKey:UserID" json:"user"`
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"-"`
}
