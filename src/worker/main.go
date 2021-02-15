package main

import (
	"shared"
	"worker/service"
)

func main() {
	if err := shared.Db.CreateDbIfNotExists(false); err != nil {
		panic(err)
	}
	service.StartService()
}
