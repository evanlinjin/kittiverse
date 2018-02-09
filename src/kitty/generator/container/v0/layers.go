package v0

import (
	"github.com/kittycash/kittiverse/src/kitty/generator/container"
	"github.com/kittycash/kittiverse/src/kitty/generator/container/common"
	"github.com/kittycash/kittiverse/src/kitty/genetics"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"image"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	PrefixAccessory = "accessory"
)

type Layers struct {
	LayerTypes       []LayersOfType
	Breeds           []string
	layerTypesByName map[string]*LayersOfType `enc:"-"`
	breedsByName     map[string]struct{}      `enc:"-"`
}

func NewLayersContainer() *Layers {
	return &Layers{
		layerTypesByName: make(map[string]*LayersOfType),
		breedsByName:     make(map[string]struct{}),
	}
}

func (lc *Layers) Version() uint16 {
	return version
}

func (lc *Layers) Import(raw []byte) error {
	// Check raw size.
	if len(raw) < common.VersionLen {
		return common.ErrInvalidSize
	}
	// Check version.
	var ver uint16
	if e := encoder.DeserializeRaw(raw[:common.VersionLen], &ver); e != nil {
		return e
	}
	if ver != version {
		return common.ErrInvalidVersion
	}
	// Load data.
	if e := encoder.DeserializeRaw(raw[common.VersionLen:], lc); e != nil {
		return e
	}
	// Prepare maps.
	lc.layerTypesByName = make(map[string]*LayersOfType)
	for i, v := range lc.LayerTypes {
		lc.layerTypesByName[v.OfType] = &lc.LayerTypes[i]
	}
	lc.breedsByName = make(map[string]struct{})
	for _, v := range lc.Breeds {
		lc.breedsByName[v] = struct{}{}
	}
	return nil
}

func (lc *Layers) Export() []byte {
	return append(encoder.Serialize(version), encoder.Serialize(lc)...)
}

func (lc *Layers) Compile(rootDir string, images container.Images) error {
	// Get layer types.
	if e := initLayerTypes(lc, rootDir); e != nil {
		return e
	}
	// Get breeds.
	if e := initBreeds(lc, rootDir); e != nil {
		return e
	}
	// Get layers.
	if e := initLayers(lc, rootDir, images); e != nil {
		return e
	}
	return nil
}

func (lc *Layers) GetAlleleRanges() genetics.AlleleRanges {
	getRange := func(pos genetics.DNAPos) genetics.AlleleRange {
		return genetics.AlleleRange{
			Min: 0,
			Max: uint16(len(lc.layerTypesByName[pos.String()].Attributes) - 1),
		}
	}
	return genetics.AlleleRanges{
		Breed: genetics.AlleleRange{
			Min: 0,
			Max: uint16(len(lc.Breeds) - 1),
		},
		BodyAttribute: getRange(genetics.DNABodyAttrPos),
		BodyColorA:    getRange(genetics.DNABodyColorAPos),
		BodyColorB:    getRange(genetics.DNABodyColorBPos),
		BodyPattern:   getRange(genetics.DNABodyPatternPos),
		EarsAttribute: getRange(genetics.DNAEarsAttrPos),
		EyesAttribute: getRange(genetics.DNAEyesAttrPos),
		EyesColor:     getRange(genetics.DNAEyesColorPos),
		NoseAttribute: getRange(genetics.DNANoseAttrPos),
		TailAttribute: getRange(genetics.DNATailAttrPos),
	}
}

func (lc *Layers) GenerateKitty(dna genetics.DNA) (image.Image, error) {
	out := image.NewRGBA(image.Rect(0, 0, common.XpxLen, common.YpxLen))
	// TODO: Complete.
	return out, nil
}

/*
	<<< MEMBER HELPERS >>>
*/

func (lc *Layers) addLayerType(ltName string) error {
	if _, has := lc.layerTypesByName[ltName]; has {
		return common.ErrAlreadyExists
	}
	lc.LayerTypes = append(lc.LayerTypes, LayersOfType{
		OfType:           ltName,
		layersByKey:      make(map[attributeKey]*Layer),
		attributesByName: make(map[string]struct{}),
	})
	lc.layerTypesByName[ltName] = &lc.LayerTypes[len(lc.LayerTypes)-1]
	return nil
}

func (lc *Layers) addBreed(bName string) error {
	if _, has := lc.breedsByName[bName]; has {
		return common.ErrAlreadyExists
	}
	lc.Breeds = append(lc.Breeds, bName)
	lc.breedsByName[bName] = struct{}{}
	return nil
}

/*
	<<< HELPERS >>>
*/

func initLayerTypes(lc *Layers, rootDir string) error {
	subDirs, e := ioutil.ReadDir(rootDir)
	if e != nil {
		return e
	}
	for _, dir := range subDirs {
		if dir.IsDir() == false {
			continue
		}
		dirName := dir.Name()
		//if strings.HasPrefix(dirName, PrefixAccessory+"_") {
		//	continue
		//}
		if e := lc.addLayerType(dirName); e != nil {
			return e
		}
	}
	return nil
}

func initBreeds(lc *Layers, rootDir string) error {
	for _, lt := range lc.LayerTypes {
		subDirs, e := ioutil.ReadDir(path.Join(rootDir, lt.OfType))
		if e != nil {
			return e
		}
		for _, dir := range subDirs {
			if dir.IsDir() == false {
				continue
			}
			if e := lc.addBreed(dir.Name()); e != nil {
				return e
			}
		}
	}
	return nil
}

