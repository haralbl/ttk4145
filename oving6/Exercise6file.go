package main

import (
	."fmt"
	"io/ioutil"
)

func main() {
	
}

func read() {
	var text = []byte
	file, err = ReadFile("alive")
	if err != nil {
		println("error write")
	}
}

func write() {
	var text = []byte
	err = WriteFile("alive", text, ModeDir)
	if err != nil {
		println("error read")
	}
}


