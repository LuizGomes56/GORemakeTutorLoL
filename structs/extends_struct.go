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
	ChampionName string         `json:"championName"`
	Champion     TargetChampion `json:"champion"`
	BaseStats    GameCoreStats  `json:"baseStats"`
	BonusStats   GameCoreStats  `json:"bonusStats"`
	Team         string         `json:"team"`
	Tool         string         `json:"tool"`
	Relevant     GameRelevant   `json:"relevant"`
	Skin         uint8          `json:"skin"`
	Items        []string       `json:"items"`
}

type ExtendsPlayer struct {
	Champion      TargetChampion `json:"champion,omitempty"`
	ChampionStats GameCoreStats  `json:"championStats,omitempty"`
	BaseStats     GameCoreStats  `json:"baseStats,omitempty"`
	BonusStats    GameCoreStats  `json:"bonusStats,omitempty"`
	Damage        ExtendsDamage  `json:"damage,omitempty"`
}

type ExtendsPlayerDamage map[string]struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Type  string  `json:"type"`
	Name  string  `json:"name,omitempty"`
	Area  bool    `json:"area,omitempty"`
	Onhit bool    `json:"onhit,omitempty"`
}

type ExtendsDamage struct {
	Abilities ExtendsPlayerDamage `json:"abilities"`
	Runes     ExtendsPlayerDamage `json:"runes"`
	Items     ExtendsPlayerDamage `json:"items"`
	Spell     ExtendsPlayerDamage `json:"spells"`
}
