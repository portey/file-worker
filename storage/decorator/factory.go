package decorator

import "github.com/portey/file-worker/storage/decorator/uniquer"

type MemoryFactory struct {
}

func NewMemoryUniquerFactory() UniquerFactory {
	return &MemoryFactory{}
}

func (f *MemoryFactory) NewMatcher() UniqueMatcher {
	return uniquer.New()
}
