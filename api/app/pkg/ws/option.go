package ws

type OptionFunc func(*Manager)

func (m *Manager) With(opt ...OptionFunc) *Manager {
	for _, o := range opt {
		o(m)
	}
	return m
}

func WithProcessorFabric(fabric PipeProcessorFabric) OptionFunc {
	return func(m *Manager) {
		m.processorFabric = fabric
	}
}
