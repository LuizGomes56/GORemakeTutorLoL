package services

import (
	"fmt"
	"golang/functions"
	"golang/structs"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

var LOCAL_CHAMPION structs.LocalChampion
var LOCAL_ITEMS = functions.FetchFile[structs.LocalItems]("effects/items")
var LOCAL_RUNES = functions.FetchFile[structs.LocalRunes]("effects/runes")
var LOCAL_STATS = functions.FetchFile[structs.LocalStats]("cache/stats")

func assign_champion(allPlayers []structs.GamePlayer) {
	for i := range allPlayers {
		val := &allPlayers[i]
		champ := ChampionAPI(&val.ChampionName)
		val.Champion = champ
	}
}

func Calculate(data *structs.GameProps, tool_item string) structs.GameProps {
	assign_champion(data.AllPlayers)

	for _, player := range data.AllPlayers {
		if player.SummonerName == data.ActivePlayer.SummonerName {
			data.ActivePlayer.Team = player.Team
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

			{
				path, exists := LOCAL_STATS[tool_item]
				if !exists {
					log.Fatal("Error at func Calculate() on tool_item; LOCAL_STATS: " + tool_item)
				}
				_, active := LOCAL_ITEMS[tool_item]

				data.ActivePlayer.Tool = structs.GameToolInfo{
					Id:     tool_item,
					Name:   path.Name,
					Active: active,
					Gold:   uint16(path.Gold.Total),
					Raw:    path.Stats.Raw,
				}
			}
			break
		}
	}

	for i := range data.AllPlayers {
		player := &data.AllPlayers[i]
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

			stats := all_stats(*player, data.ActivePlayer)

			relv := data.ActivePlayer.Relevant

			player.Damage = structs.ExtendsDamage{
				Abilities: ability_damage(stats, data.ActivePlayer.Abilities),
				Runes:     rune_damage(stats, relv.Runes.Min),
				Items:     item_damage(stats, relv.Items.Min),
				Spell:     spell_damage(player.Level, relv.Spell.Min),
			}

			cloned, err := functions.StructuredClone(data.ActivePlayer)
			if err != nil {
				log.Fatal("Erro ao clonar os dados de ActivePlayer: ", err)
			}
			player.Tool = tool_damage(cloned, *player, tool_item)
		}
	}

	return *data
}

func assign_stats(item string, active_player *structs.GameActivePlayer) map[string]float64 {
	stats := active_player.ChampionStats.IntoHashMap()

	if res, ok := LOCAL_STATS[item]; ok {
		for key, val := range res.Stats.Mod {
			if _, ok := stats[key]; ok {
				if v, ok := val.(float64); ok {
					stats[key] += v
				} else if v, ok := val.(string); ok {
					v = strings.ReplaceAll(v, "%", "")
					if r, err := strconv.ParseFloat(v, 64); err == nil {
						stats[key] -= r
					}
				}
			}
		}
	}
	return stats
}

func evaluate_change(next structs.ExtendsPlayerDamage, curr structs.ExtendsPlayerDamage) structs.ExtendsPlayerDamage {
	var max float64
	if next.Max != nil && curr.Max != nil {
		max = *next.Max - *curr.Max
	}

	return structs.ExtendsPlayerDamage{
		Min:   next.Min - curr.Min,
		Max:   &max,
		Type:  next.Type,
		Name:  next.Name,
		Area:  next.Area,
		Onhit: next.Onhit,
	}
}

func process_change(at string, val structs.ExtendsDamageReturn, min_at structs.ExtendsDamageReturn, change_at *structs.ExtendsDamageReturn, sum *float64) {
	for k, v := range val {
		if curr, ok := min_at[k]; ok {
			result := evaluate_change(v, curr)
			*sum += result.Min
			if result.Max != nil {
				*sum += *result.Max
			}
			(*change_at)[k] = result
		} else {
			fmt.Printf("Key is %s, and was not found in both directions on %s\n", k, at)
			continue
		}
	}
}

func find_change(max structs.ExtendsDamage, min structs.ExtendsDamage, sum *float64) *structs.ExtendsDamage {
	change := structs.ExtendsDamage{
		Abilities: structs.ExtendsDamageReturn{},
		Runes:     structs.ExtendsDamageReturn{},
		Items:     structs.ExtendsDamageReturn{},
		Spell:     structs.ExtendsDamageReturn{},
	}
	min_hashmap := min.ToHashMap()
	max_hashmap := max.ToHashMap()

	for key, val := range max_hashmap {
		switch key {
		case "abilities":
			process_change(key, val, min_hashmap[key], &change.Abilities, sum)
		case "runes":
			process_change(key, val, min_hashmap[key], &change.Runes, sum)
		case "items":
			process_change(key, val, min_hashmap[key], &change.Items, sum)
		case "spell":
			process_change(key, val, min_hashmap[key], &change.Spell, sum)
		}
	}

	return &change
}

