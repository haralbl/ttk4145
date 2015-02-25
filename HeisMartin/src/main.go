package main

import (
	."fmt"
	"status"
	"network"
	"time"
	//"timer"
	"driver"
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

	//lagt til
	upButtonChan			:= make(chan int)
	downButtonChan			:= make(chan int)
	commandButtonChan		:= make(chan int)
	floorChan				:= make(chan int)

	//go over denne?
	status.Initialize()
	

	
	go network.Send(sendChan)
	go network.Receive(receiveChan)
	go network.SendAliveMessage()
	go network.ReceiveAliveMessage(receiveAliveMessageChan)
	//go timer.ElevatorTimer(elevatorTimerChan)
	go status.CheckAliveElevators(receiveAliveMessageChan, elevatorTimerChan)
	go driver.Test()

	//poller lagt til
	go driver.UpButtonPoller(upButtonChan)
	go driver.DownButtonPoller(downButtonChan)
	go driver.CommandButtonPoller(commandButtonChan)
	go driver.FloorPoller(floorChan)
	
	
	
	
	
	//lese knappetrykk
	go func(upButtonChan chan int) {
		for {
			//Println("hei")
			Println(<- upButtonChan)
	}
	}(upButtonChan)

	go func(downButtonChan chan int) {
		for {
			Println(<- downButtonChan)
		}
	}(downButtonChan)
	
	go func(commandButtonChan chan int) {
		for {
			Println(<- commandButtonChan)
		}
	}(commandButtonChan)
	
	go func(floorChan chan int) {
		for {
			Println(<- floorChan)
		}
	}(floorChan)
	
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