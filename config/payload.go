package config

import "time"

type Payload struct {
	ID      string    `json:"id"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
