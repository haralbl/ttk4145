package main

import (
	."fmt"
	"net"
	"time"
)



func count() {
	var counter int = 0
	var data = make([]byte, 256)
	BROADCAST_IPv4 := net.IPv4(129, 241, 187, 255)
	port := 58017
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	if err != nil {
		Printf("error SendAliveMessage 1")
	}
	
	for {
		counter++
		Println(counter)
		
		data[0] = byte(counter)
		_, err := socket.Write(data)
		if err != nil {
			Printf("error SendAliveMessage 2")
		}
		time.Sleep(1000*time.Millisecond)
	}
}

func SendAliveMessage() {
	BROADCAST_IPv4 := net.IPv4(129, 241, 187, 255)
	port := 57017
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP: BROADCAST_IPv4,
		Port: port,
	})
	if err != nil {
		Printf("error SendAliveMessage 1")
	}
	
	for {
		
		data := []byte(GetLocalIP())
		_, err := socket.Write(data)
		if err != nil {
			Printf("error SendAliveMessage 2")
		}
		time.Sleep(100*time.Millisecond)
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

func main() {
	doneChan := make(chan string)
	
	go SendAliveMessage()
	go count()
	
	println(<-doneChan)
}
