package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	// "time"

	// "github.com/adrianmo/go-nmea"
	"github.com/google/uuid"
	"github.com/sjsafranek/logger"
)

const (
	ProjectName  = "nmea-parser"
	VersionMajor = 0
	VersionMinor = 1
	VersionPatch = 0
)

const (
	DEFAULT_HOST        = "127.0.0.1"
	DEFAULT_PORT        = 4352
	DEFAULT_CONFIG_FILE = "config.toml"
)

var (
	CONFIG_FILE = DEFAULT_CONFIG_FILE
	config      = NewConfig()
)

//
// type RMC struct {
// 	nmea.RMC
// }
//
// func (self *RMC) GetTimestamp() (time.Time, error) {
// 	// Create timestamp
// 	year := fmt.Sprintf("%02d", self.Date.YY)
// 	month := fmt.Sprintf("%02d", self.Date.MM)
// 	day := fmt.Sprintf("%02d", self.Date.DD)
// 	hour := fmt.Sprintf("%02d", self.Time.Hour)
// 	minute := fmt.Sprintf("%02d", self.Time.Minute)
// 	second := fmt.Sprintf("%02d", self.Time.Second)
// 	datetimeString := fmt.Sprintf("20%v-%v-%vT%v:%v:%v.%vZ", year, month, day, hour, minute, second, self.Time.Millisecond)
// 	return time.Parse(time.RFC3339, datetimeString)
// }
//
// func (self *RMC) GetSpeed() float64 {
// 	return self.Speed * 0.514444
// }

// func parse(data string) *Event {
// 	message, err := nmea.Parse(data)
// 	if err != nil {
// 		// log.Fatal(err)
// 		// logger.Warn(err)
// 		// logger.Debug(message)
// 	}
//
// 	if message.DataType() == nmea.TypeRMC {
// 		rmc := message.(nmea.RMC)
//
// 		// Create timestamp
// 		year := fmt.Sprintf("%02d", rmc.Date.YY)
// 		month := fmt.Sprintf("%02d", rmc.Date.MM)
// 		day := fmt.Sprintf("%02d", rmc.Date.DD)
// 		hour := fmt.Sprintf("%02d", rmc.Time.Hour)
// 		minute := fmt.Sprintf("%02d", rmc.Time.Minute)
// 		second := fmt.Sprintf("%02d", rmc.Time.Second)
// 		datetimeString := fmt.Sprintf("20%v-%v-%vT%v:%v:%v.%vZ", year, month, day, hour, minute, second, rmc.Time.Millisecond)
// 		timestamp, err := time.Parse(time.RFC3339, datetimeString)
// 		if err != nil {
// 			fmt.Println("Error while parsing date :", err)
// 		}
//
// 		// Knots to Meters per Second
// 		speed := rmc.Speed * 0.514444
//
// 		// Create times eries event
// 		return &Event{
// 			// TripID:    tripId,
// 			Timestamp: timestamp,
// 			Longitude: rmc.Longitude,
// 			Latitude:  rmc.Latitude,
// 			Speed:     speed,
// 			Heading:   rmc.Course,
// 		}
// 	}
//
// 	return nil
// 	// logger.Debug(message.DataType())
// 	// logger.Info(message)
// }

func init() {
	var report_version bool
	host := DEFAULT_HOST
	port := DEFAULT_PORT

	flag.BoolVar(&report_version, "version", false, "Version")
	flag.StringVar(&host, "host", DEFAULT_HOST, "NMEA Reciever host")
	flag.IntVar(&port, "port", DEFAULT_PORT, "NMEA Reciever port")
	flag.StringVar(&CONFIG_FILE, "config", DEFAULT_CONFIG_FILE, "Configuration file")
	flag.Parse()

	if report_version {
		fmt.Printf("%v v%v.%v.%v\n", ProjectName, VersionMajor, VersionMinor, VersionPatch)
		os.Exit(0)
	}

	err := config.Fetch(CONFIG_FILE)
	if nil != err {
		config.Host = host
		config.Port = port
		config.Save(CONFIG_FILE)
	}

	config.Host = host
	config.Port = port
}

func main() {

	tripId := uuid.New().String()

	// Connect to NMEA Reciever
	logger.Info("Connecting to data stream")
	conn, err := Open(config.Host, config.Port)
	if nil != err {
		log.Fatal(err)
	}

	// Open database connection
	logger.Info("Opening database")
	db, err := NewDatabase()
	if nil != err {
		log.Fatal(err)
	}
	defer db.Close()

	// Read messages from NMEA Reciever
	logger.Info("Reading from data stream")
	quit := make(chan bool)
	conn.ReadAll(quit, func(message string, err error) {
		if nil != err {
			quit <- true
			log.Fatal(err)
		}

		// Parse NMAE RMC messages
		if strings.HasPrefix(message, "$GPRMC") {
			logger.Info(message)

			rmc, err := parseRMC(message)
			if nil != err {
				logger.Error(err)
				return
			}

			event := Event{
				TripID:    tripId,
				Timestamp: rmc.GetTimestamp(),
				Longitude: rmc.GetLongitude(),
				Latitude:  rmc.GetLatitude(),
				Speed:     rmc.GetSpeed(),
				Heading:   rmc.GetHeading(),
			}

			err = db.InsertEvent(&event)
			if nil != err {
				log.Fatal(err)
			}
		}

	})

}
