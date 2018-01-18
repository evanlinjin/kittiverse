package incubator

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

const (
	ImageMin = 0
	ImageMax = 1200
)

var (
	ErrReturn    = errors.New("returned action, no error")
	ErrAccessory = errors.New("is accessory")
	ErrDNAPart   = errors.New("is dna part")

	rootDir    = "/home/evan/skycoin/ivan/kittycash/Kitties"
	rootDirMux sync.RWMutex
)

func GetRootDir() string {
	rootDirMux.RLock()
	defer rootDirMux.RUnlock()

	return rootDir
}

func SetRootDir(path string) error {
	rootDirMux.Lock()
	defer rootDirMux.Unlock()

	var e error
	if path, e = filepath.Abs(path); e != nil {
		return e
	} else if _, e = os.Stat(path); e != nil {
		return e
	} else {
		return nil
	}
}

func init() {
	for _, part := range kittyParts {
		kittyPartsByName.Store(part.folderName, part)
	}
}
