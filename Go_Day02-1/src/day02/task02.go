package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

func main() {
	console := bufio.NewScanner(os.Stdin)
	arg := os.Args[1:]

	for console.Scan() {
		temple := console.Text()
		command := append(arg[1:], temple)
		cmd := exec.Command(arg[0], command...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			log.Fatal("Invalid command", cmd)
		}
	}

	err := console.Err()
	if err != nil {
		log.Fatal(err)
	}
}
