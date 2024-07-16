package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	DateFormat     = "2006-01-02"
	TimeOnlyFormat = "15:04:05"
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
		parsedTime, err := time.Parse(TimeOnlyFormat, string(v))
		if err != nil {
			return err
		}
		*t = TimeOnly{parsedTime}
		return nil
	case string:
		parsedTime, err := time.Parse(TimeOnlyFormat, v)
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
	return t.Format(fmt.Sprintf(`"%s"`, TimeOnlyFormat)), nil
}

func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	str := string(b)
	parsedTime, err := time.Parse(fmt.Sprintf(`"%s"`, TimeOnlyFormat), str)
	if err != nil {
		return err
	}
	*t = TimeOnly{parsedTime}
	return nil
}

// MarshalJSON customizes the JSON marshalling for TimeOnly
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format(TimeOnlyFormat))), nil
}

func (t TimeOnly) Before(u TimeOnly) bool {
	return t.Time.Before(u.Time)
}

func (t TimeOnly) After(u TimeOnly) bool {
	return t.Time.After(u.Time)
}

func (t TimeOnly) Add(d time.Duration) TimeOnly {
	t.Time = t.Time.Add(d)
	return t
}

type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	parsedDate, err := time.Parse(fmt.Sprintf(`"%s"`, DateFormat), string(b))
	if err != nil {
		return err
	}
	*d = DateOnly{parsedDate}
	return nil
}

// MarshalJSON customizes the JSON marshalling for TimeOnly
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.Format(DateFormat))), nil
}

func (d DateOnly) Add(dur time.Duration) DateOnly {
	d.Time = d.Time.Add(dur)
	return d
}
