package finiteStateMachine

import (
	"defines"
	"math"
)

func costFunction(floor int, buttonType int) (cheapestElevator int) {
	var costs[defines.NumberOfElevators]int
	for i:=0; i<defines.NumberOfElevators; i++ {
		costs[i] = 0
	}
	
	for i:=0; i<defines.NumberOfElevators; i++ {
		
		// Check number of orders
		for j:=0; j<defines.NumberOfFloors; j++ {
			if ElevatorStatus.OrdersUp[i][j] == 1 || ElevatorStatus.OrdersDown[i][j] == 1 || ElevatorStatus.OrdersOut[i][j] == 1 {
				costs[i] += 2
			}
		}
		// Check distances
		costs[i] += int(math.Abs(float64(floor) - float64(ElevatorStatus.PreviousFloors[i])))
		
		// Check if order in same direction in front of elevator
		if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.UP && buttonType == defines.BUTTON_CALL_UP {
			costs[i] += 2*defines.NumberOfFloors
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.DOWN && buttonType == defines.BUTTON_CALL_DOWN {
			costs[i] += 2*defines.NumberOfFloors
			
		// Check if order in opposite direction in front of elevator	
		} else if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.UP && buttonType == defines.BUTTON_CALL_DOWN {
			costs[i] += 5*defines.NumberOfFloors
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.DOWN && buttonType == defines.BUTTON_CALL_UP {
			costs[i] += 5*defines.NumberOfFloors

		// Check if order in opposite direction behind elevator
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.UP && buttonType == defines.BUTTON_CALL_DOWN {
			costs[i] += 8*defines.NumberOfFloors
		} else if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.DOWN && buttonType == defines.BUTTON_CALL_UP {
			costs[i] += 8*defines.NumberOfFloors

		// Check if order in same direction behind elevator
		} else if floor < ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.UP && buttonType == defines.BUTTON_CALL_UP {
			costs[i] += 11*defines.NumberOfFloors
		} else if floor > ElevatorStatus.PreviousFloors[i] && ElevatorStatus.Directions[i] == defines.DOWN && buttonType == defines.BUTTON_CALL_DOWN {
			costs[i] += 11*defines.NumberOfFloors
		}
		if buttonType == defines.BUTTON_CALL_UP && ElevatorStatus.OrdersUp[i][floor] == 1 {
			costs[i] = 0
		}
		if buttonType == defines.BUTTON_CALL_DOWN && ElevatorStatus.OrdersDown[i][floor] == 1 {
			costs[i] = 0
		}
	}
	
	cheapestElevator	= 0
	cheapestCost		:= 10000
	
	for i:=0; i<defines.NumberOfElevators; i++ {
		if costs[i] < int(cheapestCost) && ElevatorStatus.ActiveElevators[i] != "empty" {
			cheapestCost = costs[i]
			cheapestElevator = i
		}
	}
	return
}














