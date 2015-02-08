package timer

import (
	"time"
)

var (
	timer0 time.Time
	timer1 time.Time
	timer2 time.Time
)

func elevatorTimer(elevatorTimerChan chan int) {
	elevatorN := -1
	//when elevatorN <- elevatorTimerChan
		switch elevatorN {
		case 0:
			timer0 = time.Now()
		case 1:
			timer1 = time.Now()
		case 2:
			timer2 = time.Now()
		}
	//when time.Now() - timerN > someTime
		elevatorTimerChan <- elevatorN
}

