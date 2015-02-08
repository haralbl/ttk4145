package main

import (
	."fmt"
	"./network"
	"./status"
)

var (
	message string
)

func main() {
	sendChan	:= make(chan string)
	receiveChan := make(chan string)
	doneChan	:= make(chan string)
	
	status.Initialize()
	
	go network.Send(sendChan)
	go network.Receive(receiveChan)
	
	
	message = status.GetStatus()
	println(message)
	sendChan <- message
	
	
	

	
	
	
	
	
	
	
	
	
	
	Println(<-doneChan)
}
