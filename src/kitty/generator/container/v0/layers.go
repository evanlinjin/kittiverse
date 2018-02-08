package v0

import (
	"github.com/kittycash/kittiverse/src/kitty/generator/container"
	"github.com/kittycash/kittiverse/src/kitty/generator/container/common"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"io/ioutil"
	"path"
	"strings"
	"os"
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

func (lc *Layer) GenerateKitty() {

}

/*
	<<< MEMBER HELPERS >>>
*/

func (lc *Layers) addLayerType(ltName string) error {
	if _, has := lc.layerTypesByName[ltName]; has {
		return common.ErrAlreadyExists
	}
	lc.LayerTypes = append(lc.LayerTypes, LayersOfType{
		OfType:      ltName,
		layersByKey: make(map[attributeKey]*Layer),
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
					fullPath      = path.Join(rootDir, lt.OfType, breed, file.Name())
					fullName      = strings.TrimSuffix(file.Name(), ".png")
					splitName     = strings.Split(fullName, "_")

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
					 layer, e = lt.Add(Layer{
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
				layer.ensurePartsCount(partIndex+1)
				switch {
				case isArea:
					layer.Parts[partIndex][0] = imgHash
				case isOutline:
					layer.Parts[partIndex][1] = imgHash
				}
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
	OfType      string
	Layers      []Layer
	layersByKey map[attributeKey]*Layer `enc:"-"` // aka, by attribute and breed
}

func (lt *LayersOfType) Init() {
	if lt.layersByKey == nil {
		lt.layersByKey = make(map[attributeKey]*Layer)
	}
	for i := 0; i < len(lt.Layers); i++ {
		layer := &lt.Layers[i]
		lt.layersByKey[layer.key()] = layer
	}
}

func (lt *LayersOfType) Add(layer Layer) (*Layer, error) {
	if _, has := lt.layersByKey[layer.key()]; has {
		return nil, common.ErrAlreadyExists
	}
	lt.Layers = append(lt.Layers, layer)
	lt.layersByKey[layer.key()] = &lt.Layers[len(lt.Layers)-1]
	return &lt.Layers[len(lt.Layers)-1], nil
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
			make([][2]cipher.SHA256, n - len(a.Parts))...)
	}
}

func (a *Layer) key() attributeKey {
	return newAttributeKey(a.OfAttribute, a.OfBreed)
}

type attributeKey string

func newAttributeKey(attribute, breed string) attributeKey {
	return attributeKey(attribute + "_" + breed)
}
