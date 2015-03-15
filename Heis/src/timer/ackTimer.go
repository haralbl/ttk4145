package timer

import (
	."fmt"
	"time"
)

const (
	timeUntilAckTimeout = 300*time.Millisecond
)
var (
	ackTimer		time.Time
	ackTimerFlag	bool = false
)

func AckTimer(ackTimerChan chan string, ackTimeoutChan chan string, ackResetChan chan string) {
	go checkAckTimer(ackTimeoutChan)
	for {
		select {
		case <- ackTimerChan:
			ackTimer = time.Now()	
			ackTimerFlag = true
		case <- ackTimeoutChan:
			ackTimerChan <- "ackTimeout"
		case <- ackResetChan:
			ackTimerFlag = false
		}
	}
}

func checkAckTimer(ackTimeoutChan chan string) {
	for {
		if (int64(time.Since(ackTimer)) > int64(timeUntilAckTimeout)) && (ackTimerFlag == true) {
			ackTimerFlag = false
			ackTimeoutChan <- "timeout"
			Println("					ackTimer timeout")
		}
		time.Sleep(10*time.Millisecond)
	}
}

