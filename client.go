package main

type Client interface {
	ReadLine() (string, error)
	ReadAll(chan bool, func(string, error))
}
