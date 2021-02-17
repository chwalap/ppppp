package main

import (
	"log"

	"weather/service"
)

func main() {
	log.Printf("Weather service startup\n")
	service.StartService()
}
