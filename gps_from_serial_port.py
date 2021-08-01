import io
import uuid
import sqlite3
import argparse
from datetime import date
from datetime import datetime

import serial
import pynmea2

# serial_port = "/dev/ttyACM0"

def main(port):

    # Open database connection
    with sqlite3.connect('db.sqlite3') as conn:

        # Get database cursor
        cursor = conn.cursor()

        # Create table(s)
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS history (
                trip_id             TEXT,
                event_timestamp     TIMESTAMP,
                longitude           DOUBLE PRECISION,
                latitude            DOUBLE PRECISION,
                speed				REAL,
                heading				REAL
            );
        ''')

        ser = serial.Serial(port, 9600, timeout=5.0)
        sio = io.TextIOWrapper(io.BufferedRWPair(ser, ser))

        tripId = str(uuid.uuid4())

        while True:
            try:
                line = sio.readline()
                msg = pynmea2.parse(line)
                if type(msg) is pynmea2.types.talker.RMC:
                    # Build event
                    timestamp = datetime.combine(msg.datestamp, msg.timestamp).isoformat()
                    event = (tripId, timestamp, msg.longitude, msg.latitude, msg.spd_over_grnd, msg.true_course)

                    print(event)

                    # Insert a row of data
                    cursor.execute('''
                        INSERT INTO history (trip_id, event_timestamp, longitude, latitude, speed, heading)
                            VALUES (?,?,?,?,?,?);
                        ''', (event))

                    # Save (commit) the changes
                    conn.commit()


            except serial.SerialException as e:
                print('Device error: {}'.format(e))
                break
            except pynmea2.ParseError as e:
                print('Parse error: {}'.format(e))
                continue


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Service for collecting location data on edge devices')
    parser.add_argument('-port', type=str, required=True, help='NMEA device serial port')
    args, unknown = parser.parse_known_args()
    main(args.port)
