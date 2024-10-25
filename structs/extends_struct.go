package structs

type ExtendsActivePlayer struct {
	ChampionName string `json:"championName"`
	// Champion
	BaseStats  GameCoreStats `json:"baseStats"`
	BonusStats GameCoreStats `json:"bonusStats"`
	Team       string        `json:"team"`
	Tool       string        `json:"tool"`
	Relevant   GameRelevant  `json:"relevant"`
}

/*#[derive(Debug, Clone, Deserialize, Default, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct GamePlayerDamage {
    pub min: f64,
    pub max: Option<f64>,
    #[serde(rename = "type")]
    pub damage_type: String,
    pub name: Option<String>,
    pub area: Option<bool>,
    pub onhit: Option<bool>,
}*/

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
