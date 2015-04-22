package timer

import (
	."fmt"
	"time"
	"os"
)

const (
	timeUntilStuck = 10*time.Second
)

var (
	enable int = 0
	stuckTimer time.Time
	stuckTimerFlag bool = false
)

func StuckTimer(resetStuckTimerChan chan string, enableStuckTimerChan chan int, stuckTimeoutChan chan int) {
	go checkStuckTimer(stuckTimeoutChan)
	
	for {
		select {
		case <- resetStuckTimerChan:
			stuckTimer = time.Now()
			
		case enable = <- enableStuckTimerChan:
			if enable == 1 {
				stuckTimer = time.Now()
				stuckTimerFlag = true
			} else {
				stuckTimerFlag = false
			}
			
		case <- stuckTimeoutChan:
			os.Exit(1)
		}
	}
}

func checkStuckTimer(stuckTimeoutChan chan int) {
	for {
		if (int64(time.Since(stuckTimer)) > int64(timeUntilStuck)) && (stuckTimerFlag == true) {
			//stuckTimerFlag = false
			Println("stuckTimer timeout")
			stuckTimeoutChan <- 1
		}
		time.Sleep(100*time.Millisecond)
	}
}





