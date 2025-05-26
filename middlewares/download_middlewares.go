package middlewares

import (
	logM "project/models/log_models"
	"time"

	"gorm.io/gorm"
)

func LogDownloadActivity(db *gorm.DB, userID, bookID uint) error {
	activity := logM.ActivityLog{
		UserID:    userID,
		Action:    "download",
		ItemID:    bookID,
		CreatedAt: time.Now(),
	}
	return db.Create(&activity).Error
}

func IncrementDownloadCounter(db *gorm.DB) error {
	var counter logM.Counter
	err := db.Where("name = ?", "downloads").First(&counter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If counter doesn't exist, create it
			counter = logM.Counter{Name: "downloads", Count: 1}
			return db.Create(&counter).Error
		}
		return err
	}

	// Increment the counter
	counter.Count++
	return db.Save(&counter).Error
}
