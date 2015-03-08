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

func doorTimer(doorTimerChan chan string, doorTimeoutChan chan int) {
	go checkTimer(doorTimeoutChan)
	for {
		select {
		case <- doorTimerChan:
			doorTimer = time.Now()
			doorTimerFlag = true
		case <- doorTimeoutChan:
			doorTimerChan <- "doorTimeout"
		}
	}
}

func checkTimer(doorTimeoutChan chan int) {
	for {
		if (int64(time.Since(timer0)) > int64(timeUntilTimeout)) && (timer0flag == true) {
			timer0flag = false
			elevatorTimeoutChan <- 0
			Println("								timer 0 timeout")
		}
	}
}






