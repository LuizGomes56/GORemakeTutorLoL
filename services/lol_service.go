package services

import (
	"golang/functions"
	"golang/structs"
)

var IDS_CACHE = functions.FetchFile[map[string]map[string]string]("cache/ids")
var ITEM_CACHE = functions.FetchFile[structs.RiotItems]("cache/items")

func get_champion(name *string) string {
	for key, val := range IDS_CACHE {
		for _, v := range val {
			if v == *name {
				return key
			}
		}
	}
	return "TargetDummy"
}

// func ChampionAPI(name *string) (structs.TargetChampion, err) {
// return nil
// }
