package models

import "time"

type Bucket struct {
	Count     uint      `json:"count"`
	Key       string    `json:"key"`
	StartTime time.Time `json:"start_time"`
}
