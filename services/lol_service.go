package services

import (
	"fmt"
	"golang/functions"
	"golang/structs"
)

var IDS_CACHE = functions.FetchFile[map[string]map[string]string]("cache/ids")
var ITEM_CACHE = functions.FetchFile[structs.RiotItems]("cache/item")
var CHAMPION_CACHE = make(map[string]structs.TargetChampion)

func ItemAPI(name *string) structs.RiotItem {
	path, exists := ITEM_CACHE.Data[*name]

	if !exists {
		panic(fmt.Sprintf("ItemAPI -> Item %s not found", *name))
	}

	return path
}

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

func ChampionAPI(name *string) structs.TargetChampion {
	id := get_champion(name)
	for key, val := range CHAMPION_CACHE {
		if key == id {
			return val
		}
	}
	data := functions.FetchFile[structs.RiotChampion](fmt.Sprintf("cache/champions/%s", id))

	path := data.Data[id]

	target := structs.TargetChampion{
		Id:      id,
		Name:    path.Name,
		Stats:   path.Stats,
		Spells:  path.Spells,
		Passive: path.Passive,
	}

	CHAMPION_CACHE[id] = target
	return target
}
