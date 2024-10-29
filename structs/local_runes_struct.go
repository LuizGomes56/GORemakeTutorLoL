package structs

type LocalRunes = map[string]struct {
	Name string  `json:"name"`
	Type string  `json:"type"`
	Min  Format  `json:"min"`
	Max  *Format `json:"max,omitempty"`
}
