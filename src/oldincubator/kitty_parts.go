package oldincubator

type KittyPart int

func (kp KittyPart) Specs() *KittyPartSpecs {
	return kittyParts[kp]
}

const (
	Tail KittyPart = iota
	Body
	Ears
	Head
	Eyes
	Brows
	Nose
	Cap
	Collar
	kittyPartsCount
)

type KittyPartSpecs struct {
	folderName string // Folder name of the part.
	fieldName  string // KittyConfig field folderName of the part.
	accessory  bool   // Is accessory.
}

func (kps *KittyPartSpecs) FolderName() string { return kps.folderName }
func (kps *KittyPartSpecs) FieldName() string  { return kps.fieldName }
func (kps *KittyPartSpecs) IsAccessory() bool  { return kps.accessory }

var (
	kittyParts = [...]*KittyPartSpecs{
		Body:   {"body", "Body", false},
		Brows:  {"brows", "Brows", false},
		Cap:    {"cap", "Cap", true},
		Collar: {"collar", "Collar", true},
		Ears:   {"ears", "Ears", false},
		Eyes:   {"eyes", "Eyes", false},
		Head:   {"head", "Head", false},
		Nose:   {"nose", "Nose", false},
		Tail:   {"tail", "Tail", false},
	}
)

type KittyPartAction func(part KittyPart) error

func RangeKittyParts(action KittyPartAction) error {
	for id := KittyPart(0); id < kittyPartsCount; id++ {
		if e := action(id); e != nil {
			switch e {
			case ErrReturn:
				return nil
			default:
				return e
			}
		}
	}
	return nil
}
