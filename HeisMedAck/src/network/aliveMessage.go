package network

import (
	."fmt"
	"net"
	"time"
)

func SendAliveMessage() {
	BROADCAST_IPv4	:= net.IPv4(129, 241, 187, 255)
	port			:= 57017
	socket, err		:= net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	
	if err != nil {
		Printf("error SendAliveMessage 1")
	}
	for {
		data	:= []byte(GetLocalIP())
		_, err	:= socket.Write(data)
		
		if err != nil {
			Printf("error SendAliveMessage 2")
		}
		time.Sleep(100*time.Millisecond)
	}
}

func ReceiveAliveMessage(receiveAliveMessageChan chan string) {
	addr, _		:= net.ResolveUDPAddr("udp4", ":57017")
	socket, err := net.ListenUDP("udp4", addr)
	
	if err != nil {
		Printf("error ReceiveAliveMessage 1")
	}
	for {
		data	:= make([]byte, 256)
		_,_,err := socket.ReadFromUDP(data)
		    
		if err != nil {
			Printf("error ReceiveAliveMessage 2")
		}
		receiveAliveMessageChan <- string(data[:256])
	}
}

func GetLocalIP() (localIP string) {
	addrs, err := net.InterfaceAddrs()
	
    if err != nil {
    	Println(err)
    }
    for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
		    	localIP = ipnet.IP.String()
			}
		}
    }
	return
}










