package timer

import (
	//."fmt"
	"time"
)

const (
	doorOpenTime = 3*time.Second
)

var (
	doorTimer time.Time
	doorTimerFlag bool = false
)

func DoorTimer(doorTimerChan chan int, doorTimeoutChan chan int) {
	go checkTimer(doorTimeoutChan)
}

func checkDoorTimer(doorTimeoutChan chan int) {
	for {
		/*if (int64(time.Since(timer0)) > int64(timeUntilTimeout)) && (timer0flag == true) {
			timer0flag = false
			Println("								doortimer timeout")
		}*/
	}
}
