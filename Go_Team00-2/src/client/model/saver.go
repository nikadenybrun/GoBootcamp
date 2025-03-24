package model

import (
	"fmt"
)

type SaverMean struct {
	Count     int64
	Mean      float64
	Deviation float64
}

func (s SaverMean) String() string {
	return fmt.Sprintf("%d,%f,%f\n", s.Count, s.Mean, s.Deviation)
}

type SaverAnomaly struct {
	Frequency float64
	Diff      float64
}

func (s SaverAnomaly) String() string {
	return fmt.Sprintf("%f,%f\n", s.Frequency, s.Diff)
}

type SaverEntry struct {
	SessionId string
	Frequency float64
	Timestamp int64
}

func (s SaverEntry) String() string {
	return fmt.Sprintf("%s,%f,%d\n", s.SessionId, s.Frequency, s.Timestamp)
}

type SaverMessage interface {
	String() string
}
