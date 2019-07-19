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
