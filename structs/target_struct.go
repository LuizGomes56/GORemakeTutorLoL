package structs

type RiotChampionTarget struct {
	Id      string                    `json:"id"`
	Name    string                    `json:"name"`
	Spells  []RiotChampionStats       `json:"spells"`
	Stats   []RiotChampionTargetSpell `json:"stats"`
	Passive RiotChampionPassive       `json:"passive"`
}

type RiotChampionTargetSpell struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cooldown    []float64 `json:"cooldown"`
}

type RiotItemTarget struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Stats       RiotItemStats   `json:"stats"`
	Gold        RiotItemGold    `json:"gold"`
	Maps        map[string]bool `json:"maps"`
	From        []string        `json:"from,omitempty"`
}

type AllStatsMultiplier struct {
	Magic    float64 `json:"magic"`
	Physical float64 `json:"physical"`
	General  float64 `json:"general"`
}

type AllStatsAdaptative struct {
	AdaptativeType string  `json:"adaptativeType"`
	Ratio          float64 `json:"ratio"`
}

type TargetReplacements map[string]float64

type AllStatsChampionStats struct {
	MaxHealth               float64 `json:"maxHealth"`
	Armor                   float64 `json:"armor"`
	MagicResist             float64 `json:"magicResist"`
	AttackDamage            float64 `json:"attackDamage"`
	ResourceMax             float64 `json:"resourceMax"`
	AbilityPower            float64 `json:"abilityPower"`
	CurrentHealth           float64 `json:"currentHealth"`
	AttackRange             float64 `json:"attackRange"`
	CritChance              float64 `json:"critChance"`
	CritDamage              float64 `json:"critDamage"`
	PhysicalLethality       float64 `json:"physicalLethality"`
	ArmorPenetrationPercent float64 `json:"armorPenetrationPercent"`
	MagicPenetrationPercent float64 `json:"magicPenetrationPercent"`
	MagicPenetrationFlat    float64 `json:"magicPenetrationFlat"`
}

type TargetChampion struct {
	Id      string              `json:"id"`
	Name    string              `json:"name"`
	Stats   RiotChampionStats   `json:"stats"`
	Spells  RiotChampionSpell   `json:"spells"`
	Passive RiotChampionPassive `json:"passive"`
}

type AllStatsActivePlayer struct {
	ID            string                `json:"id"`
	Level         uint8                 `json:"level"`
	Form          string                `json:"form"`
	Multiplier    AllStatsMultiplier    `json:"multiplier"`
	Adaptative    AllStatsAdaptative    `json:"adaptative"`
	ChampionStats AllStatsChampionStats `json:"championStats"`
	BaseStats     GameCoreStats         `json:"baseStats"`
	BonusStats    GameCoreStats         `json:"bonusStats"`
}

type AllStatsRealStats struct {
	Armor       float64 `json:"armor"`
	MagicResist float64 `json:"magicResist"`
}

type AllStatsPlayer struct {
	Multiplier    AllStatsMultiplier `json:"multiplier"`
	RealStats     AllStatsRealStats  `json:"realStats"`
	ChampionStats GameCoreStats      `json:"championStats"`
	BaseStats     GameCoreStats      `json:"baseStats"`
	BonusStats    GameCoreStats      `json:"bonusStats"`
}

type AllStatsProperty struct {
	OverHealth    float64 `json:"overHealth"`
	MissingHealth float64 `json:"missingHealth"`
	ExcessHealth  float64 `json:"excessHealth"`
	Steelcaps     float64 `json:"steelcaps"`
	Rocksolid     float64 `json:"rocksolid"`
	Randuin       float64 `json:"randuin"`
}

type TargetAllStats struct {
	ActivePlayer AllStatsActivePlayer `json:"activePlayer"`
	Player       AllStatsPlayer       `json:"player"`
	Property     AllStatsProperty     `json:"property"`
}
