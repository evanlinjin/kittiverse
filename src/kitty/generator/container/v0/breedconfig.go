package v0

type LayerShift struct {
}

type LayerTypeName string

type BreedConfig struct {
	BreedName  string
	LayerTypes map[string]LayerShift
}
