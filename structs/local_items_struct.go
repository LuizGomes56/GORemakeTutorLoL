package structs

type Format struct {
	Melee  string `json:"melee"`
	Ranged string `json:"ranged"`
}

type LocalItems = map[string]struct {
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Min    Format    `json:"min"`
	Max    *Format   `json:"max,omitempty"`
	Onhit  bool      `json:"onhit,omitempty"`
	Effect []float64 `json:"effect,omitempty"`
}
