package services

import (
	"fmt"
	"golang/functions"
	"golang/structs"
	"strings"
)

var LOCAL_CHAMPION structs.LocalChampion
var LOCAL_ITEMS = functions.FetchFile[structs.LocalItems]("effects/items")
var LOCAL_RUNES = functions.FetchFile[structs.LocalRunes]("effects/runes")

func assign_champion(allPlayers []structs.GamePlayer) {
	for i := range allPlayers {
		val := &allPlayers[i]
		champ := ChampionAPI(&val.ChampionName)
		val.Champion = champ
	}
}

func Calculate(data *structs.GameProps) {
	assign_champion(data.AllPlayers)

	for _, player := range data.AllPlayers {
		if player.SummonerName == data.ActivePlayer.SummonerName {
			data.ActivePlayer.Champion = player.Champion
			data.ActivePlayer.ChampionName = player.ChampionName
			data.ActivePlayer.Skin = player.SkinId

			LOCAL_CHAMPION = functions.FetchFile[structs.LocalChampion](fmt.Sprintf("champions/%s", player.Champion.Id))

			base := structs.Base(player.Champion.Stats, float64(player.Level))

			data.ActivePlayer.BaseStats = base
			data.ActivePlayer.BonusStats = base.Bonus(data.ActivePlayer.ChampionStats.Core())

			{
				items := make([]string, 0, 6)

				for _, item := range player.Items {
					items = append(items, fmt.Sprintf("%d", item.ItemId))
				}

				data.ActivePlayer.Items = items

				data.ActivePlayer.Relevant = structs.GameRelevant{
					Abilities: filter_abilities(),
					Items:     filter_items(items),
					Runes:     filter_runes(data.ActivePlayer.FullRunes.GeneralRunes),
					Spell:     filter_spells(player.SummonerSpells),
				}
			}
			break
		}
	}
	fmt.Printf("%+v\n", data.AllPlayers[0])
}

func filter_spells(spells structs.SummonerSpells) structs.GameRelevantProps {
	min := make([]string, 0, 2)
	max := make([]string, 0)

	valid := func(str string) {
		val := "SummonerDot"
		if strings.Contains(str, val) {
			min = append(min, val)
		}
	}

	valid(spells.SummonerSpellOne.RawDescription)
	valid(spells.SummonerSpellTwo.RawDescription)

	return structs.GameRelevantProps{
		Min: min,
		Max: max,
	}
}

func filter_runes(general []structs.GeneralRunes) structs.GameRelevantProps {
	min := make([]string, 0, 6)
	max := make([]string, 0, 6)

	runes := make([]string, 0, 6)

	for _, val := range general {
		runes = append(runes, string(val.Id))
	}

	for _, key := range runes {
		r, exists := LOCAL_RUNES[key]
		if !exists {
			continue
		}
		if r.Max != nil {
			max = append(max, key)
		}
		min = append(min, key)
	}
	return structs.GameRelevantProps{
		Min: min,
		Max: max,
	}
}

func filter_items(items []string) structs.GameRelevantProps {
	min := make([]string, 0, 7)
	max := make([]string, 0, 7)
	for _, key := range items {
		item, exists := LOCAL_ITEMS[key]
		if !exists {
			continue
		}
		if item.Max != nil {
			max = append(max, key)
		}
		min = append(min, key)
	}
	return structs.GameRelevantProps{
		Min: min,
		Max: max,
	}
}

func filter_abilities() structs.GameRelevantProps {
	min := make([]string, 0, 8)
	max := make([]string, 0, 8)
	for key, val := range LOCAL_CHAMPION {
		if val.Max != nil {
			max = append(max, key)
		}
		min = append(min, key)
	}
	min = append(min, "A", "C")
	return structs.GameRelevantProps{
		Min: min,
		Max: max,
	}
}
