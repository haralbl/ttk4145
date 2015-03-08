package status

import (
	."fmt"
	"strings"
	"strconv"
	"network"
	"math"
	"driver"
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
		newOrderUp[elevator][floor] = 1
	case 2:
		newOrderUp[elevator][floor] = 1
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
	for i:=0; i<numberOfElevators; i++ {
		temp,_	= strconv.Atoi(field[i])
		previousFloors[i] = previousFloors[i] | temp
	}
	field = strings.Split(statusFields[6], " ")
	for i:=0; i<numberOfElevators; i++ {
		temp,_		= strconv.Atoi(field[i])
		inFloor[i] = inFloor[i] | temp
	}
	field = strings.Split(statusFields[7], " ")
	for i:=0; i<numberOfElevators; i++ {
		temp,_		= strconv.Atoi(field[i])
		directions[i] = directions[i] | temp
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

func handleMessage(sendChan chan string, resetAckChan chan string, elevator int, floor int, buttonType int, order int, messageType string){
	switch (messageType) {
	case "ack":
		Println("received ack")
		resetAckChan <- "reset"
		
		// legge til ordre du ikke selv skal ta
		
	case "newOrder":
		if elevator == 0{
			Println("received new order to handle myself")
			sendChan <- wrapMessage("ack", buttonType, elevator, floor)
		
			// oppdater status	
			switch (buttonType) {
			case 0:
				if ordersUp[elevator][floor] == 0 {
					sendChan <- wrapMessage("ack", buttonType, elevator, floor)
					ordersUp[elevator][floor] = 1
					event_newOrder()
				}
			case 1:
				if ordersDown[elevator][floor] == 0 {
					sendChan <- wrapMessage("ack", buttonType, elevator, floor)
					ordersDown[elevator][floor] = 1
					event_newOrder()
				}
			case 2:
				if ordersDown[elevator][floor] == 0 {
					sendChan <- wrapMessage("ack", buttonType, elevator, floor)
					ordersDown[elevator][floor] = 1
					event_newOrder()
				}
			}
		}
		
	case "floorReached":
		Println("received floorReached")
		
	case "orderCompleted":
		Println("received orderCompleted")
		// oppdater status
		switch (buttonType) {
		case 0:
			ordersUp[elevator][floor] = 0
		case 1:
			ordersDown[elevator][floor] = 0
		case 2:
			ordersOut[elevator][floor] = 0
		}
	}
}

func Initialize() {
	for i:=0; i<numberOfElevators; i++ {
		activeElevators[i] = "empty"
	} 
	activeElevators[0] = network.GetLocalIP()
	
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

func CheckAliveElevators(receiveAliveMessageChan chan string, elevatorTimerChan chan int) {
	var elevatorIP	string
	var elevatorN	int
	for {
		select {
		case elevatorIP = <- receiveAliveMessageChan:
			elevatorIP = elevatorIP[0:15]
			elevatorN = isElevatorInList(elevatorIP)
			if elevatorN == -1 {
				elevatorN = addElevator(elevatorIP)
			}
			elevatorTimerChan <- elevatorN
		case elevatorN = <- elevatorTimerChan:
			removeElevator(elevatorN)
		}
	}
}

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
	
	if buttonType == 2 {
		costs[0] = 0
	}
	
	for i:=0; i<numberOfElevators; i++ {
		if costs[i] < int(cheapestCost) {
			cheapestElevator = i
		}
	}
	return
}

func EventHandler(sendChan chan string, upButtonChan chan int, downButtonChan chan int,
					commandButtonChan chan int, floorChan chan int, ackTimerChan chan string,
					receiveChan chan string, resetAckChan chan string, ackCheckChan chan string,
					ackTimeoutChan chan string, doorTimerChan chan int) {
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
			chosenElevator = costFunction(button, 2)
			ackTimerChan <- "acktimer"
			sendChan <- wrapMessage("newOrder", 2, chosenElevator, button)
		case ack = <- ackTimeoutChan:
			ack = ack ////////////////// fjærn!!!!
			// ta ordre selv
		case message = <- receiveChan:
			elevator, floor, buttonType, order, messageType := unwrapMessage(message)
			handleMessage(sendChan, resetAckChan, elevator, floor, buttonType, order, messageType)
		case previousFloors[0] = <- floorChan:
			sendChan <- wrapMessage("floorReached", 0, 0, previousFloors[0])
			event_floorReached()
		//case ordre fullføres
			//nextDirection()
		//case ack mottatt
		case <- doorTimerChan:
			event_doorTimerOut()
		}
	}
}
	
func event_newOrder() {
	switch (state) {
	case IDLE:
		//set button lamp
		if nextDirection() == STOP {  
			state = IDLE;
		} else if nextDirection() == UP {
			driver.Set_motor_direction(UP);
			directions[0] = UP;
 			state = MOVING;
		} else if nextDirection() == DOWN {
			driver.Set_motor_direction(DOWN);
			directions[0] = DOWN;
			state = MOVING;
		} else {
			Printf("ERROR, event_newOrder: nextDirection returns invalid value")
		}
	case DOOR_OPEN:
		//set button lamp
	case MOVING:
		//set button lamp
	}
}

func event_floorReached() {
	// set floor lights
	switch state {
	case IDLE:
		Printf("ERROR, event: floorReached when not moving!")
	case DOOR_OPEN:
		Printf("ERROR, event: floorReached when not moving!")
	case MOVING:
		if shouldStop() == 1 {
			driver.Set_motor_direction(STOP)
			driver.Set_door_open_lamp(1)
			//starte door timer
			state = DOOR_OPEN
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
		} else if nextDirection() == UP {
			driver.Set_motor_direction(UP);
			directions[0] = UP;
 			state = MOVING;
		} else if nextDirection() == DOWN {
			driver.Set_motor_direction(DOWN);
			directions[0] = DOWN;
			state = MOVING;
		} else {
			Printf("ERROR, event_doorTimerOut: nextDirection returns invalid value")
		}
	case MOVING:
		// close door? (shouldnt happen)
	}
}

func shouldStop() int{
	if directions[0] == UP {
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



















