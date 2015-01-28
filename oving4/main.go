package main

import (
	"fmt"
	."./network"
)

var (
	message string
)

func main() {
	sendChan := make(chan string)
	doneChan := make(chan string)
	
	go Send(sendChan)
	go Receive()
	
	message = "It's a trap!!!!!!!!!!!!!!!"
	sendChan <- message
	
	
	
	
	
	
	
	fmt.Println(<-doneChan)
}
