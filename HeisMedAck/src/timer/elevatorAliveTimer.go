package timer

import (
	."fmt"
	"time"
)

const (
	timeUntilTimeout = 700*time.Millisecond
)
var (
	timer0		time.Time
	timer1		time.Time
	timer2		time.Time
	timer0flag	bool = false
	timer1flag	bool = false
	timer2flag	bool = false
)

func ElevatorTimer(elevatorTimerChan chan int, elevatorTimeoutChan chan int) {
	elevatorN := -1
	go checkTimer(elevatorTimeoutChan)
	for {
		select {
		case elevatorN = <- elevatorTimerChan:
			switch elevatorN {
			case 0:
				timer0 = time.Now()
				timer0flag = true
			case 1:
				timer1 = time.Now()
				timer1flag = true
			case 2:
				timer2 = time.Now()
				timer2flag = true
			}
		case elevatorN = <- elevatorTimeoutChan:
			elevatorTimerChan <- elevatorN
		}
	}
}

func checkTimer(elevatorTimeoutChan chan int) {
	for {
		if (int64(time.Since(timer0)) > int64(timeUntilTimeout)) && (timer0flag == true) {
			timer0flag = false
			elevatorTimeoutChan <- 0
			Println("								timer 0 timeout")
		}
		if (int64(time.Since(timer1)) > int64(timeUntilTimeout)) && (timer1flag == true) {
			timer1flag = false
			elevatorTimeoutChan <- 1
			Println("								timer 1 timeout")
		}
		if (int64(time.Since(timer2)) > int64(timeUntilTimeout)) && (timer2flag == true) {
			timer2flag = false
			elevatorTimeoutChan <- 2
			Println("								timer 2 timeout")
		}
		time.Sleep(10*time.Millisecond)
	}
}











