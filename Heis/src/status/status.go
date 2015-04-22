package status

import (
	."fmt"
	"network"
	"math"
	"driver"
	"encoding/json"
	"structDefine"
)

const (
	NumberOfElevators	= 3
	numberOfFloors		= 4
	
	IDLE				= 0
	DOOR_OPEN			= 1
	MOVING				= 2
	
	UP 					= 1
	DOWN 				= -1
	STOP				= 0
)

var (
	ElevatorStatus structDefine.ElevatorStatus_t
)

func Initialize(initChan chan string, floorChan chan int) {	
	for i:=0; i<NumberOfElevators; i++ {
		ElevatorStatus.ActiveElevators[i] = "empty"
	} 
	ElevatorStatus.ActiveElevators[0] = network.GetLocalIP()
	driver.Set_motor_direction(DOWN)
	
	for {
		select {
		case ElevatorStatus.PreviousFloors[0] = <- floorChan:
			driver.Set_floor_indicator(ElevatorStatus.PreviousFloors[0])
			if ElevatorStatus.PreviousFloors[0] == 0 {
				driver.Set_motor_direction(STOP)
				ElevatorStatus.State = IDLE
				initChan <- "Finished init"
				return
			}
		}
	}
}

func PrintStatus(status structDefine.ElevatorStatus_t) {
	Printf("Active elevators: %s, %s, %s\n", status.ActiveElevators[0],	status.ActiveElevators[1],	status.ActiveElevators[2])
	Printf("PreviousFloors:	  %d, %d, %d\n", status.PreviousFloors[0],	status.PreviousFloors[1],	status.PreviousFloors[2])
	Printf("InFloor:		  %d, %d, %d\n", status.InFloor[0],			status.InFloor[1],			status.InFloor[2])
	Printf("Directions:		  %d, %d, %d\n", status.Directions[0],		status.Directions[1],		status.Directions[2])
	
	Printf("\nElevator0\nOrdersUp	OrdersDown	OrdersOut\n")
	for i:=0; i<numberOfFloors; i++ {
		Printf("%d		%d		%d\n", status.OrdersUp[0][i], status.OrdersDown[0][i], status.OrdersOut[0][i])
	}
	
	Printf("\nElevator1\nOrdersUp	OrdersDown	OrdersOut\n")
	for i:=0; i<numberOfFloors; i++ {
		Printf("%d		%d		%d\n", status.OrdersUp[1][i], status.OrdersDown[1][i], status.OrdersOut[1][i])	
	}
	
	Printf("\nElevator2\nOrdersUp	OrdersDown	OrdersOut\n")
	for i:=0; i<numberOfFloors; i++ {
		Printf("%d		%d		%d\n", status.OrdersUp[2][i], status.OrdersDown[2][i], status.OrdersOut[2][i])	
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
	var receivedStatus structDefine.ElevatorStatus_t
	err := json.Unmarshal(message[0:getMessageLength(message)], &receivedStatus)
	if err != nil {
		Println(err)
	}
	
	currentIPtoUpdate				:= ""
	currentPositionInReceivedStatus := 0
	
	for i:=0; i<NumberOfElevators; i++ {
		currentIPtoUpdate = ElevatorStatus.ActiveElevators[i]
		currentPositionInReceivedStatus = -1
		for j:=0; j<NumberOfElevators; j++ {
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
			for j:=0; j<numberOfFloors; j++ {
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
			var tempStatus structDefine.ElevatorStatus_t
			err := json.Unmarshal(data[0:len(data)], &tempStatus)
			if err != nil {
				Println(err)
			}
			
			if tempStatus.MessageType == "newOrder" {
				var elevator int = -1
				for i:=0; i<NumberOfElevators; i++ {
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

func handleMessage(sendChan chan []byte, doorTimerChan chan string, elevatorIP string, floor int, buttonType int, MessageType string, enableStuckTimerChan chan int) {
	switch (MessageType) {
	case "":
		Println("json is eating the message, and shitting out an empty status")
	case "ack":
		Println("received ack")
		
		var elevator int = -1
		for i:=0; i<NumberOfElevators; i++ {
			if elevatorIP == ElevatorStatus.ActiveElevators[i] {
				elevator = i
			}
		}
		if elevator == -1 {
			Println("It appears to be a fuckup in handlemessage() > case ack. I will take it myself.")
			elevator = 0
		}
		
		switch(buttonType) {
		case 0:
			ElevatorStatus.OrdersUp[elevator][floor] = 1
			driver.Set_button_lamp(buttonType, floor, 1)
		case 1:
			ElevatorStatus.OrdersDown[elevator][floor] = 1
			driver.Set_button_lamp(buttonType, floor, 1)
		case 2:
			ElevatorStatus.OrdersOut[elevator][floor] = 1
		}
		
	case "newOrder":
		if elevatorIP == ElevatorStatus.ActiveElevators[0] {
			Println("received new order to handle myself")
			
			driver.Set_button_lamp(buttonType, floor, 1)
			sendChan <- wrapMessage("ack", buttonType, elevatorIP, floor)
		
			switch (buttonType) {
			case 0:
				if ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
					ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor] = 1
					event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan)
				}
			case 1:
				if ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
					ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] = 1
					event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan)
				}
			case 2:
				if ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
					ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor] = 1
					event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan)
				}
			}
		}
		
	case "floorReached":
		Println("received floorReached")
		//ElevatorStatus.PreviousFloors[elevatorIPtoIndex(elevatorIP)] = floor
		
		//ElevatorStatus.PreviousFloors[elevatorIPtoIndex(elevatorIP)] = floor
		//ElevatorStatus.InFloor[elevatorIPtoIndex(elevatorIP)] = receivedStatus.InFloor[currentPositionInReceivedStatus]
		//ElevatorStatus.Directions[elevatorIPtoIndex(elevatorIP)] = 
	case "orderCompleted":
		Println("received orderCompleted")
		ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor]	= 0
		ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] = 0
		ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor]	= 0
		driver.Set_button_lamp(0, floor, 0)
		driver.Set_button_lamp(1, floor, 0)
		if ElevatorStatus.ActiveElevators[0] == elevatorIP {
			driver.Set_button_lamp(2, floor, 0)
		}
	case "updateAwokenElevator":
		Println("received updateAwokenElevator")
		for floor:=0; floor<numberOfFloors; floor++ {
			if ElevatorStatus.OrdersUp[0][floor] == 1 {
				driver.Set_button_lamp(0, floor, 1)
				event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan)
			}
			if ElevatorStatus.OrdersDown[0][floor] == 1 {
				driver.Set_button_lamp(1, floor, 1)
				event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan)
			}
			if ElevatorStatus.OrdersOut[0][floor] == 1 {
				driver.Set_button_lamp(2, floor, 1)
				event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan)
			}
		}
	}
}

