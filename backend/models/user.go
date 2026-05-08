package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username    string     `gorm:"uniqueIndex;not null" json:"username"`
	Email       string     `gorm:"uniqueIndex;not null" json:"email"`
	Password    string     `gorm:"not null" json:"-"`
	DisplayName string     `json:"display_name"`
	Avatar      string     `json:"avatar"`
	Bio         string     `json:"bio"`
	IsOnline    bool       `gorm:"default:false" json:"is_online"`
	LastSeenAt  *time.Time `json:"last_seen_at"`

	// Relationships
	Messages            []Message            `gorm:"foreignKey:SenderID" json:"-"`
	ConversationMembers []ConversationMember `gorm:"foreignKey:UserID" json:"-"`
}
