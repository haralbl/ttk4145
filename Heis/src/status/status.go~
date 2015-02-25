package status

import (
	//."fmt"
	"strings"
	"strconv"
	"network"
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
	
	ordersUp			[numberOfElevators][numberOfFloors]int
	ordersDown			[numberOfElevators][numberOfFloors]int
	ordersOut			[numberOfElevators][numberOfFloors]int
)

func Get() (status string) {
	status = ""
	for i:=0; i<numberOfElevators; i++ {
		status += activeElevators[i]
		status += " "
	}
	status += "\n"
	for i:=0; i<numberOfElevators; i++ {
		status += strconv.Itoa(lastPositions[i])
		status += " "
	}
	status += "\n"
	for i:=0; i<numberOfElevators; i++ {
		status += strconv.Itoa(inFloor[i])
		status += " "
	}
	status += "\n"
	for i:=0; i<numberOfElevators; i++ {
		for j:=0; j<numberOfFloors; j++ {
			status += strconv.Itoa(ordersUp[i][j])
			status += " "
			status += strconv.Itoa(ordersDown[i][j])
			status += " "
			status += strconv.Itoa(ordersOut[i][j])
			status += " "
		}
		status += "\n"
	}
	
	// Mask: used for new orders
	for i:=0; i<numberOfElevators; i++ {
		status += 0
		status += " "
	}
	status += "\n"
	for i:=0; i<numberOfElevators; i++ {
		status += 0
		status += " "
	}
	status += "\n"
	for i:=0; i<numberOfElevators; i++ {
		status += 0
		status += " "
	}
	status += "\n"
	for i:=0; i<numberOfElevators; i++ {
		for j:=0; j<numberOfFloors; j++ {
			status += 0
			status += " "
			status += 0
			status += " "
			status += 0
			status += " "
		}
		status += "\n"
	}
	return
}

func Update(status string) {
	statusFields := strings.Split(status, "\n")
	
	field = strings.Split(statusFields[1], " ")
	for i:=0; i<numberOfElevators; i++ {
		lastPositions[i],_	= strconv.Atoi(field[i])
	}
	field = strings.Split(statusFields[2], " ")
	for i:=0; i<numberOfElevators; i++ {
		inFloor[i],_		= strconv.Atoi(field[i])
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
		lastPositions[0]	= 0
		inFloor[0] 			= 0
		
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














