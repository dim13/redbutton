package main

import (
	"flag"
	"log"
	"os/exec"
	"strings"

	"dim13.org/redbutton"
)

var cmd = flag.String("cmd", "echo ok", "cmd to run")

type StateFn func(redbutton.Button) StateFn

func Init(b redbutton.Button) StateFn {
	if b == redbutton.Armed {
		log.Println("Ready...")
		return Armed
	}
	return Init
}

func Armed(b redbutton.Button) StateFn {
	if b == redbutton.Pressed {
		log.Println("Go!")
		Exec(*cmd)
		return Reset
	}
	return Init
}

func Reset(b redbutton.Button) StateFn {
	if b == redbutton.Closed {
		return Init
	}
	return Reset
}

func Exec(s string) {
	parts := strings.Fields(s)
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
}

func main() {
	flag.Parse()

	dev, err := redbutton.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	state := redbutton.Poll(dev)
	for stateFn := Init; stateFn != nil; {
		stateFn = stateFn(<-state)
	}
}
