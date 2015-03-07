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

func AckTimer(ackTimerChan chan string, ackTimeoutChan chan string, resetAckTimer chan string, 			
				ackCheckChan chan string) {
	go checkAckTimer(ackCheckChan)
	for {
		select {
		case <- ackTimerChan:
			ackTimer = time.Now()	
			ackTimerFlag = true
		case <- ackCheckChan:
			ackTimeoutChan <- "timeout"
		case <- resetAckTimer:
			ackTimerFlag = false
		}
	}
}

func checkAckTimer(ackCheckChan chan string) {
	for {
		if (int64(time.Since(ackTimer)) > int64(timeUntilAckTimeout)) && (ackTimerFlag == true) {
			ackTimerFlag = false
			ackCheckChan <- "timeout"
			Println("								ackTimer timeout")
		}
	}
}

