package main

import (
	."fmt"
	"status"
	"network"
	//"time"
	"timer"
	"driver"
	//"strings"
	"os"
	//"structDefine"
)

var (
	message string
)

func main() {
	
	sendChan				:= make(chan []byte)
	receiveChan 			:= make(chan []byte)
	receiveAliveMessageChan := make(chan string)
	elevatorTimerChan		:= make(chan int)
	elevatorTimeoutChan		:= make(chan int)
	doorTimerChan			:= make(chan string)
	doorTimeoutChan			:= make(chan int)
	//ackTimerChan			:= make(chan string)
	//ackTimeoutChan			:= make(chan string)
	//ackResetChan			:= make(chan string)
	
	doneChan				:= make(chan string)
	initChan				:= make(chan string)
	
	upButtonChan			:= make(chan int)
	downButtonChan			:= make(chan int)
	commandButtonChan		:= make(chan int)
	floorChan				:= make(chan int)
	
	checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan := make(chan []byte)
	
	if driver.Init() != 1 {
		Printf("driver.Init failed")
		os.Exit(1)
	}
	
	go driver.UpButtonPoller(upButtonChan)
	go driver.DownButtonPoller(downButtonChan)
	go driver.CommandButtonPoller(commandButtonChan)
	go driver.FloorPoller(floorChan)
	
	go status.Initialize(initChan, floorChan)
	Println(<-initChan)
	
	go network.SendManager(sendChan, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan)
	
	go network.Receive(receiveChan)
	go network.SendAliveMessage()
	go network.ReceiveAliveMessage(receiveAliveMessageChan)
	
	go timer.ElevatorTimer(elevatorTimerChan, elevatorTimeoutChan)
	go timer.DoorTimer(doorTimerChan, doorTimeoutChan)
	//go timer.AckTimer(ackTimerChan, ackTimeoutChan, ackResetChan)
	
	go status.CheckAliveElevators(receiveAliveMessageChan, elevatorTimerChan, sendChan)
	go status.EventHandler(sendChan, upButtonChan, downButtonChan, commandButtonChan, floorChan,
							/*ackTimerChan,*/ receiveChan, /*ackTimeoutChan, ackResetChan,*/
							doorTimerChan)
							
	go status.CheckIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAdded(sendChan, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan)
	
	
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
	
	
	/*for {
		time.Sleep(time.Millisecond * 1000)
		st := status.Get()
		Printf("Active: %s\n", st)
	}*/
		



	
	Println(<-doneChan)
}