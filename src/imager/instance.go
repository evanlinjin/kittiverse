package imager

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Instance struct {
	LayerTypes []LayerType

	layerTypesByName sync.Map

	dir string               // root directory.
	log *logrus.Logger       // logging.
	ic  *ImageContainer      // contains all images.
	lc  *LayerTypesContainer // contains layers.
}
