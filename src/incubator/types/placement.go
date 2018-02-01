package types

type LayerPlacement struct {
	DisplaceX uint64  `json:"displace_x"`
	DisplaceY uint64  `json:"displace_y"`
	ScaleX    float64 `json:"scale_x"`
	ScaleY    float64 `json:"scale_y"`
	Rotate    float64 `json:"rotate"`
}
