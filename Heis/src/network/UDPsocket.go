package network

import (
	."fmt"
	"net"
)

func Send(sendChan chan string){ ///////////////////////// husk å spam:
	var upButton int
	var downButton int
	var commandButton int
	var currFloor int


	BROADCAST_IPv4 := net.IPv4(129, 241, 187, 255)
	port := 58017
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	if err != nil {
		Printf("error Send 1")
	}
	
	for {		
		data := []byte(<-sendChan)
	
		_, err := socket.Write(data)
		Printf("sending\n")
		if err != nil {
			Printf("error Send 2")
		}
	}
}

func Receive(receiveChan chan string){
	addr, _ := net.ResolveUDPAddr("udp4", ":58017")
	socket, err := net.ListenUDP("udp4", addr)
	if err != nil {
		Printf("error Receive 1")}
	
	for {
		data := make([]byte, 256)
		_,_,err := socket.ReadFromUDP(data)
		    
		if err != nil {
			Printf("error Receive 2")}
		receiveChan <- string(data[:256])
	}
}








