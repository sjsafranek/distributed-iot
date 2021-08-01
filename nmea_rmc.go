package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/adrianmo/go-nmea"
)

type RMC struct {
	nmea.RMC
}

func (self *RMC) GetTimestamp() time.Time {
	// Create timestamp
	year := fmt.Sprintf("%02d", self.Date.YY)
	month := fmt.Sprintf("%02d", self.Date.MM)
	day := fmt.Sprintf("%02d", self.Date.DD)
	hour := fmt.Sprintf("%02d", self.Time.Hour)
	minute := fmt.Sprintf("%02d", self.Time.Minute)
	second := fmt.Sprintf("%02d", self.Time.Second)
	datetimeString := fmt.Sprintf("20%v-%v-%vT%v:%v:%v.%vZ", year, month, day, hour, minute, second, self.Time.Millisecond)

	// Swallowing the error to make the api match..
	timestamp, _ := time.Parse(time.RFC3339, datetimeString)
	return timestamp
}

func (self *RMC) GetSpeed() float64 {
	return self.Speed * 0.514444
}

func (self *RMC) GetHeading() float64 {
	return self.Course
}

func (self *RMC) GetLongitude() float64 {
	return self.Longitude
}

func (self *RMC) GetLatitude() float64 {
	return self.Latitude
}

func parseRMC(raw string) (*RMC, error) {
	message, err := nmea.Parse(raw)
	if err != nil {
		return &RMC{}, err
	}

	if message.DataType() == nmea.TypeRMC {
		rmc := message.(nmea.RMC)
		return &RMC{rmc}, nil
	}

	return &RMC{}, errors.New("Wrong message type")
}
