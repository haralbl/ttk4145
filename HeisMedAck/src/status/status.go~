package status

import (
	."fmt"
	"network"
	//"math"
	"driver"
	//"math/rand"
	"encoding/json"
	//"time"
	//"./Timer"
	"structDefine"

)

const (
	numberOfElevators	= 3
	numberOfFloors		= 4
	
	IDLE				= 0
	DOOR_OPEN			= 1
	MOVING				= 2
	
	UP 					= 1
	DOWN 				= -1
	STOP				= 0
)

var ElevatorStatus structDefine.ElevatorStatus_t
/*
type elevatorStatus_t struct{
	ActiveElevators		[numberOfElevators]string // IP addresses
 
	PreviousFloors		[numberOfElevators]int
	InFloor				[numberOfElevators]int
	Directions			[numberOfElevators]int

	OrdersUp			[numberOfElevators][numberOfFloors]int
	OrdersDown			[numberOfElevators][numberOfFloors]int
	OrdersOut			[numberOfElevators][numberOfFloors]int

	State				int
	
	MessageType			string
	OrderedButtonType	int
	OrderedElevator		string
	OrderedFloor		int	
}
*/
func Initialize(initChan chan string, floorChan chan int) {	
	for i:=0; i<numberOfElevators; i++ {
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
				PrintStatus(ElevatorStatus)
				return
				Println("init goFunc not exited")
			}
		}
	}
	
	
}

