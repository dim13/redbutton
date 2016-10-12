package main

import (
	"fmt"
	"log"

	"dim13.org/redbutton"
)

func main() {
	dev, err := redbutton.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	state := redbutton.Poll(dev)
	for {
		fmt.Println(<-state)
	}
}
