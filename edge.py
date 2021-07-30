

import uuid
import socket
import hashlib
import logging
import sqlite3
import argparse
from datetime import date
from datetime import datetime

import pynmea2
# import bluetooth


# Configuration
BUFF_SIZE = 1024    # 4 KiB


# Setup logging
logger = logging.getLogger()
logger.setLevel(logging.DEBUG)
streamhandler = logging.StreamHandler()
formatter = logging.Formatter("%(asctime)s [%(levelname)s] [%(threadName)s] %(filename)s %(funcName)s:%(lineno)d %(message)s", datefmt='%Y-%m-%d %H:%M:%S')
streamhandler.setFormatter(formatter)
logger.addHandler(streamhandler)



def recvall(sock):
    sock.settimeout(600)
    buff = ''
    while True:

        chunk = sock.recv(BUFF_SIZE)
        if 0 == len(chunk):
            break   # EOF
        chunk = chunk.decode()
        buff += chunk
        # Look for newline character
        if '\n' in buff:
            # Parse each message
            messages = buff.split('\n')
            if not '\n' == buff[-1]:
                buff = ''
            else:
                buff = messages[-1]
                messages.pop()
            for message in messages:
                yield pynmea2.parse(message)


def openTcpConnection(host, port):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((host, port))
    return sock


def openBluetoothConnection(host, port):
    sock = socket.socket(socket.AF_BLUETOOTH, socket.SOCK_STREAM, socket.BTPROTO_RFCOMM)
    # sock = bluetooth.BluetoothSocket( bluetooth.RFCOMM )
    sock.connect((host, port))
    return sock


def connect(protocol, host, port):
    if 'tcp' == protocol.lower():
        return openTcpConnection(host, port)
    elif 'bluetooth' == protocol.lower():
        return openBluetoothConnection(host, port)



def main(protocol, host, port):
    # Generate tripId
    tripId = str(uuid.uuid4())

    # Open database connection
    with sqlite3.connect('edge.db') as conn:

        # Get database cursor
        cursor = conn.cursor()

        # Create table(s)
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS history (
                trip_id             TEXT,
                event_timestamp     TIMESTAMP,
                longitude           DOUBLE PRECISION,
                latitude            DOUBLE PRECISION,
                altitude            INTEGER
            );
        ''')
        # speed       INTEGER,
        # heading     INTEGER,
        # temperature INTEGER,

        # Open socket to GPS Reciever

        # with openTCP as sock:
            # sock.connect((host, port))
        with connect(protocol, host, port) as sock:

            # Read off stream
            prev = ''
            for message in recvall(sock):

                # Check for GGA messages
                if type(message) is pynmea2.types.talker.GGA:

                    event = (
                        tripId,
                        datetime.combine(date.today(), message.timestamp).isoformat(),
                        round(message.longitude, 8),
                        round(message.latitude, 8),
                        int(message.altitude)
                    )

                    # Dedupe event records
                    m = hashlib.md5()
                    for item in event:
                        m.update(str(item).encode())
                    hsh = m.hexdigest()
                    if hsh == prev:
                        continue
                    prev = hsh

                    # Log event message
                    logger.info(event)

                    # Insert a row of data
                    cursor.execute("INSERT INTO history (trip_id,event_timestamp,longitude,latitude,altitude) VALUES (?,?,?,?,?)", (event))

                    # Save (commit) the changes
                    conn.commit()





if __name__ == "__main__":
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Service for collecting location data on edge devices')
    parser.add_argument('-protocol', type=str, default='tcp', help='Protocol for NMEA device connection')
    parser.add_argument('-host', type=str, default='localhost', help='Host for NMEA device')
    parser.add_argument('-port', type=int, default=4352, help='Port for NMEA device')
    args, unknown = parser.parse_known_args()
    main(args.protocol, args.host, args.port)















# def is_socket_closed(sock):
#     ''' Checks if socket is still connected
#         https://stackoverflow.com/questions/48024720/python-how-to-check-if-socket-is-still-connected
#     '''
#     try:
#         # this will try to read bytes without blocking and also without removing them from buffer (peek only)
#         data = sock.recv(16, socket.MSG_DONTWAIT | socket.MSG_PEEK)
#         if len(data) == 0:
#             return True
#     except BlockingIOError:
#         return False  # socket is open and reading from it would block
#     except ConnectionResetError:
#         return True  # socket was closed for some other reason
#     except Exception as e:
#         log.exception("unexpected exception when checking if a socket is closed")
#         return False
#     return False
