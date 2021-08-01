package main

import (
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Title      string `toml:"title"`
	DeviceID   string `toml:"device_id"`
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	SerialPort string `toml:"serial_port"`
}

func (self *Config) Fetch(file string) error {
	b, err := ioutil.ReadFile(file)
	if nil != err {
		return err
	}
	return self.Unmarshal(string(b))
}

func (self *Config) Save(file string) error {
	contents, err := self.Marshal()
	if nil != err {
		return err
	}
	return ioutil.WriteFile(file, []byte(contents), 0644)
}

func (self *Config) Unmarshal(data string) error {
	return toml.Unmarshal([]byte(data), self)
}

func (self Config) Marshal() (string, error) {
	b, err := toml.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), nil
}

func NewConfig() *Config {
	return &Config{
		Title:    "Node",
		DeviceID: uuid.New().String(),
	}
}