func elevatorIPtoIndex(elevatorIP string) (elevatorIndex int) {
	elevatorIndex = 0
	for i:=0; i<NumberOfElevators; i++ {
		if elevatorIP == ElevatorStatus.ActiveElevators[i] {
			elevatorIndex = i
		}
	}
	return
}

func isElevatorInList(elevatorIP string) int {
	for i:=0; i<NumberOfElevators; i++ {
		if ElevatorStatus.ActiveElevators[i] == elevatorIP {
			return i
		}
	}
	return -1
}

func addElevator(elevatorIP string) int {
	alreadyAdded	:= false
	full			:= true
	nextIndex		:= 0
	for i:=NumberOfElevators-1; i>-1; i-- {
		if ElevatorStatus.ActiveElevators[i] == elevatorIP {
			alreadyAdded = true
		}
		if ElevatorStatus.ActiveElevators[i] == "empty" {
			full = false
			nextIndex = i
		}
	}
	if !alreadyAdded && !full {
		ElevatorStatus.ActiveElevators[nextIndex] = elevatorIP
		return nextIndex
	}
	return -1
}		

func removeElevator(elevatorN int) {
	ElevatorStatus.ActiveElevators[elevatorN] = "empty"
}

func CheckAliveElevators(receiveAliveMessageChan chan string, elevatorTimerChan chan int, sendChan chan []byte) {
	var elevatorIP	string
	var elevatorN	int
	for {
		select {
		case elevatorIP = <- receiveAliveMessageChan:
			elevatorIP = elevatorIP[0:15]
			elevatorN = isElevatorInList(elevatorIP)
			if elevatorN == -1 {
				elevatorN = addElevator(elevatorIP)
				sendChan <- wrapMessage("updateAwokenElevator", 0, "", 0)
			}
			elevatorTimerChan <- elevatorN
			
		case elevatorN = <- elevatorTimerChan:
			removeElevator(elevatorN)
			var lowestIPindex = 0
			var lowestIP = ElevatorStatus.ActiveElevators[0]
			for i:=0; i<NumberOfElevators; i++ {
				if ElevatorStatus.ActiveElevators[i] < lowestIP {
					lowestIP = ElevatorStatus.ActiveElevators[i]
					lowestIPindex = i
				}
			}	
			for i:=0; i<numberOfFloors; i++ {
				if ElevatorStatus.OrdersUp[elevatorN][i] == 1	{
					if lowestIPindex == 0 {
						sendChan <- wrapMessage("newOrder", 0, lowestIP, i)
					}
					ElevatorStatus.OrdersUp[elevatorN][i] = 0
				}
				if ElevatorStatus.OrdersDown[elevatorN][i] == 1	{
					if lowestIPindex == 0 {
						sendChan <- wrapMessage("newOrder", 1, lowestIP, i)
					}
					ElevatorStatus.OrdersDown[elevatorN][i] = 0
				}
			}
		}
	}
}

