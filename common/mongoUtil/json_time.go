package mongoUtil

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

var (
	timeZone = "Asia/Ho_Chi_Minh"
	timeFmt  = "2006-01-02T15:04:05.999999+00:00"
)

// SetTimeFormatter - Set time format layout. Default: 2006-01-02T15:04:05.999999-07:00
func SetTimeFormatter(layout string) {
	timeFmt = layout
}

type JSONTime time.Time

// MarshalJSON - Implement method MarshalJSON to output time with in formatted
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Add(timeAddDuration).Format(timeFmt))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	ti, err := time.Parse(timeFmt, strings.Replace(string(data), "\"", "", -1))

	*t = JSONTime(ti)

	if err != nil {
		return err
	}
	return nil
}

func (t *JSONTime) String() string {
	return time.Time(*t).Format("2006-01-02 15:04:05.999999")
}

// Value - This method for mapping JSONTime to datetime data type in sql
func (t *JSONTime) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return time.Time(*t).Format("2006-01-02 15:04:05.999999"), nil
}

// Scan - This method for scanning JSONTime from datetime data type in sql
func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if v, ok := value.(time.Time); ok {
		*t = JSONTime(v)
		return nil
	}

	return errors.New("invalid Scan Source")
}

func (t *JSONTime) GetBSON() (interface{}, error) {
	if t == nil {
		return nil, nil
	}
	return time.Time(*t), nil
}

func (t *JSONTime) SetBSON(raw bson.RawValue) error {
	var tm time.Time
	tm = raw.Time()
	*t = JSONTime(tm)
	return nil
}
