package models

import "time"

type Time time.Time

func (this Time) MarshalJSON() ([]byte, error) {
	//json格式化的时候必须携带时区信息
	return []byte(`"`+time.Time(this).Format(time.RFC3339)+`"`), nil
}

func (this Time) String() string {
	//string格式化的时候必须携带时区信息
	return time.Time(this).Format(time.RFC3339)
}
