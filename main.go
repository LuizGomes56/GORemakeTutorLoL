package main

import (
	"golang/functions"
	"golang/services"
	"golang/structs"
)

func main() {
	game := functions.FetchFile[structs.GameProps]("test")

	services.Calculate(game)
}