func initLayers(lc *Layers, rootDir string, images container.Images) error {
	for _, lt := range lc.LayerTypes {
		ltDir := path.Join(rootDir, lt.OfType)
		bDirs, e := ioutil.ReadDir(ltDir)
		if e != nil {
			return e
		}
		for _, bDir := range bDirs {
			if bDir.IsDir() == false {
				continue
			}
			breed := bDir.Name()
			// collect attributes.
			files, e := ioutil.ReadDir(path.Join(ltDir, bDir.Name()))
			if e != nil {
				return e
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if strings.HasSuffix(file.Name(), ".png") == false {
					continue
				}
				var (
					fullPath  = path.Join(rootDir, lt.OfType, breed, file.Name())
					fullName  = strings.TrimSuffix(file.Name(), ".png")
					splitName = strings.Split(fullName, "_")

					attributeName = splitName[0]
					partIndex     = 0
					isArea        = false
					isOutline     = false
				)
				for i := 1; i < len(splitName); i++ {
					v := splitName[i]
					switch {
					case v == "left":
						partIndex = 0
					case v == "right":
						partIndex = 1
					case strings.HasPrefix(v, "part"):
						partIndex = getPartIndex(v)
					case v == "area":
						isArea = true
					case v == "outline":
						isOutline = true
					}
				}
				// Checks.
				if isArea && isOutline || !isArea && !isOutline {
					log.WithField("layer_type", lt.OfType).
						WithField("breed", breed).
						WithField("full_attribute", fullName).
						Error("invalid 'area' and 'outline' combination")
					continue
				}
				// Extract image.
				f, e := os.Open(fullPath)
				if e != nil {
					log.WithError(e).Error("failed to open image file")
				}
				imgRaw, e := ioutil.ReadAll(f)
				if e != nil {
					log.WithError(e).Error("failed to read image file")
				}
				imgHash, e := images.Add(imgRaw)
				if e != nil {
					log.WithField("layer_type", lt.OfType).
						WithField("breed", breed).
						WithField("full_attribute", fullName).
						WithError(e).Error("failed to add image to container")
					continue
				}
				// Append.
				layer, ok := lt.Get(newAttributeKey(attributeName, breed))
				if !ok {
					layer, e = lt.addLayer(Layer{
						OfAttribute: attributeName,
						OfBreed:     breed,
					})
					if e != nil {
						log.WithField("layer_type", lt.OfType).
							WithField("breed", breed).
							WithField("full_attribute", fullName).
							WithError(e).Error("failed to add attribute")
						continue
					}
				}
				layer.ensurePartsCount(partIndex + 1)
				switch {
				case isArea:
					layer.Parts[partIndex][0] = imgHash
				case isOutline:
					layer.Parts[partIndex][1] = imgHash
				}
				// Ensure attribute.
				lt.addAttribute(attributeName)
			}
		}
	}
	return nil
}

func getPartIndex(str string) int {
	p := strings.TrimPrefix(str, "part")
	return int([]byte(p)[0] - 65)
}

/*
	<<< TYPES >>>
*/

type LayersOfType struct {
	OfType           string
	Layers           []Layer
	Attributes       []string
	layersByKey      map[attributeKey]*Layer `enc:"-"` // aka, by attribute and breed
	attributesByName map[string]struct{}     `enc:"-"`
}

func (lt *LayersOfType) Init() {
	lt.layersByKey = make(map[attributeKey]*Layer)
	for i, v := range lt.Layers {
		lt.layersByKey[v.key()] = &lt.Layers[i]
	}
	lt.attributesByName = make(map[string]struct{})
	for _, v := range lt.Attributes {
		lt.attributesByName[v] = struct{}{}
	}
}

func (lt *LayersOfType) addLayer(layer Layer) (*Layer, error) {
	if _, has := lt.layersByKey[layer.key()]; has {
		return nil, common.ErrAlreadyExists
	}
	lt.Layers = append(lt.Layers, layer)
	lt.layersByKey[layer.key()] = &lt.Layers[len(lt.Layers)-1]
	return &lt.Layers[len(lt.Layers)-1], nil
}

func (lt *LayersOfType) addAttribute(attrName string) error {
	if _, has := lt.attributesByName[attrName]; has {
		return common.ErrAlreadyExists
	}
	lt.Attributes = append(lt.Attributes, attrName)
	lt.attributesByName[attrName] = struct{}{}
	return nil
}

func (lt *LayersOfType) Get(key attributeKey) (*Layer, bool) {
	v, has := lt.layersByKey[key]
	if !has {
		return nil, false
	}
	return v, true
}

// Layer represents a kitty layer.
// Field "Parts" represents a slice of image hashes in pairs, in which;
// 		1. each slice element represents a "part", and
//		2. within each part, the hash pair consists of;
//			[0] representing the layer "area".
//			[1] representing the layer "outline".
type Layer struct {
	OfAttribute string
	OfBreed     string
	Parts       [][2]cipher.SHA256
}

func (a *Layer) ensurePartsCount(n int) {
	if len(a.Parts) < n {
		a.Parts = append(a.Parts,
			make([][2]cipher.SHA256, n-len(a.Parts))...)
	}
}

func (a *Layer) key() attributeKey {
	return newAttributeKey(a.OfAttribute, a.OfBreed)
}

type attributeKey string

func newAttributeKey(attribute, breed string) attributeKey {
	return attributeKey(attribute + "_" + breed)
}
