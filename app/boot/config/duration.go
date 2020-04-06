package config

import "time"

type Duration struct {
	time.Duration
}

func (this *Duration) UnmarshalText(text []byte) error {
	var err error
	this.Duration, err = time.ParseDuration(string(text))
	return err
}

func (this Duration) MarshalText() (text []byte, err error) {
	s := this.Duration.String()
	return []byte(s), nil
}

func (this Duration) String() string {
	return this.Duration.String()
}