func tool_change(max structs.ExtendsDamage, min structs.ExtendsDamage) structs.ExtendsToolChange {
	sum := 0.0
	dif := find_change(max, min, &sum)
	return structs.ExtendsToolChange{
		Dif: dif,
		Sum: sum,
	}
}

func tool_damage(active_player structs.GameActivePlayer, player structs.GamePlayer, item string) structs.ExtendsTool {
	assigned_stats := assign_stats(item, &active_player)

	current_stats := structs.FromHashMapCamel(assigned_stats)

	active_player.ChampionStats = current_stats
	active_player.BonusStats = active_player.BaseStats.Bonus(active_player.ChampionStats.Core())

	stats := all_stats(player, active_player)

	max := structs.ExtendsDamage{
		Abilities: ability_damage(stats, active_player.Abilities),
		Runes:     rune_damage(stats, active_player.Relevant.Runes.Min),
		Items:     item_damage(stats, active_player.Relevant.Items.Min),
		Spell:     spell_damage(player.Level, active_player.Relevant.Spell.Min),
	}

	change := tool_change(max, player.Damage)

	result := structs.ExtendsTool{
		Sum: change.Sum,
		Dif: change.Dif,
		Max: max,
	}
	return result
}

func rune_damage(stats structs.TargetAllStats, runes []string) map[string]structs.ExtendsPlayerDamage {
	result := make(map[string]structs.ExtendsPlayerDamage, 8)
	form := stats.ActivePlayer.Form
	for _, r := range runes {
		element, exists := LOCAL_RUNES[r]
		if !exists {
			continue
		}
		var min_str string
		switch form {
		case "melee":
			min_str = element.Min.Melee
		case "ranged":
			min_str = element.Min.Ranged
		default:
			log.Fatal("Error at func rune_damage() on rune: Form is not defined." + r)
		}
		min, _ := evaluate(min_str, nil, stats, nil)

		result[r] = structs.ExtendsPlayerDamage{
			Min:  min,
			Type: element.Type,
			Name: &element.Name,
		}
	}
	return result
}

func item_damage(stats structs.TargetAllStats, items []string) map[string]structs.ExtendsPlayerDamage {
	result := make(map[string]structs.ExtendsPlayerDamage, 8)
	form := stats.ActivePlayer.Form
	for _, r := range items {
		element, exists := LOCAL_ITEMS[r]
		if !exists {
			continue
		}
		var min_str string
		var max_str *string
		var total *float64
		if element.Effect != nil {
			total = &element.Effect[stats.ActivePlayer.Level-1]
		}
		switch form {
		case "melee":
			min_str = element.Min.Melee
		case "ranged":
			min_str = element.Min.Ranged
		default:
			log.Fatal("Error at func rune_damage() on rune: Form is not defined." + r)
		}
		extras := make(map[string]float64, 1)
		if total != nil {
			extras["total"] = *total
		}
		min, max := evaluate(min_str, max_str, stats, &extras)
		result[r] = structs.ExtendsPlayerDamage{
			Min:   min,
			Max:   max,
			Type:  element.Type,
			Name:  &element.Name,
			Onhit: &element.Onhit,
		}
	}
	return result
}

func spell_damage(level uint8, spells []string) map[string]structs.ExtendsPlayerDamage {
	result := make(map[string]structs.ExtendsPlayerDamage, 1)
	for _, spell := range spells {
		if spell == "SummonerDot" {
			name := "Ignite"
			result[spell] = structs.ExtendsPlayerDamage{
				Min:  50.0 + 20.0*float64(level),
				Type: "true",
				Name: &name,
			}
		}
	}
	return result
}

