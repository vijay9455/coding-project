package repository

import (
	"calendly/app/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type MeetingRepositoryInterface interface {
	GetForUser(ctx context.Context, db *gorm.DB, userID string, startTime, endTime time.Time) ([]*models.Meeting, error)
}

func NewMeetingRepository() MeetingRepositoryInterface {
	return &meetingRepository{}
}

type meetingRepository struct{}

func (repo *meetingRepository) GetForUser(ctx context.Context, db *gorm.DB, userID string, startTime, endTime time.Time) ([]*models.Meeting, error) {
	var meetings []*models.Meeting
	err := db.Table("meetings").Joins("inner join meeting_participants on meeting_participants.meeting_id = meetings.id").
		Where("meetings.end_time >= ? or meetings.start_time <= ?", startTime, endTime).Where("meeting_participants.user_id = ?", userID).
		Order("start_time, end_time asc").
		Find(&meetings).Error
	if err != nil {
		return nil, err
	}

	return meetings, nil
}
