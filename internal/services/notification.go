package services

import (
	"fmt"
	"log"
	"time"

	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService() *NotificationService {
	return &NotificationService{db: database.GetDB()}
}

func (s *NotificationService) Create(userID uint64, title, message, notifType, priority string, metadata string) (*models.Notification, error) {
	if metadata == "" {
		metadata = "{}"
	}
	n := &models.Notification{
		UserID:           userID,
		Title:            title,
		Message:          message,
		NotificationType: notifType,
		Priority:         priority,
		Metadata:         metadata,
		CreatedAt:        time.Now(),
	}
	if err := s.db.Create(n).Error; err != nil {
		log.Printf("notification create failed: %v (user=%d title=%s)", err, userID, title)
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	log.Printf("notification created: id=%d user=%d title=%s", n.ID, userID, title)
	return n, nil
}

func (s *NotificationService) List(userID uint64, onlyUnread bool) ([]models.Notification, error) {
	var items []models.Notification
	q := s.db.Where("user_id = ?", userID).Order("created_at DESC")
	if onlyUnread {
		q = q.Where("is_read = ?", false)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	return items, nil
}

func (s *NotificationService) MarkRead(userID, id uint64) error {
	res := s.db.Model(&models.Notification{}).Where("user_id = ? AND id = ?", userID, id).Update("is_read", true)
	if res.Error != nil {
		return fmt.Errorf("failed to mark read: %w", res.Error)
	}
	return nil
}
