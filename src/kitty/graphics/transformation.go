package layer

type Transformation struct {
	ShiftX uint16  `json:"shift_x"`
	ShiftY uint16  `json:"shift_y"`
	ScaleX float32 `json:"scale_x"`
	ScaleY float32 `json:"scale_y"`
}
