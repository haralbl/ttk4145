package network

import (
	."fmt"
	"net"
	//"structDefine"
	//"encoding/json"
)

const (
	numberOfRetries int = 1 ////////////////////////////// ØK TIL 5
)

var (
	lastMsgLength int
)

func Send(sendChan chan []byte){

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
		data := <-sendChan
		for i:=0; i<numberOfRetries; i++ {
			_, err := socket.Write(data)
			if err != nil {
				Printf("error Send 2")
			}
			
		}
	}
}

func Receive(receiveChan chan []byte){
	addr, _ := net.ResolveUDPAddr("udp4", ":58017")
	socket, err := net.ListenUDP("udp4", addr)
	//var ReceivedStatus structDefine.ElevatorStatus_t
	if err != nil {
		Printf("error Receive 1")}
	
	for {
		data := make([]byte, 4096)
		length,_,err := socket.ReadFromUDP(data)
		lastMsgLength = length
		/*messageLengthChan <- msgLength*/
		if err != nil {
			Printf("error Receive 2")}
		//println("hola señor")
		//println(string(data[:4096]))

		//json.Unmarshal(data, &ReceivedStatus)
		//receiveChan <- ReceivedStatus
		receiveChan <- data
	}
}

func GetLastMsgLength() int {
	return lastMsgLength
}



