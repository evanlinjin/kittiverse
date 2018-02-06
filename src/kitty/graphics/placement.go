package layer

type Placement struct {
	CoordX  uint64  `json:"coord_x"` // x coordinate
	CoordY  uint64  `json:"coord_y"` // y coordinate
	ScaleX  float64 `json:"scale_x"` // x scale factor
	ScaleY  float64 `json:"scale_y"` // y scale factor
	Rotate  float64 `json:"rotate"`  // clockwise rotation in radians
	Opacity float64 `json:"opacity"` // opacity factor 0.0 - 1.0 (TODO: Implement)
}
