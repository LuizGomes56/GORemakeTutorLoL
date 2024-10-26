package services

import (
	"fmt"
	"golang/structs"
)

func assign_champion(data *structs.GameProps) {
	allp := data.AllPlayers
	for i := 0; i < len(allp); i++ {
		if allp[i].SummonerName == data.ActivePlayer.SummonerName {
			data.ActivePlayer.Team = allp[i].Team
		}
		// champ, err := go ChampionAPI(allp[i].ChampionName)
		// if err != nil {
		// 	log.Fatalf("Error when calling ChampionAPI: %s", err)
		// }
		// allp[i].Champion = champ
	}
}

func Calculate(data structs.GameProps) {
	assign_champion(&data)
	fmt.Printf("%+v\n", data.ActivePlayer.ChampionStats)
}
