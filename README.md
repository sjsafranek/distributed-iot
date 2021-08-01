# distributed-iot














https://miloserdov.org/?p=3762
https://unix.stackexchange.com/questions/522265/how-to-get-gps-data-from-android-phone-to-debian


# BlueNMEA

https://github.com/MaxKellermann/BlueNMEA
http://max.kellermann.name/projects/blue-nmea/
http://max.kellermann.name/download/blue-nmea/BlueNMEA-2.1.3.apk

```bash
$ sudo python3 -m pip install --upgrade pynmea2

$ sudo apt install adb
$ adb devices
$ adb forward tcp:4352 tcp:4352
```

sudo python3 -m pip install --upgrade PyBluez






sudo python3 -m pip install --upgrade pynmea2
sudo python3 -m pip install --upgrade pySerial









telnet 192.168.0.21 4352




for port in {1..30}; do
	echo $port
	python3 edge.py -protocol bluetooth -host 'F8:E6:1A:DB:59:84' -port $port
done



bluetooth
$ hcitool dev
