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
	DEFAULT_SERIAL_PORT = ""
)

var (
	CONFIG_FILE = DEFAULT_CONFIG_FILE
	config      = NewConfig()
)

func init() {
	var report_version bool
	host := DEFAULT_HOST
	port := DEFAULT_PORT
	serialPort := DEFAULT_SERIAL_PORT

	flag.BoolVar(&report_version, "version", false, "Version")
	flag.StringVar(&serialPort, "serial_port", DEFAULT_SERIAL_PORT, "NMEA Serial Port")
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
	config.SerialPort = serialPort
}

func main() {

	tripId := uuid.New().String()

	// Connect to NMEA Reciever
	logger.Info("Connecting to data stream")
	var conn Client
	var err error
	// if "" == config.SerialPort {
	conn, err = OpenStream(config.Host, config.Port)
	if nil != err {
		log.Fatal(err)
	}
	// } else {
	// 	conn, err = OpenSerial(config.SerialPort)
	// 	if nil != err {
	// 		log.Fatal(err)
	// 	}
	// }

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

		// logger.Debug(message)

	})

}
