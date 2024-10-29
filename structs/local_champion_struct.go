package structs

type LocalChampion = map[string]struct {
	Type string   `json:"type"`
	Area bool     `json:"area"`
	Min  []string `json:"min"`
	Max  []string `json:"max,omitempty"`
}
