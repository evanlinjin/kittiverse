package container

type Layers interface {
	Version() uint16
	Import(raw []byte) error
	Export() []byte
	Compile(rootDir string, images Images) error
}
