package river

import (
	"github.com/google/blueprint"
)

type LibraryProducer interface {
	LibraryFileName() string
}

func IsLibraryProducer(module blueprint.Module) bool {
	_, ok := module.(LibraryProducer)
	return ok
}
