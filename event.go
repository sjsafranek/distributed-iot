package main

import (
	"time"
)

type Event struct {
	TripID    string
	Timestamp time.Time
	Longitude float64
	Latitude  float64
	Speed     float64
	Heading   float64
}
