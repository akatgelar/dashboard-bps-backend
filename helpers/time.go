package helpers

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

// Time is a time.Time that marshals to JSON as "yyyy-MM-dd HH:mm:ss"
// (no T/Z). A zero value marshals to an empty string.
type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + t.Format(timeFormat) + `"`), nil
}

func (t *Time) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		t.Time = v
	case *time.Time:
		if v == nil {
			t.Time = time.Time{}
		} else {
			t.Time = *v
		}
	case string:
		parsed, err := time.ParseInLocation(timeFormat, v, time.Local)
		if err != nil {
			t.Time = time.Time{}
		} else {
			t.Time = parsed
		}
	default:
		return fmt.Errorf("cannot scan type %T into models.Time", value)
	}
	return nil
}

func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}
