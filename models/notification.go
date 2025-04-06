package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey"`
	UserID  uint   `json:"user_id" gorm:"not null"`
	Message string `json:"message"`
	Status  string `json:"status" gorm:"type:enum('unread', 'read');default:'unread'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}