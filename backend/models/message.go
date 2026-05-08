package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageType string

const (
	TextMessage  MessageType = "text"
	ImageMessage MessageType = "image"
	FileMessage  MessageType = "file"
	VideoMessage MessageType = "video"
	AudioMessage MessageType = "audio"
)

type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ConversationID uint        `gorm:"not null;index" json:"conversation_id"`
	SenderID       uint        `gorm:"not null;index" json:"sender_id"`
	Content        string      `json:"content"`
	MessageType    MessageType `gorm:"type:varchar(20);default:'text'" json:"message_type"`
	FileURL        string      `json:"file_url,omitempty"`
	IsRead         bool        `gorm:"default:false" json:"is_read"`
	ReadAt         *time.Time  `json:"read_at,omitempty"`

	// Relationships
	Sender       User         `gorm:"foreignKey:SenderID" json:"sender"`
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"-"`
}