func costFunction(floor int, buttonType int) (cheapestElevator int) {
	var costs[NumberOfElevators]int
	for i:=0; i<NumberOfElevators; i++ {
		costs[i] = 0
	}
	
	for i:=0; i<NumberOfElevators; i++ {
		
		// Check number of orders
		for j:=0; j<numberOfFloors; j++ {
			if ElevatorStatus.OrdersUp[i][j] == 1 || ElevatorStatus.OrdersDown[i][j] == 1 || ElevatorStatus.OrdersOut[i][j] == 1 {
				costs[i] += 2
			}
		}
		// Check distances
		costs[i] += int(math.Abs(float64(floor) - float64(ElevatorStatus.PreviousFloors[i])))
		
		// Check if order in same direction in front of elevator
		if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == UP && buttonType == 0 {
			costs[i] += 2*numberOfFloors
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == DOWN && buttonType == 1 {
			costs[i] += 2*numberOfFloors
			
		// Check if order in opposite direction in front of elevator	
		} else if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == UP && buttonType == 1 {
			costs[i] += 5*numberOfFloors
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == DOWN && buttonType == 0 {
			costs[i] += 5*numberOfFloors

		// Check if order in opposite direction behind elevator
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == UP && buttonType == 1 {
			costs[i] += 8*numberOfFloors
		} else if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == DOWN && buttonType == 0 {
			costs[i] += 8*numberOfFloors

		// Check if order in same direction behind elevator
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == UP && buttonType == 0 {
			costs[i] += 11*numberOfFloors
		} else if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == DOWN && buttonType == 1 {
			costs[i] += 11*numberOfFloors
		}
	}
	cheapestElevator	= 0
	cheapestCost		:= 10000
	
	for i:=0; i<NumberOfElevators; i++ {
		if costs[i] < int(cheapestCost) && ElevatorStatus.ActiveElevators[i] != "empty" {
			cheapestCost = costs[i]
			cheapestElevator = i
		}
	}
	Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nkostnader liste")
	Println(costs)
	Println(cheapestElevator)
	Printf("======================================================\n\n\n\n\n\n\n\n\n")
	return
}

/*func costFunction(floor int, buttonType int) (cheapestElevator int) {
	var costs[NumberOfElevators]int
	for i:=0; i<NumberOfElevators; i++ {
		costs[i] = 0
	}
	
	for i:=0; i<NumberOfElevators; i++ {
		for j:=0; j<numberOfFloors; j++ {
			// Check number of orders
			if ElevatorStatus.OrdersUp[i][j] == 1 || ElevatorStatus.OrdersDown[i][j] == 1 || ElevatorStatus.OrdersOut[i][j] == 1 {
				costs[i] += 10
			}
		}
		// Check if direction towards order
		if (floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == 1) || (floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == 0) {
			costs[i] += 5*numberOfFloors
		}
		// Check distances
		costs[i] += int(2*math.Abs(float64(floor) - float64(ElevatorStatus.PreviousFloors[i])))
		// Check if same direction as order
		if buttonType != ElevatorStatus.Directions[i] {
			costs[i] += 5*numberOfFloors
		}
	}
	cheapestElevator	= 0
	cheapestCost		:= 10000
	
	for i:=0; i<NumberOfElevators; i++ {
		if costs[i] < int(cheapestCost) && ElevatorStatus.ActiveElevators[i] != "empty" {
			cheapestCost = costs[i]
			cheapestElevator = i
		}
	}
	Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nkostnader liste")
	Println(costs)
	Println(cheapestElevator)
	Printf("======================================================\n\n\n\n\n\n\n\n\n")
	return
}*/

func EventHandler(sendChan chan []byte, upButtonChan chan int, downButtonChan chan int,
					commandButtonChan chan int, floorChan chan int, receiveChan chan []byte,
					doorTimerChan chan string, resetStuckTimerChan chan string, enableStuckTimerChan chan int) {
	var button			int
	var chosenElevator	string
	var message			[]byte
	for {
		select {
		case button = <- upButtonChan:
			chosenElevator = ElevatorStatus.ActiveElevators[costFunction(button, 0)]
			sendChan <- wrapMessage("newOrder", 0, chosenElevator, button)
			
		case button = <- downButtonChan:
			chosenElevator = ElevatorStatus.ActiveElevators[costFunction(button, 1)]
			sendChan <- wrapMessage("newOrder", 1, chosenElevator, button)
			
		case button = <- commandButtonChan:
			chosenElevator = ElevatorStatus.ActiveElevators[0]
			sendChan <- wrapMessage("newOrder", 2, chosenElevator, button)
			

		case message = <- receiveChan:
			elevator, floor, buttonType, MessageType := unwrapMessage(message)
			handleMessage(sendChan, doorTimerChan, elevator, floor, buttonType, MessageType, enableStuckTimerChan)
			
		case ElevatorStatus.PreviousFloors[0] = <- floorChan:
			StateOfShouldStop := shouldStop()
			sendChan <- wrapMessage("floorReached", 0, "", ElevatorStatus.PreviousFloors[0])
			resetStuckTimerChan <- "reset"
			event_floorReached(StateOfShouldStop, doorTimerChan)
			
		case <- doorTimerChan:
			Printf("door timer finished\n")
			sendChan <- wrapMessage("orderCompleted", 0, ElevatorStatus.ActiveElevators[0], ElevatorStatus.PreviousFloors[0])
			event_doorTimerOut(enableStuckTimerChan)
		}
	}
}
	
func event_newOrder(sendChan chan []byte, doorTimerChan chan string, enableStuckTimerChan chan int) {
	switch (ElevatorStatus.State) {
	case IDLE:
		enableStuckTimerChan <- 1
		if nextDirection() == STOP {
			driver.Set_door_open_lamp(1)
			doorTimerChan <- "start"
			ElevatorStatus.State = DOOR_OPEN
			Print("State = DOOR_OPEN\n")
			sendChan <- wrapMessage("orderCompleted", 0, "", ElevatorStatus.PreviousFloors[0])
		} else if nextDirection() == UP {
			driver.Set_motor_direction(UP)
			ElevatorStatus.Directions[0] = UP
 			ElevatorStatus.State = MOVING 			
 			Print("State = MOVING\n")
		} else if nextDirection() == DOWN {
			driver.Set_motor_direction(DOWN)
			ElevatorStatus.Directions[0] = DOWN
			ElevatorStatus.State = MOVING
			Print("State = MOVING\n")
		} else {
			Printf("ERROR, event_newOrder: nextDirection returns invalid value")
		}
	case DOOR_OPEN:
		// do nothing
	case MOVING:
		// do nothing
	}
}

func event_floorReached(StateOfShouldStop int, doorTimerChan chan string) {
	driver.Set_floor_indicator(ElevatorStatus.PreviousFloors[0])
	switch ElevatorStatus.State {
	case IDLE:
		Printf("ERROR, event: floorReached when not moving!")
	case DOOR_OPEN:
		Printf("ERROR, event: floorReached when not moving!")
	case MOVING:
		if StateOfShouldStop == 1 {
			driver.Set_motor_direction(STOP)
			driver.Set_door_open_lamp(1)
			doorTimerChan <- "start"
			ElevatorStatus.State = DOOR_OPEN
			Print("State = DOOR_OPEN\n")
		}
	}
}

func event_doorTimerOut(enableStuckTimerChan chan int) {
	switch ElevatorStatus.State {
	case IDLE:
		// do nothing
	case DOOR_OPEN:
		driver.Set_door_open_lamp(0)
		if nextDirection() == STOP {  
			ElevatorStatus.State = IDLE
			enableStuckTimerChan <- 0
			Print("State = IDLE\n")
		} else if nextDirection() == UP {
			driver.Set_motor_direction(UP)
			ElevatorStatus.Directions[0] = UP
 			ElevatorStatus.State = MOVING
 			Print("State = MOVING\n")
		} else if nextDirection() == DOWN {
			driver.Set_motor_direction(DOWN)
			ElevatorStatus.Directions[0] = DOWN
			ElevatorStatus.State = MOVING
			Print("State = MOVING\n")
		} else {
			Printf("ERROR, event_doorTimerOut: nextDirection returns invalid value")
		}
	case MOVING:
		// do nothing
	}
}

func shouldStop() int {
	if ElevatorStatus.Directions[0] == UP {
		if (ElevatorStatus.OrdersUp[0][ElevatorStatus.PreviousFloors[0]] | ElevatorStatus.OrdersOut[0][ElevatorStatus.PreviousFloors[0]]) != 0 {
			return 1
		}
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<numberOfFloors; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return 0
			}
		}
	}
	if ElevatorStatus.Directions[0] == DOWN {
		if(ElevatorStatus.OrdersDown[0][ElevatorStatus.PreviousFloors[0]] | ElevatorStatus.OrdersOut[0][ElevatorStatus.PreviousFloors[0]])!=0 {
			return 1
		}
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if(ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0{
				return 0
			}
		}
	}
	return 1
}

func nextDirection() int {
	if ElevatorStatus.PreviousFloors[0] == 0 {
	
		Println("YOLO etasje 0")
	
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return UP
			}
		}
	} else if ElevatorStatus.PreviousFloors[0] == numberOfFloors {
	
		Println("YOLO etasje 3")
	
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return DOWN
 			}
		}
	} else if (ElevatorStatus.Directions[0] == UP){
	
		Println("YOLO retning opp")
	
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return UP
			}
		}
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return DOWN;
			}
		}
	} else if (ElevatorStatus.Directions[0] == DOWN){
	
		Println("YOLO retning ned")
	
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return DOWN
			}
		}
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return UP
			}
		}
	}
	return STOP
}















