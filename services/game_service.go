package services

import (
	"fmt"
	"golang/functions"
	"golang/structs"
	"math"
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
				items := make([]string, 0, len(player.Items))

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

	for i, player := range data.AllPlayers {
		if player.Team != data.ActivePlayer.Team {
			base := structs.Base(player.Champion.Stats, float64(player.Level))

			items := make([]string, 0, len(player.Items))

			for _, item := range player.Items {
				items = append(items, fmt.Sprintf("%d", item.ItemId))
			}

			player.BaseStats = base

			champion_stats := player_stats(base, items)

			player.ChampionStats = champion_stats
			player.BonusStats = base.Bonus(champion_stats)

			stats := all_stats(player, data.ActivePlayer)

			// relv := data.ActivePlayer.Relevant

			// player.Damage = structs.ExtendsDamage{
			// 	Abilities: ability_damage(stats, data.ActivePlayer.Abilities),
			// 	Runes:     rune_damage(stats, relv.Runes.Min),
			// 	Items:     item_damage(stats, relv.Items.Min),
			// 	Spell:     spell_damage(player.Level, relv.Spell.Min),
			// }

			if i == 0 {
				fmt.Printf("%+v\n", stats)
			}
		}
	}
	// fmt.Printf("%+v\n", data.AllPlayers[0])
}

func ability_damage(stats structs.TargetAllStats, abilities structs.GameAbilities) {

}

func all_stats(player structs.GamePlayer, activePlayer structs.GameActivePlayer) structs.TargetAllStats {
	acs := activePlayer.ChampionStats
	abs := activePlayer.BonusStats
	abt := activePlayer.BaseStats

	pcs := player.ChampionStats
	pbs := player.BonusStats
	pbt := player.BaseStats

	acpMod, pphyMod, pmagMod, pgenMod := 1.0, 1.0, 1.0, 1.0

	rar := math.Max(0, pcs.Armor*acs.ArmorPenetrationPercent-acs.ArmorPenetrationFlat)
	rmr := math.Max(0, pcs.MagicResist*acs.MagicPenetrationPercent-acs.MagicPenetrationFlat)

	physical := 100 / (100 + rar)
	magic := 100 / (100 + rmr)

	adp := 0.35*abs.AttackDamage >= 0.2*acs.AbilityPower

	add := 0.0

	if adp {
		add = physical
	} else {
		add = magic
	}

	ohp := pcs.MaxHealth / acs.MaxHealth
	ehp := pcs.MaxHealth - acs.MaxHealth
	mshp := 1 - acs.CurrentHealth/acs.MaxHealth
	exhp := 1.0

	if ehp > 2500 {
		exhp = 2500
	} else if ehp < 0 {
		exhp = 0
	}

	rel := activePlayer.Relevant

	if functions.Includes(rel.Runes.Min, "8299") {
		x := 0.0
		if mshp > 0.7 {
			x = 0.11
		} else if mshp >= 0.4 {
			x = 0.2*mshp - 0.03
		}
		acpMod += x
	}

	if functions.Includes(rel.Items.Min, "4015") {
		acpMod += exhp / (220000 / 15)
	}

	var form string

	if acs.AttackRange > 350 {
		form = "ranged"
	} else {
		form = "melee"
	}

	var adaptative string

	if adp {
		adaptative = "physical"
	} else {
		adaptative = "magic"
	}

	stcaps, rand, rock, overhp := 1.0, 1.0, 1.0, 0.0

	items := make([]string, 0, len(player.Items))

	for _, item := range player.Items {
		items = append(items, fmt.Sprintf("%d", item.ItemId))
	}

	if functions.Includes(items, "3143") {
		rand = 0.7
	}

	if functions.Includes(items, "3143", "3110", "3082") {
		rock += (pcs.MaxHealth / 1000) * 3.5
	}

	if functions.Includes(items, "3047") {
		stcaps = 0.88
	}

	if ohp < 1.1 {
		overhp = 0.65
	} else if ohp > 2 {
		overhp = 2
	}

	return structs.TargetAllStats{
		ActivePlayer: structs.AllStatsActivePlayer{
			ID:         activePlayer.Champion.Id,
			Level:      uint8(activePlayer.Level),
			Form:       form,
			Multiplier: structs.AllStatsMultiplier{Magic: magic, Physical: physical, General: acpMod},
			Adaptative: structs.AllStatsAdaptative{AdaptativeType: adaptative, Ratio: add},
			ChampionStats: structs.AllStatsChampionStats{
				AbilityPower: acs.AbilityPower,
				AttackDamage: acs.AttackDamage,
				AttackRange:  acs.AttackRange,
				Armor:        acs.Armor,
				MagicResist:  acs.MagicResist,
				CritChance:   acs.CritChance,
				CritDamage:   acs.CritDamage,
				MaxHealth:    acs.MaxHealth,
			},
			BaseStats:  abt,
			BonusStats: abs,
		},
		Player: structs.AllStatsPlayer{
			Multiplier:    structs.AllStatsMultiplier{Magic: pmagMod, Physical: pphyMod, General: pgenMod},
			RealStats:     structs.AllStatsRealStats{MagicResist: rmr, Armor: rar},
			ChampionStats: pcs,
			BaseStats:     pbt,
			BonusStats:    pbs,
		},
		Property: structs.AllStatsProperty{
			ExcessHealth:  exhp,
			MissingHealth: mshp,
			OverHealth:    overhp,
			Steelcaps:     stcaps,
			Rocksolid:     rock,
			Randuin:       rand,
		},
	}
}

func player_stats(base structs.GameCoreStats, items []string) structs.GameCoreStats {
	for _, item := range items {
		res := ItemAPI(&item)
		stats := res.Stats
		for key, val := range stats {
			switch key {
			case "FlatHPPoolMod":
				base.MaxHealth += val
			case "FlatMPPoolMod":
				base.ResourceMax += val
			case "FlatMagicDamageMod":
				base.AbilityPower += val
			case "FlatArmorMod":
				base.Armor += val
			case "FlatSpellBlockMod":
				base.MagicResist += val
			case "FlatPhysicalDamageMod":
				base.AbilityPower += val
			default:
				continue
			}
		}
	}
	return base
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
	min := make([]string, 0, len(items))
	max := make([]string, 0, len(items))
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
