package log_models

import "time"

// ActivityLog model for storing download activities
type ActivityLog struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	Action    string    `gorm:"not null"`
	ItemID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

// Counter model for storing download counts
type Counter struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"uniqueIndex;type:varchar(255);not null"` // Specify a maximum length
	Count uint   `gorm:"not null"`
}
