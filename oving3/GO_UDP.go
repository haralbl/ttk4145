package main

import (
		."fmt"
		 "net"
		 "time"
)

func send(doneChan chan string){
	BROADCAST_IPv4 := net.IPv4(129, 241, 187, 255)
	port := 30000
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	if err != nil {
		Printf("error send")}
	
	for {
		data := []byte("Error 51: Connection to server IP failed.")
		
		_, err := socket.Write(data)
		
		if err != nil {
			Printf("error receive 2")}
		
		Printf("sent bytes\n")
		time.Sleep(1000*time.Millisecond)
	}
	doneChan<- "done sending"
}

func receive(doneChan chan string){
	port := 30000
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),	
		Port: port,
	})
	if err != nil {
		Printf("error receive 1")}
	
	for {
		data := make([]byte, 256)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		    
		if err != nil {
			Printf("error receive 2")}
		
		Printf("received %d bytes from %s: %s\n", read, remoteAddr, data)
	}
	doneChan<- "done receiving"
}

func main() {
	doneChan := make(chan string)
	
	go send(doneChan)
	go receive(doneChan)
	
	Println(<-doneChan)
	Println(<-doneChan)
}
