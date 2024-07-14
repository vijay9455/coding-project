package models

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type User struct {
	ID    string `gorm:"primaryKey"`
	Email string

	FirstName, LastName  string
	CreatedAt, UpdatedAt *time.Time
}

func (user *User) BeforeCreate(_ *gorm.DB) (err error) {
	if user.ID != "" {
		return nil
	}

	kid, err := ksuid.NewRandomWithTime(time.Now())
	if err != nil {
		return err
	}

	user.ID = fmt.Sprintf("usr_%s", kid.String())
	return nil
}
