package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

const (
	startOffset = 0x00156FBD
	endOffset   = 0x00156FC1
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Please provide the exe file path as an argument.")
	}

	exePath := os.Args[1]
	data, err := ioutil.ReadFile(exePath)
	if err != nil {
		log.Fatal(err)
	}

	newData := make([]byte, len(data))
	offset := startOffset
	for i := range data {
		if offset >= endOffset {
			offset = startOffset
		}
		newData[i] = data[i] + byte(math.Abs(float64(offset-startOffset)))
		offset++
	}

	err = ioutil.WriteFile(exePath, newData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Exe file has been obfuscated.")
}
