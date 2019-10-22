package main

import (
	"fmt"
	"log"

	"github.com/dim13/redbutton"
)

func main() {
	dev, err := redbutton.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	for ev := range redbutton.Poll(dev, redbutton.PollInterval) {
		fmt.Println(ev)
	}
}
