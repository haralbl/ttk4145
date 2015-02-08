package status

import (
	."fmt"
	"net"
	"strings"
	"strconv"
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

func GetStatus() (status string) {
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
	return
}

func UpdateStatus(status string) {
	statusFields := strings.Split(status, "\n")
	
	field := strings.Split(statusFields[0], " ")
	for i:=0; i<numberOfElevators; i++ {
		activeElevators[i]	= field[i]
	}
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




