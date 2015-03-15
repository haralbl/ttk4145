package status

import (
	."fmt"
	"strings"
	"strconv"
	"network"
	//"math"
	"driver"
	"math/rand"
	//"time"
	//"./Timer"
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

var (
	activeElevators		[numberOfElevators]string // IP addresses
	
	previousFloors		[numberOfElevators]int
	inFloor				[numberOfElevators]int
	directions			[numberOfElevators]int
	
	ordersUp			[numberOfElevators][numberOfFloors]int
	ordersDown			[numberOfElevators][numberOfFloors]int
	ordersOut			[numberOfElevators][numberOfFloors]int
	
	state				int
)

func wrapMessage(messageType string, buttonType int, elevator int, floor int) (message string) {
	message = ""
	message += messageType + "\n"
	
	var newOrderUp		[numberOfElevators][numberOfFloors]int
	var newOrderDown	[numberOfElevators][numberOfFloors]int
	var newOrderOut		[numberOfElevators][numberOfFloors]int
	switch buttonType {
	case 0:
		newOrderUp[elevator][floor] = 1
	case 1:
		newOrderDown[elevator][floor] = 1
	case 2:
		newOrderOut[elevator][floor] = 1
	}
	for i:=0; i<numberOfElevators; i++ {
		for j:=0; j<numberOfFloors; j++ {
			message += strconv.Itoa(newOrderUp[i][j])
			message += " "
			message += strconv.Itoa(newOrderDown[i][j])
			message += " "
			message += strconv.Itoa(newOrderOut[i][j])
			message += " "
		}
		message += "\n"
	}
	
	for i:=0; i<numberOfElevators; i++ {
		message += activeElevators[i]
		message += " "
	}
	message += "\n"
	for i:=0; i<numberOfElevators; i++ {
		message += strconv.Itoa(previousFloors[i])
		message += " "
	}
	message += "\n"
	for i:=0; i<numberOfElevators; i++ {
		message += strconv.Itoa(inFloor[i])
		message += " "
	}
	message += "\n"
	for i:=0; i<numberOfElevators; i++ {
		message += strconv.Itoa(inFloor[i])
		message += " "
	}
	message += "\n"
	for i:=0; i<numberOfElevators; i++ {
		for j:=0; j<numberOfFloors; j++ {
			message += strconv.Itoa(ordersUp[i][j])
			message += " "
			message += strconv.Itoa(ordersDown[i][j])
			message += " "
			message += strconv.Itoa(ordersOut[i][j])
			message += " "
		}
		message += "\n"
	}
	return
}

func unwrapMessage(message string) (elevator int, floor int, buttonType int, order int, messageType string){
	statusFields := strings.Split(message, "\n")
	
	// update status
	temp := 0
	field := strings.Split(statusFields[5], " ")
	for i:=1; i<numberOfElevators; i++ {
		temp,_	= strconv.Atoi(field[i])
		previousFloors[i] = temp
	}
	field = strings.Split(statusFields[6], " ")
	for i:=1; i<numberOfElevators; i++ {
		temp,_		= strconv.Atoi(field[i])
		inFloor[i] = temp
	}
	field = strings.Split(statusFields[7], " ")
	for i:=1; i<numberOfElevators; i++ {
		temp,_		= strconv.Atoi(field[i])
		directions[i] = temp
	}
	for i:=0; i<numberOfElevators; i++ {
		field = strings.Split(statusFields[8+i], " ")
		for j:=0; j<numberOfFloors; j++ {
			temp,_	= strconv.Atoi(field[3*j+0])
			ordersUp[i][j] = ordersUp[i][j] | temp
			temp,_	= strconv.Atoi(field[3*j+1])
			ordersDown[i][j] = ordersDown[i][j] | temp
			temp,_	= strconv.Atoi(field[3*j+2])
			ordersOut[i][j] = ordersOut[i][j] | temp
		}
	}
	
	// add new info to status
	messageType = statusFields[0]
	
	for i:=0; i<numberOfElevators; i++ {
		field = strings.Split(statusFields[1+i], " ")
		for j:=0; j<numberOfFloors; j++ {
			order,_ = strconv.Atoi(field[3*j+0])
			if order == 1 {
				elevator = i
				floor = j
				buttonType = 0
			}
			order,_ = strconv.Atoi(field[3*j+1])
			if order == 1 {
				elevator = i
				floor = j
				buttonType = 1
			}
			order,_ = strconv.Atoi(field[3*j+2])
			if order == 1 {
				elevator = i
				floor = j
				buttonType = 2
			}
		}
	}
	
	// sjekk opp mot IP liste i message
	ipList := strings.Split(statusFields[4], " ")
	elevatorIP := ipList[elevator]
	for i:= 0; i<numberOfElevators; i++{
		if activeElevators[i] == elevatorIP {
			elevator = i
		}
	}
	return
}

//flytte til network?
func handleMessage(sendChan chan string, ackResetChan chan string, doorTimerChan chan string, elevator int, floor int, buttonType int, order int, messageType string){
	switch (messageType) {
	case "ack":
		Println("received ack")
		
		ackResetChan <- "Reset acktimer"
		// legge til ordre du ikke selv skal ta
		//sette knappelys hvis ikke du tar ordren selv
		
	case "newOrder":
		if elevator == 0{
			Println("received new order to handle myself")
			
			driver.Set_button_lamp(buttonType, floor, 1)
		
			sendChan <- wrapMessage("ack", buttonType, elevator, floor)
		
			switch (buttonType) {
			case 0:
				if ordersUp[elevator][floor] == 0 {
					sendChan <- wrapMessage("ack", buttonType, elevator, floor)
					ordersUp[elevator][floor] = 1
					event_newOrder(sendChan, doorTimerChan)
				}
			case 1:
				if ordersDown[elevator][floor] == 0 {
					sendChan <- wrapMessage("ack", buttonType, elevator, floor)
					ordersDown[elevator][floor] = 1
					event_newOrder(sendChan, doorTimerChan)
				}
			case 2:
				if ordersDown[elevator][floor] == 0 {
					sendChan <- wrapMessage("ack", buttonType, elevator, floor)
					ordersOut[elevator][floor] = 1
					event_newOrder(sendChan, doorTimerChan)
				}
			}
		}
		
	case "floorReached":
		Println("received floorReached")
		
	case "orderCompleted":
		Println("received orderCompleted")
		
		ordersUp[elevator][floor]		= 0
		ordersDown[elevator][floor] 	= 0
		ordersOut[elevator][floor]		= 0
		driver.Set_button_lamp(0, floor, 0)
		driver.Set_button_lamp(1, floor, 0)
		driver.Set_button_lamp(2, floor, 0)
	case "updateAwokenElevator":
		Println("received updateAwokenElevator")
	}
}

func Initialize(initChan chan string, floorChan chan int) {
	for i:=0; i<numberOfElevators; i++ {
		activeElevators[i] = "empty"
	} 
	activeElevators[0] = network.GetLocalIP()
	
	driver.Set_motor_direction(DOWN)
	for {
		select {
		case previousFloors[0] = <- floorChan:
			driver.Set_floor_indicator(previousFloors[0])
			if previousFloors[0] == 0 {
				driver.Set_motor_direction(STOP)
				initChan <- "Finished init"
				return
				Println("init goFunc not exited")
			}
		}
	}
	state = IDLE
}

func isElevatorInList(elevatorIP string) int {
	for i:=0; i<numberOfElevators; i++ {
		if activeElevators[i] == elevatorIP {
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
		if activeElevators[i] == elevatorIP {
			alreadyAdded = true
		}
		if activeElevators[i] == "empty" {
			full = false
			nextIndex = i
		}
	}
	if !alreadyAdded && !full {
		activeElevators[nextIndex] = elevatorIP
		return nextIndex
	}
	return -1
}		

func removeElevator(elevatorN int) {
	activeElevators[elevatorN] = "empty"
}

func CheckAliveElevators(receiveAliveMessageChan chan string, elevatorTimerChan chan int, sendChan chan string) {
	var elevatorIP	string
	var elevatorN	int
	for {
		select {
		case elevatorIP = <- receiveAliveMessageChan:
			elevatorIP = elevatorIP[0:15]
			elevatorN = isElevatorInList(elevatorIP)
			if elevatorN == -1 {
				elevatorN = addElevator(elevatorIP)
				sendChan <- wrapMessage("updateAwokenElevator", 0, 0, 0)
			}
			elevatorTimerChan <- elevatorN
		case elevatorN = <- elevatorTimerChan:
			removeElevator(elevatorN)
		}
	}
}
func costFunction(floor int, buttonType int) int {
	floor = floor
	buttonType = buttonType 
	return rand.Intn(1)
}/*
func costFunction(floor int, buttonType int) (cheapestElevator int) {
	var costs[numberOfElevators]int
	for i:=0; i<numberOfElevators; i++ {
		costs[i] = 0
	}
	
	for i:=0; i<numberOfElevators; i++ {
		for j:=0; j<numberOfFloors; j++ {
			// Check number of orders
			if 	ordersUp[i][j] == 1 || ordersDown[i][j] == 1 || ordersOut[i][j] == 1 {
				costs[i] += 10
			}
		}
		// Check if direction towards order
		if (floor > previousFloors[i] && directions[i] == 1) || (floor < previousFloors[i] && directions[i] == 0) {
			costs[i] += 5*numberOfFloors
		}
		// Check distances
		costs[i] += int(2*math.Abs(float64(floor) - float64(previousFloors[i])))
		// Check if same direction as order
		if buttonType != directions[i] {
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
func EventHandler(sendChan chan string, upButtonChan chan int, downButtonChan chan int,
					commandButtonChan chan int, floorChan chan int, ackTimerChan chan string,
					receiveChan chan string, ackTimeoutChan chan string, ackResetChan chan string,
					doorTimerChan chan string) {
	var ack				string
	var button			int
	//var currentFloor	int
	var chosenElevator	int
	var message			string
	for {
		select {
		case button = <- upButtonChan:
			chosenElevator = costFunction(button, 0)
			ackTimerChan <- "acktimer"
			sendChan <- wrapMessage("newOrder", 0, chosenElevator, button)
			
		case button = <- downButtonChan:
			chosenElevator = costFunction(button, 1)
			ackTimerChan <- "acktimer"
			sendChan <- wrapMessage("newOrder", 1, chosenElevator, button)
			
		case button = <- commandButtonChan:
			chosenElevator = 0
			sendChan <- wrapMessage("newOrder", 2, chosenElevator, button)
			
		case ack = <- ackTimerChan:
			ack = ack //hack
			// ta ordre selv
			
		case message = <- receiveChan:
			elevator, floor, buttonType, order, messageType := unwrapMessage(message)
			handleMessage(sendChan, ackResetChan, doorTimerChan, elevator, floor, buttonType, order, messageType)
			
		case previousFloors[0] = <- floorChan:
			stateOfShouldStop := shouldStop()
			if stateOfShouldStop == 1 {
				sendChan <- wrapMessage("orderCompleted", 0, 0, previousFloors[0])
			} else {
				sendChan <- wrapMessage("floorReached", 0, 0, previousFloors[0])
			}
			event_floorReached(stateOfShouldStop, doorTimerChan)
			
		case <- doorTimerChan:
			Printf("door timer finished\n")
			event_doorTimerOut()
		}
	}
}
	
func event_newOrder(sendChan chan string, doorTimerChan chan string) {
	switch (state) {
	case IDLE:
		//set button lamp
		if nextDirection() == STOP {  
			//state = IDLE;
			//Print("state = IDLE\n")
			driver.Set_door_open_lamp(1)
			doorTimerChan <- "start"
			state = DOOR_OPEN
			Print("state = DOOR_OPEN\n")
			sendChan <- wrapMessage("orderCompleted", 0, 0, previousFloors[0])
		} else if nextDirection() == UP {
			driver.Set_motor_direction(UP);
			directions[0] = UP;
 			state = MOVING;
 			Print("state = MOVING\n")
		} else if nextDirection() == DOWN {
			driver.Set_motor_direction(DOWN);
			directions[0] = DOWN;
			state = MOVING;
			Print("state = MOVING\n")
		} else {
			Printf("ERROR, event_newOrder: nextDirection returns invalid value")
		}
	case DOOR_OPEN:
		//set button lamp
	case MOVING:
		//set button lamp
	}
}

func event_floorReached(stateOfShouldStop int, doorTimerChan chan string) {
	driver.Set_floor_indicator(previousFloors[0])
	switch state {
	case IDLE:
		Printf("ERROR, event: floorReached when not moving!")
	case DOOR_OPEN:
		Printf("ERROR, event: floorReached when not moving!")
	case MOVING:
		if stateOfShouldStop == 1 {
			driver.Set_motor_direction(STOP)
			driver.Set_door_open_lamp(1)
			doorTimerChan <- "start"
			state = DOOR_OPEN
			Print("state = DOOR_OPEN\n")
		}
	}
}

func event_doorTimerOut() {
	switch state {
	case IDLE:
		// close door? (shouldnt happen)
	case DOOR_OPEN:
		driver.Set_door_open_lamp(0)
		if nextDirection() == STOP {  
			state = IDLE;
			Print("state = IDLE\n")
		} else if nextDirection() == UP {
			driver.Set_motor_direction(UP)
			directions[0] = UP
 			state = MOVING
 			Print("state = MOVING\n")
		} else if nextDirection() == DOWN {
			driver.Set_motor_direction(DOWN)
			directions[0] = DOWN
			state = MOVING
			Print("state = MOVING\n")
		} else {
			Printf("ERROR, event_doorTimerOut: nextDirection returns invalid value")
		}
	case MOVING:
		// close door? (shouldnt happen)
	}
}

func shouldStop() int {
	if directions[0] == UP {
		Printf("ordersUp: %d, OrdersDown: %d \n", ordersUp[0][previousFloors[0]], ordersDown[0][previousFloors[0]])
		if (ordersUp[0][previousFloors[0]] | ordersOut[0][previousFloors[0]]) != 0 {
			return 1;
		}
		for floor:=previousFloors[0]+1; floor<numberOfFloors; floor++ {
			if (ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0 {
				return 0;
			}
		}
	}
	if directions[0] == DOWN {
		if(ordersDown[0][previousFloors[0]] | ordersOut[0][previousFloors[0]])!=0 {
			return 1;
		}
		for floor:=0; floor<previousFloors[0]; floor++ {
			if(ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0{
				return 0;
			}
		}
	}
	return 1;
}

func nextDirection() int {
	if previousFloors[0] == 0 {
		for floor:=previousFloors[0]+1; floor<4; floor++ {
			if (ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0 {
				return UP;
			}
		}
	} else if previousFloors[0] == numberOfFloors {
		for floor:=0; floor<previousFloors[0]; floor++ {
			if (ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0 {
				return DOWN;
 			}
		}
	} else if (directions[0] == UP){
		for floor:=previousFloors[0]+1; floor<4; floor++ {
			if (ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0 {
				return UP;
			}
		}
		for floor:=0; floor<previousFloors[0]; floor++ {
			if (ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0 {
				return DOWN;
			}
		}
	} else if (directions[0] == DOWN){
		for floor:=0; floor<previousFloors[0]; floor++ {
			if (ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0 {
				return DOWN;
			}
		}
		for floor:=previousFloors[0]+1; floor<4; floor++ {
			if (ordersUp[0][floor] | ordersDown[0][floor] | ordersOut[0][floor]) != 0 {
				return UP;
			}
		}
	}
	return STOP;
}



















