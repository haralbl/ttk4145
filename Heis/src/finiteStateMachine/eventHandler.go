package finiteStateMachine

import (
	."fmt"
	"driver"
	"defines"
)

var (
	ElevatorStatus defines.ElevatorStatus_t
)

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

func handleMessage(sendChan chan []byte, doorTimerChan chan string, elevatorIP string, floor int, buttonType int, MessageType string, enableStuckTimerChan chan int) {
	switch (MessageType) {
	case "":
		Println("json is eating the message, and shitting out an empty status")
	case "ack":
		Println("received ack")
		
		var elevator int = -1
		for i:=0; i<defines.NumberOfElevators; i++ {
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
		for floor:=0; floor<defines.NumberOfFloors; floor++ {
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
	
func event_newOrder(sendChan chan []byte, doorTimerChan chan string, enableStuckTimerChan chan int) {
	switch (ElevatorStatus.State) {
	case defines.IDLE:
		enableStuckTimerChan <- 1
		if nextDirection() == defines.STOP {
			driver.Set_door_open_lamp(1)
			doorTimerChan <- "start"
			ElevatorStatus.State = defines.DOOR_OPEN
			Print("State = defines.DOOR_OPEN\n")
			sendChan <- wrapMessage("orderCompleted", 0, "", ElevatorStatus.PreviousFloors[0])
		} else if nextDirection() == defines.UP {
			driver.Set_motor_direction(defines.UP)
			ElevatorStatus.Directions[0] = defines.UP
 			ElevatorStatus.State = defines.MOVING 			
 			Print("State = defines.MOVING\n")
		} else if nextDirection() == defines.DOWN {
			driver.Set_motor_direction(defines.DOWN)
			ElevatorStatus.Directions[0] = defines.DOWN
			ElevatorStatus.State = defines.MOVING
			Print("State = defines.MOVING\n")
		} else {
			Printf("ERROR, event_newOrder: nextDirection returns invalid value")
		}
	case defines.DOOR_OPEN:
		// do nothing
	case defines.MOVING:
		// do nothing
	}
}

func event_floorReached(StateOfShouldStop int, doorTimerChan chan string) {
	driver.Set_floor_indicator(ElevatorStatus.PreviousFloors[0])
	switch ElevatorStatus.State {
	case defines.IDLE:
		Printf("ERROR, event: floorReached when not moving!")
	case defines.DOOR_OPEN:
		Printf("ERROR, event: floorReached when not moving!")
	case defines.MOVING:
		if StateOfShouldStop == 1 {
			driver.Set_motor_direction(defines.STOP)
			driver.Set_door_open_lamp(1)
			doorTimerChan <- "start"
			ElevatorStatus.State = defines.DOOR_OPEN
			Print("State = defines.DOOR_OPEN\n")
		}
	}
}

func event_doorTimerOut(enableStuckTimerChan chan int) {
	switch ElevatorStatus.State {
	case defines.IDLE:
		// do nothing
	case defines.DOOR_OPEN:
		driver.Set_door_open_lamp(0)
		if nextDirection() == defines.STOP {  
			ElevatorStatus.State = defines.IDLE
			enableStuckTimerChan <- 0
			Print("State = defines.IDLE\n")
		} else if nextDirection() == defines.UP {
			driver.Set_motor_direction(defines.UP)
			ElevatorStatus.Directions[0] = defines.UP
 			ElevatorStatus.State = defines.MOVING
 			Print("State = defines.MOVING\n")
		} else if nextDirection() == defines.DOWN {
			driver.Set_motor_direction(defines.DOWN)
			ElevatorStatus.Directions[0] = defines.DOWN
			ElevatorStatus.State = defines.MOVING
			Print("State = defines.MOVING\n")
		} else {
			Printf("ERROR, event_doorTimerOut: nextDirection returns invalid value")
		}
	case defines.MOVING:
		// do nothing
	}
}

func shouldStop() int {
	if ElevatorStatus.Directions[0] == defines.UP {
		if (ElevatorStatus.OrdersUp[0][ElevatorStatus.PreviousFloors[0]] | ElevatorStatus.OrdersOut[0][ElevatorStatus.PreviousFloors[0]]) != 0 {
			return 1
		}
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<defines.NumberOfFloors; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return 0
			}
		}
	}
	if ElevatorStatus.Directions[0] == defines.DOWN {
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
				return defines.UP
			}
		}
	} else if ElevatorStatus.PreviousFloors[0] == defines.NumberOfFloors {
	
		Println("YOLO etasje 3")
	
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.DOWN
 			}
		}
	} else if (ElevatorStatus.Directions[0] == defines.UP){
	
		Println("YOLO retning opp")
	
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.UP
			}
		}
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.DOWN;
			}
		}
	} else if (ElevatorStatus.Directions[0] == defines.DOWN){
	
		Println("YOLO retning ned")
	
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.DOWN
			}
		}
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<4; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.UP
			}
		}
	}
	return defines.STOP
}











