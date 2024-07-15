package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TimeOnly struct {
	time.Time
}

func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		*t = TimeOnly{time.Time{}}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = TimeOnly{v}
		return nil
	case []byte:
		parsedTime, err := time.Parse("15:04:05-07:00", string(v))
		if err != nil {
			return err
		}
		*t = TimeOnly{parsedTime}
		return nil
	case string:
		parsedTime, err := time.Parse("15:04:05-07:00", v)
		if err != nil {
			return err
		}
		*t = TimeOnly{parsedTime}
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into TimeOnly", value)
	}
}

// Value implements the driver.Valuer interface for TimeOnly
func (t TimeOnly) Value() (driver.Value, error) {
	return t.Format("15:04:05-07:00"), nil
}

func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	str := string(b)
	parsedTime, err := time.Parse(`"15:04:05-07:00"`, str)
	if err != nil {
		return err
	}
	*t = TimeOnly{parsedTime}
	return nil
}

// MarshalJSON customizes the JSON marshalling for TimeOnly
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format("15:04:05-07:00"))), nil
}
