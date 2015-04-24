package finiteStateMachine

import (
	."fmt"
	"driver"
	"defines"
	"os"
	"os/exec"
)

var (
	ElevatorStatus defines.ElevatorStatus_t
)

func EventHandler(sendChan chan []byte, upButtonChan chan int, downButtonChan chan int,
					commandButtonChan chan int, floorReachedChan chan int, floorLeftChan chan string, receiveChan chan []byte,
					doorTimerChan chan string, resetStuckTimerChan chan string, enableStuckTimerChan chan int) {
	var button			int
	var chosenElevator	string
	var message			[]byte
	for {
		select {
		case button = <- upButtonChan:
			chosenElevator = ElevatorStatus.ActiveElevators[costFunction(button, defines.BUTTON_CALL_UP)]
			sendChan <- wrapMessage("newOrder", defines.BUTTON_CALL_UP, chosenElevator, button)
			
		case button = <- downButtonChan:
			chosenElevator = ElevatorStatus.ActiveElevators[costFunction(button, defines.BUTTON_CALL_DOWN)]
			sendChan <- wrapMessage("newOrder", defines.BUTTON_CALL_DOWN, chosenElevator, button)
			
		case button = <- commandButtonChan:
			chosenElevator = ElevatorStatus.ActiveElevators[0]
			sendChan <- wrapMessage("newOrder", defines.BUTTON_COMMAND, chosenElevator, button)

		case message = <- receiveChan:
			elevator, floor, buttonType, MessageType := unwrapMessage(message)
			messageHandler(sendChan, doorTimerChan, elevator, floor, buttonType, MessageType, enableStuckTimerChan)
			
		case ElevatorStatus.PreviousFloors[0] = <- floorReachedChan:
			event_floorReached(sendChan, resetStuckTimerChan, doorTimerChan)
			
		case <- floorLeftChan:
			event_floorLeft(sendChan)
			
		case <- doorTimerChan:
			event_doorTimerOut(sendChan, enableStuckTimerChan)
		}
	}
}

func messageHandler(sendChan chan []byte, doorTimerChan chan string, elevatorIP string, floor int, buttonType int, MessageType string, enableStuckTimerChan chan int) {
	switch (MessageType) {
	case "":
		Println("received empty message")
		
	case "floorReached":
		Println("received floorReached")
		
	case "ack":
		event_ackReceived(elevatorIP, floor, buttonType)
		
	case "newOrder":
		event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan, elevatorIP, floor, buttonType)
		
	case "orderCompleted":
		event_orderCompleted(elevatorIP, floor)
		
	case "updateAwokenElevator":
		event_wokeUp(sendChan, doorTimerChan, enableStuckTimerChan, elevatorIP, floor, buttonType)
	}
}
	
