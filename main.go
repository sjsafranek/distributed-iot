package main

import (
	// "bufio"
	"flag"
	"fmt"
	"log"
	// "net"
	"os"
	"time"

	"github.com/adrianmo/go-nmea"
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
	tripId      = ""
)

func GetTripId() string {
	return tripId
}

func parse(data string, writeQueue chan *Event) {
	message, err := nmea.Parse(data)
	if err != nil {
		// log.Fatal(err)
		// logger.Warn(err)
		// logger.Debug(message)
	}

	if message.DataType() == nmea.TypeRMC {
		m := message.(nmea.RMC)

		// Create timestamp
		year := fmt.Sprintf("%02d", m.Date.YY)
		month := fmt.Sprintf("%02d", m.Date.MM)
		day := fmt.Sprintf("%02d", m.Date.DD)
		hour := fmt.Sprintf("%02d", m.Time.Hour)
		minute := fmt.Sprintf("%02d", m.Time.Minute)
		second := fmt.Sprintf("%02d", m.Time.Second)
		datetimeString := fmt.Sprintf("20%v-%v-%vT%v:%v:%v.%vZ", year, month, day, hour, minute, second, m.Time.Millisecond)
		timestamp, err := time.Parse(time.RFC3339, datetimeString)
		if err != nil {
			fmt.Println("Error while parsing date :", err)
		}

		// Knots to Meters per Second
		speed := m.Speed * 0.514444

		// Create times eries event
		writeQueue <- &Event{
			TripID:    tripId,
			Timestamp: timestamp,
			Longitude: m.Longitude,
			Latitude:  m.Latitude,
			Speed:     speed,
			Heading:   m.Course,
		}
	}

	// logger.Debug(message.DataType())
	// logger.Info(message)
}

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

	tripId = uuid.New().String()
}

func main() {

	// Connect to NMEA Reciever
	conn, err := Open(config.Host, config.Port)
	if nil != err {
		log.Fatal(err)
	}

	// Open database connection
	db, err := NewDatabase()
	if nil != err {
		log.Fatal(err)
	}
	defer db.Close()

	// Make database write queue
	queue := make(chan *Event, 10)
	go func() {
		for event := range queue {
			logger.Info(event)
			err := db.InsertEvent(event)
			if nil != err {
				log.Fatal(err)
			}
		}
	}()

	// Read messages from NMEA Reciever
	quit := make(chan bool)
	conn.ReadAll(quit, func(message string, err error) {
		if nil != err {
			quit <- true
			log.Fatal(err)
		}
		go parse(message, queue)
	})

}
