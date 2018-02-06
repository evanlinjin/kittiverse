package oldincubator

import (
	"fmt"
	"path/filepath"
	"reflect"
)

type DNAGenSpecs struct {
	Group   uint64 `json:"group"`   // Kitty group.
	Color   int32  `json:"color"`   // Fur color.
	Pattern int32  `json:"pattern"` // Fur pattern.
	Body    int32  `json:"body"`
	Brows   int32  `json:"brows"`
	Ears    int32  `json:"ears"`
	Eyes    int32  `json:"eyes"`
	Head    int32  `json:"head"`
	Nose    int32  `json:"nose"`
	Tail    int32  `json:"tail"`
}

type ItemGenSpecs struct {
	ID    int32 `json:"id"`
	Color int32 `json:"color"`
}

type AccessoriesGenSpecs struct {
	Cap    *ItemGenSpecs `json:"cap,omitempty"`
	Collar *ItemGenSpecs `json:"collar,omitempty"`
}

type KittyGenSpecs struct {
	Version     uint64              `json:"version"`
	DNA         DNAGenSpecs         `json:"dna"`
	Accessories AccessoriesGenSpecs `json:"accessories"`

	dnaVal         reflect.Value
	accessoriesVal reflect.Value
	partsDir       string
}

func (s *KittyGenSpecs) Init() error {
	s.dnaVal = reflect.ValueOf(&s.DNA).Elem()
	s.accessoriesVal = reflect.ValueOf(&s.Accessories).Elem()
	s.partsDir = filepath.Join(GetRootDir(), "kitties",
		fmt.Sprintf("group_%d", s.DNA.Group))
	return nil
}

func (s *KittyGenSpecs) GetPartsDir() string {
	return s.partsDir
}

func (s *KittyGenSpecs) GetDNAPartID(ps *KittyPartSpecs) int32 {
	return int32(s.dnaVal.FieldByName(ps.FieldName()).Int())
}

func (s *KittyGenSpecs) SetDNAPartID(ps *KittyPartSpecs, id int32) {
	s.dnaVal.FieldByName(ps.FieldName()).SetInt(int64(id))
}

func (s *KittyGenSpecs) GetAccessory(ps *KittyPartSpecs) *ItemGenSpecs {
	return s.accessoriesVal.FieldByName(ps.FieldName()).Interface().(*ItemGenSpecs)
}

func (s *KittyGenSpecs) SetAccessory(ps *KittyPartSpecs, item *ItemGenSpecs) {
	s.accessoriesVal.Set(reflect.ValueOf(item))
}

type PartPaths struct {
	Outline    string
	OutlineAlt string
	Area       string
	Color      string
}

func (p *PartPaths) HasColor() bool {
	return p.Color != ""
}

func (s *KittyGenSpecs) GetPartPath(ps *KittyPartSpecs) *PartPaths {
	var item *ItemGenSpecs
	if ps.IsAccessory() {
		item = s.GetAccessory(ps)
	} else {
		item = &ItemGenSpecs{
			ID:    s.GetDNAPartID(ps),
			Color: -1,
		}
	}
	if item == nil || item.ID < 0 {
		return nil
	}
	return &PartPaths{
		Outline:    getOutline(s, ps, item.ID),
		OutlineAlt: getOutlineAlt(s, ps, item.ID),
		Area:       getArea(s, ps, item.ID),
		Color:      getColorPath(item.Color),
	}
}

type SkinPaths struct {
	Color   string
	Pattern string
}

func (o *SkinPaths) HasPattern() bool {
	return o.Pattern != ""
}

func (s *KittyGenSpecs) GetSkinPaths() *SkinPaths {
	return &SkinPaths{
		Color:   getColorPath(s.DNA.Color),
		Pattern: getPatternPath(s.DNA.Pattern),
	}
}

/*
	<< HELPER FUNCTIONS >>
*/

func getColorPath(id int32) string {
	if id < 0 {
		return ""
	}
	return filepath.Join(GetRootDir(), "fur", "color",
		fmt.Sprintf("%d.png", id))
}

func getPatternPath(id int32) string {
	if id < 0 {
		return ""
	}
	return filepath.Join(GetRootDir(), "fur", "pattern",
		fmt.Sprintf("%d.png", id))
}

func getOutline(s *KittyGenSpecs, ps *KittyPartSpecs, id int32) string {
	return filepath.Join(s.partsDir, ps.FolderName(),
		fmt.Sprintf("%d_outline.png", id))
}

func getOutlineAlt(s *KittyGenSpecs, ps *KittyPartSpecs, id int32) string {
	return filepath.Join(s.partsDir, ps.FolderName(),
		fmt.Sprintf("%d.png", id))
}

func getArea(s *KittyGenSpecs, ps *KittyPartSpecs, id int32) string {
	return filepath.Join(s.partsDir, ps.FolderName(),
		fmt.Sprintf("%d_area.png", id))
}
