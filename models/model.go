package models

import (
	"encoding/gob"
	"time"
)

type Time struct {
	T time.Time
}

func init() {
	gob.Register(&Time{})
}

func (this Time) MarshalJSON() ([]byte, error) {
	if this.T.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"`+this.T.Format("2006-01-02 15:04:05")+`"`), nil
}

func (this *Time) UnmarshalJSON(data []byte) (error) {
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	this.T = now
	return err
}

func (this Time) String() string {
	return this.T.Format("2006-01-02 15:04:05")
}
