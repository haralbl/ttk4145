package timer

import (
	."fmt"
	"time"
	"os"
	"os/exec"
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
		}
	}
}

func checkStuckTimer(stuckTimeoutChan chan int) {
	for {
		if (int64(time.Since(stuckTimer)) > int64(timeUntilStuck)) && (stuckTimerFlag == true) {
			Println("stuckTimer timeout")
			newElevator := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			_ = newElevator.Run()
			os.Exit(1)
		}
		time.Sleep(100*time.Millisecond)
	}
}





