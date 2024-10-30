package structs

type GameCoreStats struct {
	MaxHealth    float64 `json:"maxHealth"`
	Armor        float64 `json:"armor"`
	MagicResist  float64 `json:"magicResist"`
	AttackDamage float64 `json:"attackDamage"`
	ResourceMax  float64 `json:"resourceMax"`
	AbilityPower float64 `json:"abilityPower"`
}

func Base(base RiotChampionStats, level float64) GameCoreStats {
	formula := func(base float64, per_level float64, level float64) float64 {
		return base + per_level*(level-1.0)*(0.7025+0.0175*(level-1.0))
	}
	return GameCoreStats{
		MaxHealth:    formula(base.Hp, base.Hpperlevel, level),
		Armor:        formula(base.Armor, base.Armorperlevel, level),
		MagicResist:  formula(base.Spellblock, base.SpellBlockPerLevel, level),
		AttackDamage: formula(base.Attackdamage, base.Attackdamageperlevel, level),
		ResourceMax:  formula(base.Mp, base.Mpperlevel, level),
		AbilityPower: 0.0,
	}
}

func (base GameCoreStats) Bonus(current GameCoreStats) GameCoreStats {
	return GameCoreStats{
		MaxHealth:    current.MaxHealth - base.MaxHealth,
		Armor:        current.Armor - base.Armor,
		MagicResist:  current.MagicResist - base.MagicResist,
		AttackDamage: current.AttackDamage - base.AttackDamage,
		ResourceMax:  current.ResourceMax - base.ResourceMax,
		AbilityPower: current.AbilityPower,
	}
}

type ExtendsActivePlayer struct {
	ChampionName string         `json:"championName,omitempty"`
	Champion     TargetChampion `json:"champion,omitempty"`
	BaseStats    GameCoreStats  `json:"baseStats,omitempty"`
	BonusStats   GameCoreStats  `json:"bonusStats,omitempty"`
	Team         string         `json:"team,omitempty"`
	Tool         GameToolInfo   `json:"tool,omitempty"`
	Relevant     GameRelevant   `json:"relevant,omitempty"`
	Skin         uint8          `json:"skin,omitempty"`
	Items        []string       `json:"items,omitempty"`
}

type GameToolInfo struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Active bool           `json:"active"`
	Gold   uint16         `json:"gold"`
	Raw    map[string]any `json:"raw"`
}

type GameRelevantProps struct {
	Min []string `json:"min"`
	Max []string `json:"max"`
}

type GameRelevant struct {
	Abilities GameRelevantProps `json:"abilities"`
	Items     GameRelevantProps `json:"items"`
	Runes     GameRelevantProps `json:"runes"`
	Spell     GameRelevantProps `json:"spells"`
}

type ExtendsPlayer struct {
	Champion      TargetChampion `json:"champion,omitempty"`
	ChampionStats GameCoreStats  `json:"championStats,omitempty"`
	BaseStats     GameCoreStats  `json:"baseStats,omitempty"`
	BonusStats    GameCoreStats  `json:"bonusStats,omitempty"`
	Damage        ExtendsDamage  `json:"damage,omitempty"`
	Tool          ExtendsTool    `json:"tool,omitempty"`
}

type ExtendsToolChange struct {
	Dif *ExtendsDamage `json:"dif,omitempty"`
	Sum float64        `json:"sum"`
}

type ExtendsTool struct {
	Max ExtendsDamage       `json:"max"`
	Dif *ExtendsDamage      `json:"dif,omitempty"`
	Sum float64             `json:"sum"`
	Rec *map[string]float64 `json:"rec,omitempty"`
}

type ExtendsPlayerDamage struct {
	Min   float64  `json:"min"`
	Max   *float64 `json:"max"`
	Type  string   `json:"type"`
	Name  *string  `json:"name,omitempty"`
	Area  *bool    `json:"area,omitempty"`
	Onhit *bool    `json:"onhit,omitempty"`
}

type ExtendsDamageReturn = map[string]ExtendsPlayerDamage

type ExtendsDamage struct {
	Abilities ExtendsDamageReturn `json:"abilities"`
	Runes     ExtendsDamageReturn `json:"runes"`
	Items     ExtendsDamageReturn `json:"items"`
	Spell     ExtendsDamageReturn `json:"spell"`
}

func (curr ExtendsDamage) ToHashMap() map[string]ExtendsDamageReturn {
	return map[string]ExtendsDamageReturn{
		"abilities": curr.Abilities,
		"runes":     curr.Runes,
		"items":     curr.Items,
		"spell":     curr.Spell,
	}
}