func json_replacements(stats structs.TargetAllStats) structs.TargetReplacements {
	x := stats.ActivePlayer
	y := stats.Player
	z := stats.Property
	k := x.ChampionStats
	t := x.BaseStats
	n := x.BonusStats
	m := y.ChampionStats

	type Entry struct {
		Key   string
		Value float64
	}

	entries := []Entry{
		{"steelcapsEffect", z.Steelcaps},
		{"attackReductionEffect", z.Rocksolid},
		{"exceededHP", z.ExcessHealth},
		{"missingHP", z.MissingHealth},
		{"magicMod", x.Multiplier.Magic},
		{"physicalMod", x.Multiplier.Physical},
		{"level", float64(x.Level)},
		{"currentAP", k.AbilityPower},
		{"currentAD", k.AttackDamage},
		{"currentLethality", k.PhysicalLethality},
		{"maxHP", k.MaxHealth},
		{"maxMana", k.ResourceMax},
		{"currentMR", k.MagicResist},
		{"currentArmor", k.Armor},
		{"currentHealth", k.CurrentHealth},
		{"basicAttack", 1.0},
		{"attackSpeed", 1.0},
		{"critChance", k.CritChance},
		{"critDamage", k.CritDamage},
		{"adaptative", x.Adaptative.Ratio},
		{"baseHP", t.MaxHealth},
		{"baseMana", t.ResourceMax},
		{"baseArmor", t.Armor},
		{"baseMR", t.MagicResist},
		{"baseAD", t.AttackDamage},
		{"bonusAD", n.AttackDamage},
		{"bonusHP", n.MaxHealth},
		{"bonusArmor", n.Armor},
		{"bonusMR", n.MagicResist},
		{"expectedHealth", m.MaxHealth},
		{"expectedMana", m.ResourceMax},
		{"expectedArmor", m.Armor},
		{"expectedMR", m.MagicResist},
		{"expectedAD", m.AttackDamage},
		{"expectedBonusHealth", y.BonusStats.MaxHealth},
	}

	result := make(map[string]float64, len(entries)+1)
	for _, entry := range entries {
		result[entry.Key] = entry.Value
	}
	return result
}

func evaluate(min_str string, max_str *string, stats structs.TargetAllStats, extra *map[string]float64) (float64, *float64) {
	replacemets := json_replacements(stats)

	if extra != nil {
		for key, val := range *extra {
			replacemets[key] = val
		}
	}

	for key, val := range replacemets {
		min_str = strings.ReplaceAll(min_str, key, fmt.Sprintf("%f", val))
		if max_str == nil {
			continue
		}
		*max_str = strings.ReplaceAll(*max_str, key, fmt.Sprintf("%f", val))
	}

	min_expr, err := govaluate.NewEvaluableExpression(min_str)
	if err != nil {
		log.Fatal("Error at func evaluate() on min_expr: " + err.Error())
	}
	min_eval, err := min_expr.Evaluate(nil)
	if err != nil {
		log.Fatal("Error at func evaluate() on min: " + err.Error())
	}

	min, ok := min_eval.(float64)
	if !ok {
		log.Fatal("Error at func evaluate() on min: expected float64")
	}

	if max_str == nil {
		return min, nil
	}

	max_expr, err := govaluate.NewEvaluableExpression(*max_str)
	if err != nil {
		log.Fatal("Error at func evaluate() on max_expr: " + err.Error())
	}
	max_eval, err := max_expr.Evaluate(nil)
	if err != nil {
		log.Fatal("Error at func evaluate() on max: " + err.Error())
	}

	max, ok := max_eval.(float64)
	if !ok {
		log.Fatal("Error at func evaluate() on max: expected float64")
	}

	return min, &max
}

func ability_damage(stats structs.TargetAllStats, abilities structs.GameAbilities) map[string]structs.ExtendsPlayerDamage {
	res := make(map[string]structs.ExtendsPlayerDamage, 8)
	for key, val := range LOCAL_CHAMPION {
		letter := key[0]
		var index uint8
		switch letter {
		case 'Q':
			index = abilities.Q.AbilityLevel - 1
		case 'W':
			index = abilities.W.AbilityLevel - 1
		case 'E':
			index = abilities.E.AbilityLevel - 1
		case 'R':
			index = abilities.R.AbilityLevel - 1
		case 'P':
			index = stats.ActivePlayer.Level - 1
		default:
			log.Fatal("Error at func ability_damage() on key: " + key)
		}
		if index == 0 {
			continue
		}

		var min_str string
		var max_str *string

		if int(index) < len(val.Min) {
			min_str = val.Min[index]
		}

		if int(index) < len(val.Max) {
			max_str = &val.Max[index]
		}

		min, max := evaluate(min_str, max_str, stats, nil)

		res[key] = structs.ExtendsPlayerDamage{
			Min:  min,
			Max:  max,
			Type: val.Type,
			Area: &val.Area,
		}
	}

	acst := stats.ActivePlayer.ChampionStats
	attack := acst.AttackDamage * stats.ActivePlayer.Multiplier.Physical
	res["A"] = structs.ExtendsPlayerDamage{
		Min:  attack,
		Max:  nil,
		Type: "physical",
	}
	res["C"] = structs.ExtendsPlayerDamage{
		Min:  attack * acst.CritDamage / 100.0,
		Max:  nil,
		Type: "physical",
	}
	return res
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
