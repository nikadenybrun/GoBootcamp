package model

import "time"

type Message struct {
	SessionID string
	Frequency float64
	Timestamp time.Time
}
