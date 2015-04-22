package main

import (
	."fmt"
	"finiteStateMachine"
	"network"
	"timer"
	"driver"
	"os"
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
	resetStuckTimerChan		:= make(chan string)
	enableStuckTimerChan	:= make(chan int)
	stuckTimeoutChan		:= make(chan int)
	
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
	timer.Init()
		
	
	go driver.UpButtonPoller(upButtonChan)
	go driver.DownButtonPoller(downButtonChan)
	go driver.CommandButtonPoller(commandButtonChan)
	go driver.FloorPoller(floorChan)
	
	go finiteStateMachine.Initialize(initChan, floorChan)
	Println(<-initChan)
	
	go network.SendManager(sendChan, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan)
	go network.Receive(receiveChan)
	go network.SendAliveMessage()
	go network.ReceiveAliveMessage(receiveAliveMessageChan)
	
	go timer.ElevatorTimer(elevatorTimerChan, elevatorTimeoutChan)
	go timer.DoorTimer(doorTimerChan, doorTimeoutChan)
	go timer.StuckTimer(resetStuckTimerChan, enableStuckTimerChan, stuckTimeoutChan)
	
	
	go finiteStateMachine.CheckAliveElevators(receiveAliveMessageChan, elevatorTimerChan, sendChan)
	go finiteStateMachine.EventHandler(sendChan, upButtonChan, downButtonChan, commandButtonChan, floorChan, receiveChan, doorTimerChan, resetStuckTimerChan, enableStuckTimerChan)			
	go finiteStateMachine.CheckIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAdded(sendChan, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan)
	
	Println(<-doneChan)
}













