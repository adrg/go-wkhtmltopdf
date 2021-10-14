package pdf

import (
	"sync"
	"unsafe"
)

type objectID = unsafe.Pointer

type objectRegistry struct {
	sync.RWMutex

	objects map[objectID]interface{}
}

func newObjectRegistry() *objectRegistry {
	return &objectRegistry{
		objects: map[objectID]interface{}{},
	}
}

func (or *objectRegistry) get(id objectID) (interface{}, bool) {
	if id == nil {
		return nil, false
	}

	or.RLock()
	object, ok := or.objects[id]
	or.RUnlock()

	return object, ok
}

func (or *objectRegistry) add(id objectID, object interface{}) {
	if id == nil {
		return
	}

	or.Lock()
	or.objects[id] = object
	or.Unlock()
}

func (or *objectRegistry) remove(id objectID) {
	if id == nil {
		return
	}

	or.Lock()
	delete(or.objects, id)
	or.Unlock()
}