func PrintStatus(status structDefine.ElevatorStatus_t) {
	Printf("Active elevators: %s, %s, %s\n", status.ActiveElevators[0],status.ActiveElevators[1],status.ActiveElevators[2])
	Printf("PreviousFloors: %d, %d, %d\n", status.PreviousFloors[0],status.PreviousFloors[1],status.PreviousFloors[2])
	Printf("InFloor: %d, %d, %d\n", status.InFloor[0],status.InFloor[1],status.InFloor[2])
	Printf("Directions: %d, %d, %d\n", status.Directions[0],status.Directions[1],status.Directions[2])
	Printf("\nElevator0\nOrdersUp	OrdersDown	OrdersOut\n")
	for i:=0; i<numberOfFloors; i++ {
		Printf("%d		%d		%d\n", status.OrdersUp[0][i],status.OrdersDown[0][i],status.OrdersOut[0][i])
	}
	Printf("\nElevator1\nOrdersUp	OrdersDown	OrdersOut\n")
	for i:=0; i<numberOfFloors; i++ {
		Printf("%d		%d		%d\n", status.OrdersUp[1][i],status.OrdersDown[1][i],status.OrdersOut[1][i])	
	}
	Printf("\nElevator2\nOrdersUp	OrdersDown	OrdersOut\n")
	for i:=0; i<numberOfFloors; i++ {
		Printf("%d		%d		%d\n", status.OrdersUp[2][i],status.OrdersDown[2][i],status.OrdersOut[2][i])	
	}	
	Printf("\nStatus: %d", status.State)				
	
	Printf("\nMessageType: %s", status.MessageType)
	Printf("\nOrderedButtonType: %d", status.OrderedButtonType)	
	Printf("\nOrderedElevator: %s", status.OrderedElevator)		
	Printf("\nOrderedFloor: %d\n\n", status.OrderedFloor)
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

	//PrintStatus(ElevatorStatus)
	
	// update status
	currentIPtoUpdate		:= ""
	currentPositionInReceivedStatus := 0
	
	for i:=1; i<numberOfElevators; i++ {
		currentIPtoUpdate = ElevatorStatus.ActiveElevators[i]
		currentPositionInReceivedStatus = -1
		for j:=0; j<numberOfElevators; j++ {
			if currentIPtoUpdate == receivedStatus.ActiveElevators[j] {
				currentPositionInReceivedStatus = j
			}
		}
		if currentPositionInReceivedStatus == -1 {
			Println("received IP in message that i dont have myself")
		} else {
			ElevatorStatus.PreviousFloors[i] = receivedStatus.PreviousFloors[currentPositionInReceivedStatus]
			ElevatorStatus.InFloor[i] = receivedStatus.InFloor[currentPositionInReceivedStatus]
			ElevatorStatus.Directions[i] = receivedStatus.Directions[currentPositionInReceivedStatus]
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
	
	//Printf("msgtype: %s\nelevator: %s\nfloor: %d\nbuttonType: %d\n", MessageType, elevator, floor, buttonType)
	
	
	/*
	// update status
	currentIPtoUpdate		:= ""
	currentPositionInReceivedStatus := 0
	
	for i:=1; i<numberOfElevators; i++ {
		currentIPtoUpdate = receivedStatus.ActiveElevators[i]
		currentPositionInReceivedStatus = -1
		for j:=0; j<numberOfElevators; j++ {
			if currentIPtoUpdate == receivedStatus.ActiveElevators[j] {
				currentPositionInReceivedStatus = j
			}
		}
		if currentPositionInReceivedStatus == -1 {
			Println("received IP in message that i dont have myself")
		} else {
			//ElevatorStatus.PreviousFloors[i] = receivedStatus.PreviousFloors[currentPositionInReceivedStatus]
			//ElevatorStatus.InFloor[i] = receivedStatus.InFloor[currentPositionInReceivedStatus]
			//ElevatorStatus.Directions[i] = receivedStatus.Directions[currentPositionInReceivedStatus]
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
	*/
	
	
	
	
	return
}

func CheckIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAdded(sendChan chan []byte, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan chan []byte) {
	var data []byte
	for {
		select {
		case data = <- checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan:
			var tempStatus structDefine.ElevatorStatus_t
			Println(len(data))	
			
			err := json.Unmarshal(data[0:len(data)], &tempStatus)
			if err != nil {
				Println(err)
			}
			
			if tempStatus.MessageType == "newOrder" {
				var elevator int = -1
				for i:=0; i<numberOfElevators; i++ {
					if tempStatus.OrderedElevator == ElevatorStatus.ActiveElevators[i] {
						elevator = i
					}
				}
				buttonType 	:= tempStatus.OrderedButtonType
				floor 		:= tempStatus.OrderedFloor
				orderNotAddedFlag:= 0
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

//flytte til network?
func handleMessage(sendChan chan []byte, /*ackResetChan chan string,*/ doorTimerChan chan string, elevatorIP string, floor int, buttonType int, MessageType string){
	switch (MessageType) {
	case "":
		Println("json is eating the message, and shitting out an empty status")
	case "ack":
		//tenne lys!!!
		Println("received ack")
		
		var elevator int = -1
		for i:=0; i<numberOfElevators; i++ {
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
		
		// legge til ordre du ikke selv skal ta
		//sette knappelys hvis ikke du tar ordren selv
		
	case "newOrder":
		if elevatorIP == ElevatorStatus.ActiveElevators[0] {
			Println("received new order to handle myself")
			
			driver.Set_button_lamp(buttonType, floor, 1)
		
			sendChan <- wrapMessage("ack", buttonType, elevatorIP, floor)
		
			switch (buttonType) {
			case 0:
				if ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
					//sendChan <- wrapMessage("ack", buttonType, elevatorIP, floor) //denne må vekk
					ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor] = 1
					event_newOrder(sendChan, doorTimerChan)
				}
			case 1:
				if ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
					//sendChan <- wrapMessage("ack", buttonType, elevatorIP, floor) 
					ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] = 1
					event_newOrder(sendChan, doorTimerChan)
				}
			case 2:
				if ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
					//sendChan <- wrapMessage("ack", buttonType, elevatorIP, floor)
					ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor] = 1
					event_newOrder(sendChan, doorTimerChan)
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
		driver.Set_button_lamp(2, floor, 0)
	case "updateAwokenElevator":
		Println("received updateAwokenElevator")
	}
}

func elevatorIPtoIndex(elevatorIP string) (elevatorIndex int) {
	elevatorIndex = 0
	for i:=0; i<numberOfElevators; i++ {
		if elevatorIP == ElevatorStatus.ActiveElevators[i] {
			elevatorIndex = i
		}
	}
	return
}

func isElevatorInList(elevatorIP string) int {
	for i:=0; i<numberOfElevators; i++ {
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
	for i:=numberOfElevators-1; i>-1; i-- {
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
			PrintStatus(ElevatorStatus)
		}
	}
}
func costFunction(floor int, buttonType int) int {
	floor = floor
	buttonType = buttonType
	cheapestElevator := 0
	distanceFromTarget := numberOfFloors
	for i := 0; i < numberOfElevators; i++ {
		if ElevatorStatus.ActiveElevators[i] != "empty" {
			if floor >= ElevatorStatus.PreviousFloors[i] {
				if distanceFromTarget > floor-ElevatorStatus.PreviousFloors[i] {
					distanceFromTarget = floor-ElevatorStatus.PreviousFloors[i] //potensiell sanntidsfeil her hvis prevfloors oppdateres mens funksjonen kjøres
					cheapestElevator = i
				}
			}else {
				if distanceFromTarget > ElevatorStatus.PreviousFloors[i]-floor {
					distanceFromTarget = ElevatorStatus.PreviousFloors[i]-floor //potensiell sanntidsfeil her hvis prevfloors oppdateres mens funksjonen kjøres
					cheapestElevator = i
				}
			}
		}
	}
	return cheapestElevator 
	//return rand.Intn(3)//ElevatorStatus.ActiveElevators[1]
}
/*func costFunction(floor int, buttonType int) (cheapestElevator int) {
	var costs[numberOfElevators]int
	for i:=0; i<numberOfElevators; i++ {
		costs[i] = 0
	}
	
	for i:=0; i<numberOfElevators; i++ {
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
	cheapestCost		:= math.Inf(1)
	
	for i:=0; i<numberOfElevators; i++ {
		if costs[i] < int(cheapestCost) {
			cheapestElevator = i
		}
	}
	return
}*/

// flytte til fsm?
func EventHandler(sendChan chan []byte, upButtonChan chan int, downButtonChan chan int,
					commandButtonChan chan int, floorChan chan int, /*ackTimerChan chan string,*/
					receiveChan chan []byte, /*ackTimeoutChan chan string, ackResetChan chan string,*/
					doorTimerChan chan string) {
	var button			int
	//var currentFloor	int
	var chosenElevator	string
	var message			[]byte
	for {
		select {
		case button = <- upButtonChan:
			
			//PrintStatus(ElevatorStatus)
			
			chosenElevator = ElevatorStatus.ActiveElevators[costFunction(button, 0)]
			//ackTimerChan <- "acktimer"
			sendChan <- wrapMessage("newOrder", 0, chosenElevator, button)
			
		case button = <- downButtonChan:
			
			//PrintStatus(ElevatorStatus)
			
			chosenElevator = ElevatorStatus.ActiveElevators[costFunction(button, 1)]
			//ackTimerChan <- "acktimer"
			sendChan <- wrapMessage("newOrder", 1, chosenElevator, button)
			
		case button = <- commandButtonChan:
			
			//PrintStatus(ElevatorStatus)
			
			chosenElevator = ElevatorStatus.ActiveElevators[0]
			sendChan <- wrapMessage("newOrder", 2, chosenElevator, button)
			

		case message = <- receiveChan:
			PrintStatus(ElevatorStatus)
			
			elevator, floor, buttonType, MessageType := unwrapMessage(message)
			
			handleMessage(sendChan, /*ackResetChan,*/ doorTimerChan, elevator, floor, buttonType, MessageType)
			
		case ElevatorStatus.PreviousFloors[0] = <- floorChan:
			StateOfShouldStop := shouldStop()
			//if StateOfShouldStop != 1 {
			sendChan <- wrapMessage("floorReached", 0, "", ElevatorStatus.PreviousFloors[0])
			//}
			event_floorReached(StateOfShouldStop, doorTimerChan)
			
		case <- doorTimerChan:
			Printf("door timer finished\n")
			sendChan <- wrapMessage("orderCompleted", 0, "", ElevatorStatus.PreviousFloors[0])
			event_doorTimerOut()
		}
	}
}
	
func event_newOrder(sendChan chan []byte, doorTimerChan chan string) {
	switch (ElevatorStatus.State) {
	case IDLE:
		//set button lamp
		if nextDirection() == STOP {
			//State = IDLE;
			//Print("State = IDLE\n")
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
		//set button lamp
	case MOVING:
		//set button lamp
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

func event_doorTimerOut() {
	switch ElevatorStatus.State {
	case IDLE:
		// close door? (shouldnt happen)
	case DOOR_OPEN:
		driver.Set_door_open_lamp(0)
		if nextDirection() == STOP {  
			ElevatorStatus.State = IDLE;
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
		// close door? (shouldnt happen)
	}
}

func shouldStop() int {
	if ElevatorStatus.Directions[0] == UP {
		Printf("OrdersUp: %d, OrdersDown: %d \n", ElevatorStatus.OrdersUp[0][ElevatorStatus.PreviousFloors[0]], ElevatorStatus.OrdersDown[0][ElevatorStatus.PreviousFloors[0]])
		if (ElevatorStatus.OrdersUp[0][ElevatorStatus.PreviousFloors[0]] | ElevatorStatus.OrdersOut[0][ElevatorStatus.PreviousFloors[0]]) != 0 {
			return 1;
		}
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<numberOfFloors; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return 0;
			}
		}
	}
	if ElevatorStatus.Directions[0] == DOWN {
		if(ElevatorStatus.OrdersDown[0][ElevatorStatus.PreviousFloors[0]] | ElevatorStatus.OrdersOut[0][ElevatorStatus.PreviousFloors[0]])!=0 {
			return 1;
		}
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if(ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0{
				return 0;
			}
		}
	}
	return 1;
}

func nextDirection() int {
	if ElevatorStatus.PreviousFloors[0] == 0 {
	
		Println("YOLO etasje 0")
	
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return UP;
			}
		}
	} else if ElevatorStatus.PreviousFloors[0] == numberOfFloors {
	
		Println("YOLO etasje 3")
	
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return DOWN;
 			}
		}
	} else if (ElevatorStatus.Directions[0] == UP){
	
		Println("YOLO retning opp")
	
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return UP;
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
				return DOWN;
			}
		}
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return UP;
			}
		}
	}
	return STOP;
}



















