package main

import (
	. "fmt"
	"runtime"
	//"time"
)

var (
	i int = 0
	ownership string = "I own the resource"
)

func goroutine1(messages chan string, doneChan chan string) {
	for j := 1; j<1000000; j++ {
		ownership = <-messages
		i++
		messages <- ownership
	}
	doneChan<- "done adding"
}

func goroutine2(messages chan string, doneChan chan string) {
	for j := 1; j<1000001; j++ {
		ownership = <-messages
		i--
		messages <- ownership
	}
	doneChan<- "done subtracting"
}

func main() {
	runtime.GOMAXPROCS(4)
	
	messages := make(chan string, 1)
	doneChan := make(chan string)
	messages <- ownership
	
	go goroutine1(messages, doneChan)
	go goroutine2(messages, doneChan)
	
	
	
	//time.Sleep(10000*time.Millisecond)
	
	Println(<-doneChan)
	Println(<-doneChan)
	
	Println("i: \t",i)
}
