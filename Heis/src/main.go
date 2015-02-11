package main

import (
	."fmt"
	"status"
	"network"
	"time"
	//"timer"
	//"driver"
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
	

	
	go network.Send(sendChan)
	go network.Receive(receiveChan)
	go network.SendAliveMessage()
	go network.ReceiveAliveMessage(receiveAliveMessageChan)
	//go timer.ElevatorTimer(elevatorTimerChan)
	go status.CheckAliveElevators(receiveAliveMessageChan, elevatorTimerChan)
	//go driver.Test()
	
	
	time.Sleep(time.Millisecond * 1000)
	
	message = status.GetStatus()
	println(message)
	sendChan <- message
	
	message = <- receiveChan
	status.Update(message)
	message = status.GetStatus()
	println(message)
	
	
	
	Println(<-doneChan)
}