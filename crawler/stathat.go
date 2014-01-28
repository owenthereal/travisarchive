package main

import (
	"time"

	sh "github.com/stathat/go"
)

type Stathat struct {
	StatName string
	Ezkey    string
}

func (s *Stathat) CountOne(time *time.Time) error {
	if s.Ezkey == "" || s.StatName == "" {
		return nil
	}

	return sh.PostEZCountTime(s.StatName, s.Ezkey, 1, time.Unix())
}
