package main

import (
	"log"
	"os"
	"os/exec"

	"dim13.org/redbutton"
)

type stateFn func(redbutton.Event) stateFn

func Init(b redbutton.Event) stateFn {
	if b == redbutton.Enabled {
		log.Println("Ready...")
		return Armed
	}
	return Init
}

func Armed(b redbutton.Event) stateFn {
	if b == redbutton.Pressed {
		log.Println("Go!")
		go Exec(os.Args[1:])
		return Reset
	}
	return Init
}

func Reset(b redbutton.Event) stateFn {
	if b == redbutton.Disabled {
		log.Println("Reset...")
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

	ev := redbutton.Poll(dev, redbutton.PollInterval)
	for state := Init; state != nil; {
		state = state(<-ev)
	}
}
