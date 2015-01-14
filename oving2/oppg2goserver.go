package main

import (
	. "fmt"
	"runtime"
	//"time"
)

var i 		int 	= 0
var message string 	= "message"

func goroutine1(messages1 chan string, doneChan chan string) {
	for j := 1; j<1000000; j++ {
		messages1 <- "increment i!"
	}
	doneChan<- "done adding"
}

func goroutine2(messages2 chan string, doneChan chan string) {
	for j := 1; j<1000001; j++ {
		messages2 <- "decrement i!"
	}
	doneChan<- "done subtracting"
}

func server(messages1 chan string, messages2 chan string, doneChan chan string) {
	
	for {
		select {
		case message = <- messages1:
			i++
		case message = <- messages2:
			i--
		}
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	
	messages1 := make(chan string, 1)
	messages2 := make(chan string, 1)
	doneChan := make(chan string)
	
	go goroutine1(messages1, doneChan)
	go goroutine2(messages2, doneChan)
	go server(messages1, messages2, doneChan)

	Println(<-doneChan)
	Println(<-doneChan)
	
	Println("i: \t",i)
}
