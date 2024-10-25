package structs

type RiotItemGold struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}

type RiotItemStats map[string]float64

type RiotItem struct {
	Name        *string            `json:"name,omitempty"`
	Gold        *RiotItemGold      `json:"gold,omitempty"`
	Description *string            `json:"description,omitempty"`
	Stats       RiotItemStats      `json:"stats"`
	Maps        *map[string]bool   `json:"maps,omitempty"`
	Effect      *map[string]string `json:"effect,omitempty"`
	From        *[]string          `json:"from,omitempty"`
}

type RiotItems struct {
	Data map[string]RiotItem `json:"data"`
}
