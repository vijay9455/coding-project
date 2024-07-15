package models

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type UserAvailability struct {
	ID        string `gorm:"primaryKey"`
	UserID    *string
	DayOfWeek int64

	StartTime, EndTime TimeOnly

	CreatedAt, UpdatedAt *time.Time
}

func (ua *UserAvailability) BeforeCreate(_ *gorm.DB) (err error) {
	if ua.ID != "" {
		return nil
	}

	kid, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		return err
	}

	ua.ID = fmt.Sprintf("ua_%s", kid.String())
	return nil
}
