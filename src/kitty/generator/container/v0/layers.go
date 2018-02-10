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
	"errors"
	"image/draw"
)

const (
	PrefixAccessory = "accessory"
)

type Layers struct {
	LayerTypes       []LayersOfType
	Breeds           []string
	layerTypesByName map[string]int `enc:"-"`
	breedsByName     map[string]int `enc:"-"`
}

func NewLayersContainer() *Layers {
	return &Layers{
		layerTypesByName: make(map[string]int),
		breedsByName:     make(map[string]int),
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
	lc.layerTypesByName = make(map[string]int)
	for i, v := range lc.LayerTypes {
		lc.layerTypesByName[v.OfType] = i
	}
	lc.breedsByName = make(map[string]int)
	for i, v := range lc.Breeds {
		lc.breedsByName[v] = i
	}
	return nil
}

func (lc *Layers) Export() []byte {
	return append(encoder.Serialize(version), encoder.Serialize(lc)...)
}

func (lc *Layers) Compile(rootDir string, images container.Images) error {
	// Get layer types.
	if e := initLayerTypes(lc, rootDir); e != nil {
		log.WithError(e).Error("failed to initiate later types")
		return e
	}
	// Get layers.
	if e := initLayers(lc, rootDir, images); e != nil {
		log.WithError(e).Error("failed to initiate layers")
		return e
	}
	return nil
}

func (lc *Layers) GetAlleleRanges() *genetics.AlleleRanges {
	getRange := func(pos genetics.DNAPos) genetics.AlleleRange {
		return genetics.AlleleRange{
			Min: genetics.Allele{}.String(),
			Max: genetics.NewAlleleFromUint16(uint16(
				len(lc.LayerTypes[lc.layerTypesByName[pos.String()]].Attributes) - 1,
			)).String(),
		}
	}
	return &genetics.AlleleRanges{
		Breed: genetics.AlleleRange{
			Min: genetics.Allele{}.String(),
			Max: genetics.NewAlleleFromUint16(uint16(
				len(lc.Breeds) - 1,
			)).String(),
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

func (lc *Layers) GenerateKitty(ic container.Images, dna genetics.DNA) (image.Image, error) {
	out := image.NewRGBA(image.Rect(0, 0, common.XpxLen, common.YpxLen))

	// Get breed.
	breed := lc.getBreed(dna.GetPhenotype(genetics.DNABreedPos))

	// Make image input common.
	iic := &imgInputCommon{lc: lc, ic: ic, breed: breed, dna: dna}

	// Generate fur.
	fur := common.EmptyImage()
	{
		bg, e := generateImage(iic, genetics.DNABodyColorAPos, nil)
		if e != nil {
			return nil, e
		}
		fg, e := generateImage(iic, genetics.DNABodyColorBPos, nil)
		if e != nil {
			return nil, e
		}
		pt, e := generateImage(iic, genetics.DNABodyPatternPos, fg)
		if e != nil {
			return nil, e
		}
		common.DrawOutline(out, bg)
		common.DrawOutline(out, pt)
	}

	// Generate ears.
	{
		ears, e := generateImage(iic, genetics.DNAEarsAttrPos, fur)
		if e != nil {
			return nil, e
		}
		common.DrawOutline(out, ears)
	}

	// Generate tail.
	{
		tail, e := generateImage(iic, genetics.DNATailAttrPos, fur)
		if e != nil {
			return nil, e
		}
		common.DrawOutline(out, tail)
	}

	// Generate body.
	{
		body, e := generateImage(iic, genetics.DNABodyAttrPos, fur)
		if e != nil {
			return nil, e
		}
		common.DrawOutline(out, body)
	}

	// Generate nose.
	{
		// TODO: Change to noseColor.
		nose, e := generateImage(iic, genetics.DNANoseAttrPos, fur)
		if e != nil {
			return nil, e
		}
		common.DrawOutline(out, nose)
	}

	// Generate eyes.
	{
		bg, e := generateImage(iic, genetics.DNAEyesColorPos, nil)
		if e != nil {
			return nil, e
		}
		fg, e := generateImage(iic, genetics.DNAEyesAttrPos, bg)
		if e != nil {
			return nil, e
		}
		common.DrawOutline(out, fg)
	}

	return out, nil
}

/*
	<<< MEMBER HELPERS >>>
*/

type layerTypeAction func(lt *LayersOfType, ltDir string, breedDirs []os.FileInfo) error

func (lc *Layers) rangeLayerTypes(rootDir string, action layerTypeAction) {
	for i := 0; i < len(lc.LayerTypes); i++ {
		var (
			lt    = &lc.LayerTypes[i]
			ltDir = path.Join(rootDir, lt.OfType)
		)
		breedDirs, e := ioutil.ReadDir(ltDir)
		if e != nil {
			log.WithField("index", i).
				WithField("layer_type", lt.OfType).
				WithField("dir", ltDir).
				WithError(e).Error("failed to read directory")
			continue
		}
		action(lt, ltDir, breedDirs)
	}
}

func (lc *Layers) addLayerType(ltName string) error {
	if _, has := lc.layerTypesByName[ltName]; has {
		return common.ErrAlreadyExists
	}
	lc.LayerTypes = append(lc.LayerTypes, LayersOfType{
		OfType:           ltName,
		layersByKey:      make(map[attributeKey]int),
		attributesByName: make(map[string]int),
	})
	lc.layerTypesByName[ltName] = len(lc.LayerTypes) - 1
	return nil
}

func (lc *Layers) addBreed(bName string) error {
	if _, has := lc.breedsByName[bName]; has {
		return common.ErrAlreadyExists
	}
	lc.Breeds = append(lc.Breeds, bName)
	lc.breedsByName[bName] = len(lc.Breeds) - 1
	return nil
}

func (lc *Layers) getLayerType(pos genetics.DNAPos) *LayersOfType {
	return &lc.LayerTypes[lc.layerTypesByName[pos.String()]]
}

func (lc *Layers) getBreed(a genetics.Allele) string {
	return lc.Breeds[a.Uint16()]
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

func initLayers(lc *Layers, rootDir string, images container.Images) error {

	lc.rangeLayerTypes(rootDir, func(lt *LayersOfType, ltDir string, bDirs []os.FileInfo) error {
		log.WithField("dir", ltDir).
			WithField("breed_count", len(bDirs)).
			Printf("ranging later type '%s'", lt.OfType)

		for _, bDir := range bDirs {
			// Skip if not directory.
			if bDir.IsDir() == false {
				continue
			}
			breed := bDir.Name()
			lc.addBreed(breed)
			// collect attributes.
			files, e := ioutil.ReadDir(path.Join(ltDir, bDir.Name()))
			if e != nil {
				return e
			}
			for _, file := range files {
				if file.IsDir() || strings.HasSuffix(file.Name(), ".png") == false {
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
					isArea, isOutline = false, true
				}
				// Extract image.
				imgRaw, e := common.GetRawImageFromFile(fullPath)
				if e != nil {
					log.WithError(e).Error("failed to read image file")
				}
				imgHash := images.GetOrAdd(imgRaw)
				// Append.
				layer := lt.getOrAddLayer(Layer{OfAttribute: attributeName, OfBreed: breed})
				layer.ensurePartsCount(partIndex + 1)
				switch {
				case isArea:
					layer.Parts[partIndex][0] = imgHash
				case isOutline:
					layer.Parts[partIndex][1] = imgHash
				}
				// Ensure attribute.
				if e := lt.addAttribute(attributeName); e != nil {
					switch e {
					case common.ErrAlreadyExists:
					default:
						log.WithField("layer_type", lt.OfType).
							WithField("breed", breed).
							WithField("full_attribute", fullName).
							WithError(e).Error("failed to add attribute")
					}
				} else {
					log.WithField("layer_type", lt.OfType).
						WithField("breed", breed).
						WithField("full_attribute", fullName).
						Infof("attribute '%s' added", attributeName)
				}
			}
		}
		return nil
	})
	return nil
}

func getPartIndex(str string) int {
	p := strings.TrimPrefix(str, "part")
	return int([]byte(p)[0] - 65)
}

type imgInputCommon struct {
	lc *Layers
	ic container.Images
	breed string
	dna genetics.DNA
}

func generateImage(c *imgInputCommon, dnaPos genetics.DNAPos, bg image.Image, ps ...int) (image.Image, error) {
	var (
		allele    = c.dna.GetPhenotype(dnaPos)
		lt        = c.lc.getLayerType(dnaPos)
		attribute = lt.Attributes[allele.Uint16()]
	)

	layer, ok := lt.get(newAttributeKey(attribute, c.breed))
	if !ok {
		layer, ok = lt.get(newAttributeKey(attribute, "default"))
		if !ok {
			log.WithField("layer_type", lt.OfType).
				WithField("breed", c.breed).
				WithField("attribute", attribute).
				Error("failed to find layer")
			return nil, errors.New("failed to find layer")
		}
	}
	return layer.generateImage(c.ic, bg, ps...)
}

func overlayImage(c *imgInputCommon, dst draw.Image, dnaPos genetics.DNAPos, bg image.Image, ps ...int) error {
	src, e := generateImage(c, dnaPos, bg, ps...)
	if e != nil {
		return e
	}
	common.DrawOutline(dst, src)
	return nil
}

/*
	<<< TYPES >>>
*/

type LayersOfType struct {
	OfType           string
	Layers           []Layer
	Attributes       []string
	layersByKey      map[attributeKey]int `enc:"-"` // aka, by attribute and breed
	attributesByName map[string]int       `enc:"-"`
}

func (lt *LayersOfType) Init() {
	lt.layersByKey = make(map[attributeKey]int)
	for i, v := range lt.Layers {
		lt.layersByKey[v.key()] = i
	}
	lt.attributesByName = make(map[string]int)
	for i, v := range lt.Attributes {
		lt.attributesByName[v] = i
	}
}

func (lt *LayersOfType) getOrAddLayer(layer Layer) *Layer {
	var key = layer.key()
	if i, has := lt.layersByKey[key]; has {
		return &lt.Layers[i]
	} else {
		lt.Layers = append(lt.Layers, layer)
		i = len(lt.Layers) - 1
		lt.layersByKey[key] = i
		return &lt.Layers[i]
	}
}

func (lt *LayersOfType) addAttribute(attrName string) error {
	if _, has := lt.attributesByName[attrName]; has {
		return common.ErrAlreadyExists
	}
	lt.Attributes = append(lt.Attributes, attrName)
	lt.attributesByName[attrName] = len(lt.Attributes) - 1
	return nil
}

func (lt *LayersOfType) get(key attributeKey) (*Layer, bool) {
	i, has := lt.layersByKey[key]
	if !has {
		return nil, false
	}
	return &lt.Layers[i], true
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

type layerPartAction func(i int, areaImg, outlineImg image.Image)

func (a *Layer) rangeParts(ic container.Images, action layerPartAction) error {
	var e error
	for i, pair := range a.Parts {
		var areaImg image.Image
		if pair[0] != (cipher.SHA256{}) {
			if areaImg, e = common.GetImage(ic, pair[0]); e != nil {
				return e
			}
		}
		var outlineImg image.Image
		if pair[1] != (cipher.SHA256{}) {
			if outlineImg, e = common.GetImage(ic, pair[1]); e != nil {
				return e
			}
		}
		action(i, areaImg, outlineImg)
	}
	return nil
}

func (a *Layer) generateImage(ic container.Images, bg image.Image, ps ...int) (image.Image, error) {
	// partsMap informs of which parts are to be included in the generated image.
	var partsMap = make(map[int]bool)
	if len(ps) == 0 {
		for i := 0; i < len(a.Parts); i++ {
			partsMap[i] = true
		}
	} else {
		for _, i := range ps {
			partsMap[i] = true
		}
	}

	out := image.NewRGBA(image.Rect(0, 0, common.XpxLen, common.YpxLen))
	e := a.rangeParts(ic, func(i int, areaImg, outlineImg image.Image) {
		if partsMap[i] {
			if areaImg != nil {
				common.DrawArea(out, bg, areaImg)
			}
			if outlineImg != nil {
				common.DrawOutline(out, outlineImg)
			}
		}
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}