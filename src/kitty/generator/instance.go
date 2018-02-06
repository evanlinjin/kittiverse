package generator

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Instance struct {
	LayerTypes []LayersOfType

	layerTypesByName sync.Map

	dir string           // root directory.
	log *logrus.Logger   // logging.
	ic  *ImageContainer  // contains all images.
	lc  *LayersContainer // contains layers.
}
