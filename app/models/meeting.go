package models

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type Meeting struct {
	ID string `gorm:"primaryKey"`

	Title, MeetingDescription string

	StartTime, EndTime   time.Time
	CreatedAt, UpdatedAt *time.Time

	MeetingParticipants []*MeetingParticipant
}

type MeetingParticipant struct {
	ID string `gorm:"primaryKey"`

	MeetingID    *string
	UserID       *string
	AcceptStatus *string

	CreatedAt, UpdatedAt *time.Time
}

func (meeting *Meeting) BeforeCreate(_ *gorm.DB) (err error) {
	if meeting.ID != "" {
		return nil
	}

	kid, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		return err
	}

	meeting.ID = fmt.Sprintf("meet_%s", kid.String())
	return nil
}

func (participant *MeetingParticipant) BeforeCreate(_ *gorm.DB) (err error) {
	if participant.ID != "" {
		return nil
	}

	kid, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		return err
	}

	participant.ID = fmt.Sprintf("prt_%s", kid.String())
	return nil
}
