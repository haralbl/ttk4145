package network

import (
	."fmt"
	"net"
)

const (
	numberOfRetries int = 1 ////////////////////////////// ØK TIL 5
)

func Send(sendChan chan string){

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
		for i:=0; i<numberOfRetries; i++ {
			_, err := socket.Write(data)
			if err != nil {
				Printf("error Send 2")
			}
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
		println(string(data[:256]))
		receiveChan <- string(data[:256])
	}
}








