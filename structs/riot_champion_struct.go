package structs

type RiotChampionPassive struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       struct {
		Full string `json:"full"`
	}
}

type RiotChampionStats struct {
	Hp                   float64 `json:"hp"`
	Hpperlevel           float64 `json:"hpperlevel"`
	Mp                   float64 `json:"mp"`
	Mpperlevel           float64 `json:"mpperlevel"`
	Armor                float64 `json:"armor"`
	Armorperlevel        float64 `json:"armorperlevel"`
	Spellblock           float64 `json:"spellblock"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	Attackrange          float64 `json:"attackrange"`
	Attackspeed          float64 `json:"attackspeed"`
	Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
	Attackdamage         float64 `json:"attackdamage"`
	Attackdamageperlevel float64 `json:"attackdamageperlevel"`
}

type RiotChampionSpell struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cooldown    []float64 `json:"cooldown"`
	Image       struct {
		Full string `json:"full"`
	}
}

type RiotChampionData struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Image struct {
		Full string `json:"full"`
	}
	Skins []struct {
		Num uint8 `json:"num"`
	}
	Stats   RiotChampionStats   `json:"stats"`
	Spells  []RiotChampionSpell `json:"spells"`
	Passive RiotChampionPassive `json:"passive"`
}

type RiotChampion struct {
	Data map[string]RiotChampionData `json:"data"`
}
