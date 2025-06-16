package models

import "time"

type Job struct {
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Pid       string        `json:"pid"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
}
