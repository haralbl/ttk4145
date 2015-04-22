package timer

import (
	."fmt"
	"time"
	"defines"
)

const (
	timeUntilTimeout = 700*time.Millisecond
)

var (
	timers		[defines.NumberOfElevators]time.Time
	timerFlags	[defines.NumberOfElevators]bool
)

func Init() {
	for i:=0; i<defines.NumberOfElevators; i++ {
		timerFlags[i] = false
	}
}

func ElevatorTimer(elevatorTimerChan chan int, elevatorTimeoutChan chan int) {
	elevatorN := -1
	go checkTimer(elevatorTimeoutChan)
	
	for {
		select {
		case elevatorN = <- elevatorTimerChan:
			timers[elevatorN] = time.Now()
			timerFlags[elevatorN] = true
			
		case elevatorN = <- elevatorTimeoutChan:
			elevatorTimerChan <- elevatorN
		}
	}
}

func checkTimer(elevatorTimeoutChan chan int) {
	for {
		for i:=0; i<defines.NumberOfElevators; i++ {
			if (int64(time.Since(timers[i])) > int64(timeUntilTimeout)) && (timerFlags[i] == true) {
				timerFlags[i] = false
				elevatorTimeoutChan <- i
				Printf("timer %d timeout", i)
			}
		}
		time.Sleep(10*time.Millisecond)
	}
}





/*package timer

import (
	."fmt"
	"time"
	"status"
)

const (
	timeUntilTimeout = 700*time.Millisecond
)

var (
	timers		[defines.NumberOfElevators]time.Time

	timerFlags	[defines.NumberOfElevators]bool
)

func Init() {
	for i:=0; i<defines.NumberOfElevators; i++ {
		timerFlags[i] = false
	}
}

func ElevatorTimer(elevatorTimerChan chan int, elevatorTimeoutChan chan int) {
	elevatorN := -1
	go checkTimer(elevatorTimeoutChan)
	
	for {
		select {
		case elevatorN = <- elevatorTimerChan:
			switch elevatorN {
			case 0:
				timers[0] = time.Now()
				timerFlags[0] = true
			case 1:
				timers[1] = time.Now()
				timerFlags[1] = true
			case 2:
				timers[2] = time.Now()
				timerFlags[2] = true
			}
		case elevatorN = <- elevatorTimeoutChan:
			elevatorTimerChan <- elevatorN
		}
	}
}

func checkTimer(elevatorTimeoutChan chan int) {
	for {
		if (int64(time.Since(timers[0])) > int64(timeUntilTimeout)) && (timerFlags[0] == true) {
			timerFlags[0] = false
			elevatorTimeoutChan <- 0
			Println("timer 0 timeout")
		}
		if (int64(time.Since(timers[1])) > int64(timeUntilTimeout)) && (timerFlags[1] == true) {
			timerFlags[1] = false
			elevatorTimeoutChan <- 1
			Println("timer 1 timeout")
		}
		if (int64(time.Since(timers[2])) > int64(timeUntilTimeout)) && (timerFlags[2] == true) {
			timerFlags[2] = false
			elevatorTimeoutChan <- 2
			Println("timer 2 timeout")
		}
		time.Sleep(10*time.Millisecond)
	}
}*/











