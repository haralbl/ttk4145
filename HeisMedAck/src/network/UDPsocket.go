package network

import (
	."fmt"
	"net"
	"time"
)

const (
	numberOfRetries int = 5
)

var (
	lastMsgLength int
)

func SendManager(sendChan chan []byte, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan chan []byte) {
	for {
		data := <-sendChan
		go send(data, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan)
	}
}

func send(data []byte, checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan chan []byte) {
	BROADCAST_IPv4	:= net.IPv4(129, 241, 187, 255)
	port			:= 58017
	socket, err		:= net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	if err != nil {
		Printf("error Send 1")
	}
	for i:=0; i<numberOfRetries; i++ {
		_, err := socket.Write(data)
		
		if err != nil {
			Printf("error Send 2")
		}
		time.Sleep(10*time.Millisecond)
	}
	checkIfOrderIsAddedToQueueAndPotentiallyTakeTheOrderMyselfIfNotAddedChan <- data
}

func Receive(receiveChan chan []byte) {
	addr, _		:= net.ResolveUDPAddr("udp4", ":58017")
	socket, err := net.ListenUDP("udp4", addr)

	if err != nil {
		Printf("error Receive 1")
	}
	for {
		data			:= make([]byte, 1024)
		length, _, err	:= socket.ReadFromUDP(data)
		
		data[1000] = byte(length%10)
		data[1001] = byte((length/10)%10)
		data[1002] = byte((length/100)%10)

		lastMsgLength = length

		if err != nil {
			Printf("error Receive 2")
		}
		receiveChan <- data
	}
}













