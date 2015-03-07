package status

import (
	//."fmt"
	"strings"
	"strconv"
	"network"
	"math"
	//"time"
	//"./Timer"
)

const (
	numberOfElevators	= 3
	numberOfFloors		= 4
)

var (
	activeElevators		[numberOfElevators]string // IP addresses
	
	lastPositions		[numberOfElevators]int
	inFloor				[numberOfElevators]int
	directions			[numberOfElevators]int
	
	ordersUp			[numberOfElevators][numberOfFloors]int
	ordersDown			[numberOfElevators][numberOfFloors]int
	ordersOut			[numberOfElevators][numberOfFloors]int
)

func createMessage(messageType int, buttonType int, elevator int, floor int) (message string) {
	message = ""
	message += strconv.Itoa(messageType) + "\n"
	
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
	message += "\n"
	
	for i:=0; i<numberOfElevators; i++ {
		message += activeElevators[i]
		message += " "
	}
	message += "\n"
	for i:=0; i<numberOfElevators; i++ {
		message += strconv.Itoa(lastPositions[i])
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

func Update(status string) {
	statusFields := strings.Split(status, "\n")
	
	field := strings.Split(statusFields[1], " ")
	for i:=0; i<numberOfElevators; i++ {
		lastPositions[i],_	= strconv.Atoi(field[i])
	}
	field = strings.Split(statusFields[2], " ")
	for i:=0; i<numberOfElevators; i++ {
		inFloor[i],_		= strconv.Atoi(field[i])
	}
	field = strings.Split(statusFields[2], " ")
	for i:=0; i<numberOfElevators; i++ {
		directions[i],_		= strconv.Atoi(field[i])
	}
	for i:=0; i<numberOfElevators; i++ {
		field = strings.Split(statusFields[3+i], " ")
		for j:=0; j<numberOfFloors; j++ {
			ordersUp[i][j],_	= strconv.Atoi(field[3*j+0])
			ordersDown[i][j],_	= strconv.Atoi(field[3*j+1])
			ordersOut[i][j],_	= strconv.Atoi(field[3*j+2])
		}
	}
}

func Initialize() {
	for i:=0; i<numberOfElevators; i++ {
		activeElevators[i] = "empty"
	} 
	activeElevators[0] = network.GetLocalIP()
	
	for i:=0; i<numberOfElevators; i++ {
		lastPositions[i]	= 0
		inFloor[i] 			= 0
		directions[i]		= 0
		
		for j:=0; j<numberOfFloors; j++ {
			ordersUp[i][j]		= 0
			ordersDown[i][j]	= 0
			ordersOut[i][j]		= 0
		}
	}
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
		if (floor > lastPositions[i] && directions[i] == 1) || (floor < lastPositions[i] && directions[i] == 0) {
			costs[i] += 5*numberOfFloors
		}
		// Check distances
		costs[i] += int(2*math.Abs(float64(floor) - float64(lastPositions[i])))
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
}

func EventHandler(sendChan chan string, upButtonChan chan int, downButtonChan chan int, commandButtonChan chan int, floorChan chan int) {
	var button			int
	//var currentFloor	int
	var chosenElevator	int
	for {
		select {
		case button = <- upButtonChan:
			chosenElevator = costFunction(button, 0)
			//spawne timer
			sendChan <- createMessage(0, 0, chosenElevator, button)
		case button = <- downButtonChan:
			chosenElevator = costFunction(button, 1)
			//spawne timer
			sendChan <- createMessage(0, 1, chosenElevator, button)
		case button = <- commandButtonChan:
			chosenElevator = costFunction(button, 2)
			//spawne timer
			sendChan <- createMessage(0, 2, chosenElevator, button)
		//case currentFloor = <- floorChan:
		//case //ordre fullføres
		
		}
	}
}