func event_newOrder(sendChan chan []byte, doorTimerChan chan string, enableStuckTimerChan chan int, elevatorIP string, floor int, buttonType int) {
	enableStuckTimerChan <- 1
	
	if elevatorIP == ElevatorStatus.ActiveElevators[0] {
		Println("received new order to handle myself")
		
		driver.Set_button_lamp(buttonType, floor, 1)
		sendChan <- wrapMessage("ack", buttonType, elevatorIP, floor)
	
		switch (buttonType) {
		case defines.BUTTON_CALL_UP:
			if ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
				ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor] = 1
			}
		case defines.BUTTON_CALL_DOWN:
			if ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
				ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] = 1
			}
		case defines.BUTTON_COMMAND:
			if ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor] == 0 {
				ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor] = 1
			}
		}
	}
	
	switch (ElevatorStatus.State) {
	case defines.IDLE:
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

func event_ackReceived(elevatorIP string, floor int, buttonType int) {
	Println("received ack")
		
	var elevator int = -1
	for i:=0; i<defines.NumberOfElevators; i++ {
		if elevatorIP == ElevatorStatus.ActiveElevators[i] {
			elevator = i
		}
	}
	if elevator == -1 {
		Println("Received ack from unknown IP. I will take the order myself.")
		elevator = 0
	}
	
	switch(buttonType) {
	case defines.BUTTON_CALL_UP:
		ElevatorStatus.OrdersUp[elevator][floor] = 1
		driver.Set_button_lamp(buttonType, floor, 1)
	case defines.BUTTON_CALL_DOWN:
		ElevatorStatus.OrdersDown[elevator][floor] = 1
		driver.Set_button_lamp(buttonType, floor, 1)
	case defines.BUTTON_COMMAND:
		ElevatorStatus.OrdersOut[elevator][floor] = 1
	}
}

func event_floorReached(sendChan chan []byte, resetStuckTimerChan chan string, doorTimerChan chan string) {
	sendChan <- wrapMessage("floorReached", 0, "", ElevatorStatus.PreviousFloors[0])
	driver.Set_floor_indicator(ElevatorStatus.PreviousFloors[0])
	resetStuckTimerChan <- "reset"
	
	switch ElevatorStatus.State {
	case defines.IDLE:
		newElevator := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
		_ = newElevator.Run()
		os.Exit(1)
	case defines.DOOR_OPEN:
		newElevator := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
		_ = newElevator.Run()
		os.Exit(1)
	case defines.MOVING:
		if shouldStop() == 1 {
			driver.Set_motor_direction(defines.STOP)
			driver.Set_door_open_lamp(1)
			doorTimerChan <- "start"
			ElevatorStatus.State = defines.DOOR_OPEN
			Print("State = defines.DOOR_OPEN\n")
		}
	}
}

func event_floorLeft(sendChan chan []byte) {
	switch ElevatorStatus.State {
	case defines.IDLE:
		newElevator := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
		_ = newElevator.Run()
		os.Exit(1)
	case defines.DOOR_OPEN:
		newElevator := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
		_ = newElevator.Run()
		os.Exit(1)
	case defines.MOVING:
		// do nothing, all is well
	}
}

func event_doorTimerOut(sendChan chan []byte, enableStuckTimerChan chan int) {
	sendChan <- wrapMessage("orderCompleted", 0, ElevatorStatus.ActiveElevators[0], ElevatorStatus.PreviousFloors[0])
	
	switch ElevatorStatus.State {
	case defines.IDLE:
		driver.Set_door_open_lamp(0)
		Printf("ERROR, event_doorTimerOut: door was open in IDLE")
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
		newElevator := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
		_ = newElevator.Run()
		os.Exit(1)
	}
}

func event_orderCompleted(elevatorIP string, floor int) {
	Println("received orderCompleted")
	
	ElevatorStatus.OrdersUp[elevatorIPtoIndex(elevatorIP)][floor]	= 0
	ElevatorStatus.OrdersDown[elevatorIPtoIndex(elevatorIP)][floor] = 0
	ElevatorStatus.OrdersOut[elevatorIPtoIndex(elevatorIP)][floor]	= 0
	driver.Set_button_lamp(0, floor, 0)
	driver.Set_button_lamp(1, floor, 0)
	if ElevatorStatus.ActiveElevators[0] == elevatorIP {
		driver.Set_button_lamp(2, floor, 0)
	}
}

func event_wokeUp(sendChan chan []byte, doorTimerChan chan string, enableStuckTimerChan chan int, elevatorIP string, floor int, buttonType int) {
	Println("received updateAwokenElevator")
	
	for floor:=0; floor<defines.NumberOfFloors; floor++ {
		if ElevatorStatus.OrdersUp[0][floor] == 1 {
			driver.Set_button_lamp(defines.BUTTON_CALL_UP, floor, 1)
			event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan, elevatorIP, floor, buttonType)
		}
		if ElevatorStatus.OrdersDown[0][floor] == 1 {
			driver.Set_button_lamp(defines.BUTTON_CALL_DOWN, floor, 1)
			event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan, elevatorIP, floor, buttonType)
		}
		if ElevatorStatus.OrdersOut[0][floor] == 1 {
			driver.Set_button_lamp(defines.BUTTON_COMMAND, floor, 1)
			event_newOrder(sendChan, doorTimerChan, enableStuckTimerChan, elevatorIP, floor, buttonType)
		}
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
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<defines.NumberOfFloors; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.UP
			}
		}
	} else if ElevatorStatus.PreviousFloors[0] == defines.NumberOfFloors {
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.DOWN
 			}
		}
	} else if (ElevatorStatus.Directions[0] == defines.UP){
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<defines.NumberOfFloors; floor++ {
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
		for floor:=0; floor<ElevatorStatus.PreviousFloors[0]; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.DOWN
			}
		}
		for floor:=ElevatorStatus.PreviousFloors[0]+1; floor<defines.NumberOfFloors; floor++ {
			if (ElevatorStatus.OrdersUp[0][floor] | ElevatorStatus.OrdersDown[0][floor] | ElevatorStatus.OrdersOut[0][floor]) != 0 {
				return defines.UP
			}
		}
	}
	return defines.STOP
}











