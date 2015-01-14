package main

import (
	. "fmt"
	"runtime"
	//"time"
)

var i int = 0

func goroutine1(messages chan int, doneChan chan string) {
	for j := 1; j<1000000; j++ {
		i = <-messages
		i++
		messages <- i
	}
	doneChan<- "done adding"
}

func goroutine2(messages chan int, doneChan chan string) {
	for j := 1; j<1000001; j++ {
		i = <-messages
		i--
		messages <- i
	}
	doneChan<- "done subtracting"
}

func main() {
	runtime.GOMAXPROCS(4)
	
	messages := make(chan int, 1)
	doneChan := make(chan string)
	messages <- i
	
	go goroutine1(messages, doneChan)
	go goroutine2(messages, doneChan)
	
	
	
	//time.Sleep(10000*time.Millisecond)
	
	Println(<-doneChan)
	Println(<-doneChan)
	
	Println("i: \t",<-messages)
}
