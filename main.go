package main

import (
	"getcet/cet"
	"log"
)

func main() {
	log.Printf("Starting CET server on port 5400 ...\n")
	log.Fatalln(cet.NewCETServer(":5400", nil))
}
