package structs

type LocalStats = map[string]LocalStatsStruct

type LocalStatsStruct struct {
	Name  string
	Stats LocalStatsHashMap
	Stack bool
	Gold  LocalStatsGold
	Maps  map[string]bool
}

type LocalStatsHashMap struct {
	Raw map[string]any
	Mod map[string]any
}

type LocalStatsGold struct {
	Base        uint32
	Purchasable bool
	Total       uint32
	Sell        uint32
}
