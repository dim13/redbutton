package main

import (
	"log"
	"os"
	"os/exec"

	"dim13.org/redbutton"
)

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
		go Exec(os.Args[1:])
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

func Exec(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ", os.Args[0], " <command>")
	}

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
