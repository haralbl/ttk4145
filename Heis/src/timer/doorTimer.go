package timer

import (
	."fmt"
	"time"
)

const (
	doorOpenTime = 3*time.Second
)

var (
	doorTimer time.Time
	doorTimerFlag bool = false
)

func DoorTimer(doorTimerChan chan string, doorTimeoutChan chan int) {
	go checkDoorTimer(doorTimeoutChan)
	
	for {
		select {
		case <- doorTimerChan:
			doorTimer = time.Now()
			Println("starting timer\n")
			doorTimerFlag = true
		case <- doorTimeoutChan:
			doorTimerChan <- "doorTimeout"
		}
	}
}

func checkDoorTimer(doorTimeoutChan chan int) {
	for {
		if (int64(time.Since(doorTimer)) > int64(doorOpenTime)) && (doorTimerFlag == true) {
			doorTimerFlag = false
			doorTimeoutChan <- 0
			Println("doorTimer timeout")
		}
		time.Sleep(100*time.Millisecond)
	}
}











