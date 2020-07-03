package uniquer

type MemoryUniquer struct {
	state map[string]struct{}
}

func New() *MemoryUniquer {
	return &MemoryUniquer{state: make(map[string]struct{})}
}

func (m *MemoryUniquer) Exists(s string) (bool, error) {
	_, ok := m.state[s]
	return ok, nil
}

func (m *MemoryUniquer) Add(s string) error {
	m.state[s] = struct{}{}
	return nil
}
