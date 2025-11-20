package main

import (
	"finTrackCLI/commands"
	"log"
)

func main() {
	err := commands.Start()
	if err != nil {
		log.Fatal(err)
	}
}
