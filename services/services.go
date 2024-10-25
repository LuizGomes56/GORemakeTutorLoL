package services

import (
	"fmt"
	"golang/structs"
)

func assign_champion() {
	fmt.Println("hello")
}

func Calculate(data structs.GameProps) {
	fmt.Printf("%+v\n", data.ActivePlayer.ChampionStats)
}
