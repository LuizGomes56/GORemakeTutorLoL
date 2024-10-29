package structs

type GameProps struct {
	ActivePlayer GameActivePlayer `json:"activePlayer"`
	AllPlayers   []GamePlayer     `json:"allPlayers"`
	Events       GameEvents       `json:"events"`
	GameData     GameData         `json:"gameData"`
}

type GameData struct {
	GameTime  float64 `json:"gameTime"`
	MapNumber uint8   `json:"mapNumber"`
}

type GameEvents struct {
	Events []struct {
		EventName  string
		KillerName string
		DragonType string
	}
}

type ChampionStats struct {
	AbilityPower            float64 `json:"abilityPower"`
	Armor                   float64 `json:"armor"`
	ArmorPenetrationFlat    float64 `json:"armorPenetrationFlat"`
	ArmorPenetrationPercent float64 `json:"armorPenetrationPercent"`
	AttackDamage            float64 `json:"attackDamage"`
	AttackRange             float64 `json:"attackRange"`
	CritChance              float64 `json:"critChance"`
	CritDamage              float64 `json:"critDamage"`
	CurrentHealth           float64 `json:"currentHealth"`
	MagicPenetrationFlat    float64 `json:"magicPenetrationFlat"`
	MagicPenetrationPercent float64 `json:"magicPenetrationPercent"`
	MagicResist             float64 `json:"magicResist"`
	MaxHealth               float64 `json:"maxHealth"`
	PhysicalLethality       float64 `json:"physicalLethality"`
	ResourceMax             float64 `json:"resourceMax"`
}

func (curr ChampionStats) Core() GameCoreStats {
	return GameCoreStats{
		MaxHealth:    curr.MaxHealth,
		Armor:        curr.Armor,
		MagicResist:  curr.MagicResist,
		AttackDamage: curr.AttackDamage,
		ResourceMax:  curr.ResourceMax,
		AbilityPower: curr.AbilityPower,
	}
}

type GeneralRunes struct {
	DisplayName string `json:"displayName"`
	Id          uint32 `json:"id"`
}

type GameAbilities struct {
	Passive struct {
		DisplayName string `json:"displayName"`
		Id          string `json:"id"`
	}
	Q, W, E, R struct {
		AbilityLevel uint8 `json:"abilityLevel"`
	}
}

type GameActivePlayer struct {
	SummonerName  string        `json:"summonerName"`
	Level         int           `json:"level"`
	Abilities     GameAbilities `json:"abilities"`
	ChampionStats ChampionStats `json:"championStats"`
	FullRunes     struct {
		GeneralRunes []GeneralRunes `json:"generalRunes"`
	} `json:"fullRunes"`
	ExtendsActivePlayer
}

type SummonerSpell struct {
	DisplayName    string `json:"displayName"`
	RawDescription string `json:"rawDescription"`
}

type SummonerSpells struct {
	SummonerSpellOne SummonerSpell `json:"summonerSpellOne"`
	SummonerSpellTwo SummonerSpell `json:"summonerSpellTwo"`
}

type GamePlayer struct {
	ChampionName string `json:"championName"`
	Level        uint8  `json:"level"`
	Position     string `json:"position"`
	SummonerName string `json:"summonerName"`
	Scores       struct {
		Assists uint8 `json:"assists"`
		Kills   uint8 `json:"kills"`
		Deaths  uint8 `json:"deaths"`
	} `json:"scores"`
	Items []struct {
		ItemId uint16 `json:"itemID"`
	}
	SummonerSpells SummonerSpells `json:"summonerSpells"`
	SkinId         uint8          `json:"skinID"`
	Team           string         `json:"team"`
	ExtendsPlayer
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

type GameToolInfo struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Active bool           `json:"active"`
	Gold   uint16         `json:"gold"`
	Raw    map[string]any `json:"raw"`
}
