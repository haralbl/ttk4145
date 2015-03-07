package main

import (
	."fmt"
	"status"
	"network"
	//"time"
	"timer"
	"driver"
	//"strings"
)

var (
	message string
)

func main() {
	
	sendChan				:= make(chan string)
	receiveChan 			:= make(chan string)
	receiveAliveMessageChan := make(chan string)
	elevatorTimerChan		:= make(chan int)
	elevatorTimeoutChan		:= make(chan int)
	doorTimerChan			:= make(chan int)
	doorTimeoutChan			:= make(chan int)
	
	doneChan				:= make(chan string)
	
	upButtonChan			:= make(chan int)
	downButtonChan			:= make(chan int)
	commandButtonChan		:= make(chan int)
	floorChan				:= make(chan int)
	
	status.Initialize()
	
	
	go network.Send(sendChan)
	go network.Receive(receiveChan)
	go network.SendAliveMessage()
	go network.ReceiveAliveMessage(receiveAliveMessageChan)
	
	go timer.ElevatorTimer(elevatorTimerChan, elevatorTimeoutChan)
	go timer.DoorTimer(doorTimerChan, doorTimeoutChan)
	
	go status.CheckAliveElevators(receiveAliveMessageChan, elevatorTimerChan)
	go status.EventHandler(sendChan, upButtonChan, downButtonChan, commandButtonChan, floorChan)
	
	//driver.Test()
	driver.Init()
	//poller lagt til
	
	go driver.UpButtonPoller(upButtonChan)
	go driver.DownButtonPoller(downButtonChan)
	go driver.CommandButtonPoller(commandButtonChan)
	go driver.FloorPoller(floorChan)
	
/*	
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
*/
	go func() {
		for {
			message = <- receiveChan
			println(message)
		}
	}()
	
	/*for {
		time.Sleep(time.Millisecond * 1000)
		st := status.Get()
		Printf("Active: %s\n", st)
	}*/
		



	
	Println(<-doneChan)
}
