package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i int = 0

func goroutine1() {
	for j := 1; j<1000000; j++ {
		i++
	}
}

func goroutine2() {
	for j := 1; j<1000000; j++ {
		i--
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	// Try doing the exercise both with and without it!
	
	go goroutine1()
	go goroutine2()
	
	time.Sleep(100*time.Millisecond)
	
	Println(i)
}
