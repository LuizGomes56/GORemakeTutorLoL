package main

import (
	"encoding/json"
	"golang/services"
	"golang/structs"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile("test.json")
	if err != nil {
		log.Fatal(err)
	}

	var game structs.GameProps

	err = json.Unmarshal(data, &game)
	if err != nil {
		log.Fatal(err)
	}

	services.Calculate(game)
}
