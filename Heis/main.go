package main

import (
	."fmt"
	"./status"
)

var (
	message string
)

func main() {
	sendChan				:= make(chan string)
	receiveChan 			:= make(chan string)
	receiveAliveMessageChan := make(chan string)
	elevatorTimerChan		:= make(chan int)
	doneChan				:= make(chan string)
	
	status.Initialize()
	
	go status.Network.Send(sendChan)
	go status.Network.Receive(receiveChan)
	go status.Network.SendAliveMessage()
	go status.Network.ReceiveAliveMessage(receiveAliveMessageChan)
	
	go status.timer.elevatorTimer(elevatorTimerChan)
	go status.CheckAliveElevators(receiveAliveMessageChan, elevatorTimerChan)
	
	
	
	
	
	
	message = status.GetStatus()
	println(message)
	sendChan <- message
	
	
	

	
	
	
	
	
	Println(<-doneChan)
}
