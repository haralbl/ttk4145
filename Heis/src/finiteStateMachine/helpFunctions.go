package finiteStateMachine

import (
	."fmt"
	"network"
	"encoding/json"
	"defines"
	"driver"
)

func Initialize(initChan chan string, floorReachedChan chan int, floorLeftChan chan string) {	
	for i:=0; i<defines.NumberOfElevators; i++ {
		ElevatorStatus.ActiveElevators[i] = "empty"
	} 
	ElevatorStatus.ActiveElevators[0] = network.GetLocalIP()
	driver.Set_motor_direction(defines.DOWN)
	
	for {
		select {
		case <- floorLeftChan: //to prevent deadlock
		
		case ElevatorStatus.PreviousFloors[0] = <- floorReachedChan:
			driver.Set_floor_indicator(ElevatorStatus.PreviousFloors[0])
			if ElevatorStatus.PreviousFloors[0] == 0 {
				driver.Set_motor_direction(defines.STOP)
				ElevatorStatus.State = defines.IDLE
				initChan <- "Finished init"
				
				PrintStatus(ElevatorStatus)
				
				return
			}
		}
	}
}

func PrintStatus(status defines.ElevatorStatus_t) {
	Printf("Active elevators: %s, %s, %s\n", status.ActiveElevators[0],	status.ActiveElevators[1],	status.ActiveElevators[2])
	Printf("PreviousFloors:	  %d, %d, %d\n", status.PreviousFloors[0],	status.PreviousFloors[1],	status.PreviousFloors[2])
	Printf("InFloor:		  %d, %d, %d\n", status.InFloor[0],			status.InFloor[1],			status.InFloor[2])
	Printf("Directions:		  %d, %d, %d\n", status.Directions[0],		status.Directions[1],		status.Directions[2])
	
	for e:=0; e<defines.NumberOfElevators; e++ {
		Printf("\nElevator %d\nOrdersUp	OrdersDown	OrdersOut\n", e)
		for i:=0; i<defines.NumberOfFloors; i++ {
			Printf("%d		%d		%d\n", status.OrdersUp[e][i], status.OrdersDown[e][i], 
			status.OrdersOut[e][i])
		}
	}
		
	Printf("\nStatus:				%d",	 status.State)				
	Printf("\nMessageType:			%s",	 status.MessageType)
	Printf("\nOrderedButtonType:	%d",	 status.OrderedButtonType)	
	Printf("\nOrderedElevator:		%s", 	 status.OrderedElevator)		
	Printf("\nOrderedFloor:			%d\n\n", status.OrderedFloor)
}


func wrapMessage(newMessageType string, buttonType int, elevator string, floor int) []byte {
	ElevatorStatus.MessageType			= newMessageType
	ElevatorStatus.OrderedButtonType	= buttonType
	ElevatorStatus.OrderedElevator		= elevator
	ElevatorStatus.OrderedFloor			= floor
	
	message := make([]byte, 1024)
	message,_ = json.Marshal(ElevatorStatus)
	
	return []byte(message)
}

func getMessageLength(message []byte) int {
	return int(message[1000]) + int(message[1001])*10 + int(message[1002])*100
}

func unwrapMessage(message []byte) (elevator string, floor int, buttonType int, MessageType string) {
	var receivedStatus defines.ElevatorStatus_t
	err := json.Unmarshal(message[0:getMessageLength(message)], &receivedStatus)
	if err != nil {
		Println(err)
	}
	
	currentIPtoUpdate				:= ""
	currentPositionInReceivedStatus := 0
	
	for i:=0; i<defines.NumberOfElevators; i++ {
		currentIPtoUpdate = ElevatorStatus.ActiveElevators[i]
		currentPositionInReceivedStatus = -1
		for j:=0; j<defines.NumberOfElevators; j++ {
			if currentIPtoUpdate == receivedStatus.ActiveElevators[j] {
				currentPositionInReceivedStatus = j
			}
		}
		if currentPositionInReceivedStatus == -1 {
			Println("received IP in message that i dont have myself")

		} else {
			if receivedStatus.ActiveElevators[currentPositionInReceivedStatus] != ElevatorStatus.ActiveElevators[0]{
				ElevatorStatus.PreviousFloors[i] = receivedStatus.PreviousFloors[currentPositionInReceivedStatus]
				ElevatorStatus.InFloor[i] = receivedStatus.InFloor[currentPositionInReceivedStatus]
				ElevatorStatus.Directions[i] = receivedStatus.Directions[currentPositionInReceivedStatus]
			}
			for j:=0; j<defines.NumberOfFloors; j++ {
				ElevatorStatus.OrdersUp[i][j]	= ElevatorStatus.OrdersUp[i][j]	| receivedStatus.OrdersUp[currentPositionInReceivedStatus][j]
				ElevatorStatus.OrdersDown[i][j] = ElevatorStatus.OrdersDown[i][j] | receivedStatus.OrdersDown[currentPositionInReceivedStatus][j]
				ElevatorStatus.OrdersOut[i][j]	= ElevatorStatus.OrdersOut[i][j] | receivedStatus.OrdersOut[currentPositionInReceivedStatus][j]
			}
		}
	}
	MessageType	= receivedStatus.MessageType
	elevator	= receivedStatus.OrderedElevator
	floor		= receivedStatus.OrderedFloor
	buttonType	= receivedStatus.OrderedButtonType
	
	return
}

func CheckIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAdded(sendChan chan []byte, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan chan []byte) {
	var data []byte
	for {
		select {
		case data = <- checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan:
			var tempStatus defines.ElevatorStatus_t
			err := json.Unmarshal(data[0:len(data)], &tempStatus)
			if err != nil {
				Println(err)
			}
			
			if tempStatus.MessageType == "newOrder" {
				var elevator int = -1
				for i:=0; i<defines.NumberOfElevators; i++ {
					if tempStatus.OrderedElevator == ElevatorStatus.ActiveElevators[i] {
						elevator = i
					}
				}
				buttonType 	:= tempStatus.OrderedButtonType
				floor 		:= tempStatus.OrderedFloor
				orderNotAddedFlag := 0
				
				if elevator != -1 {
					switch(buttonType) {
					case 0:
						if ElevatorStatus.OrdersUp[elevator][floor] != 1 {
							orderNotAddedFlag = 1
						}
					case 1:
						if ElevatorStatus.OrdersDown[elevator][floor] != 1 {
							orderNotAddedFlag = 1
						}
					case 2:
						if ElevatorStatus.OrdersOut[elevator][floor] != 1 {
							orderNotAddedFlag = 1
						}
					}
				} else {
					orderNotAddedFlag = 1
				}
				if orderNotAddedFlag == 1 {
					Println("her er jeg")
					sendChan <- wrapMessage("newOrder", buttonType, ElevatorStatus.ActiveElevators[0], floor)
				} 
			}
		}
	}
}

func elevatorIPtoIndex(elevatorIP string) (elevatorIndex int) {
	elevatorIndex = 0
	for i:=0; i<defines.NumberOfElevators; i++ {
		if elevatorIP == ElevatorStatus.ActiveElevators[i] {
			elevatorIndex = i
		}
	}
	return
}









