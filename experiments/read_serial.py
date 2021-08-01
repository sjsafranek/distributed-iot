import io
import uuid
import pynmea2
import serial
from datetime import date
from datetime import datetime

serial_port = "/dev/ttyACM0"

ser = serial.Serial(serial_port, 9600, timeout=5.0)
sio = io.TextIOWrapper(io.BufferedRWPair(ser, ser))

tripId = str(uuid.uuid4())

while True:
    try:
        line = sio.readline()
        msg = pynmea2.parse(line)
        if type(msg) is pynmea2.types.talker.RMC:
            timestamp = datetime.combine(msg.datestamp, msg.timestamp).isoformat()
            print(tripId, timestamp, msg.longitude, msg.latitude, msg.spd_over_grnd, msg.true_course)
    except serial.SerialException as e:
        print('Device error: {}'.format(e))
        break
    except pynmea2.ParseError as e:
        print('Parse error: {}'.format(e))
        continue
