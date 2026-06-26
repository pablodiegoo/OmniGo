package channel

import "sync"

// Registry provides lookup of Dispatcher implementations by channel name.
// It is safe for concurrent use.
type Registry struct {
	mu          sync.RWMutex
	dispatchers map[string]Dispatcher
}

// NewRegistry creates a registry pre-loaded with the given dispatchers.
// The map is copied — later mutations to the caller's map don't affect
// the registry.
func NewRegistry(dispatchers map[string]Dispatcher) *Registry {
	r := &Registry{
		dispatchers: make(map[string]Dispatcher, len(dispatchers)),
	}
	r.mu.Lock()
	for k, v := range dispatchers {
		r.dispatchers[k] = v
	}
	r.mu.Unlock()
	return r
}

// Get returns the dispatcher for the named channel, or nil if not registered.
func (r *Registry) Get(name string) (Dispatcher, bool) {
	r.mu.RLock()
	d, ok := r.dispatchers[name]
	r.mu.RUnlock()
	return d, ok
}

// GetOrDefault returns the dispatcher for the named channel, or fallback
// if not registered.
func (r *Registry) GetOrDefault(name string, fallback Dispatcher) Dispatcher {
	r.mu.RLock()
	d, ok := r.dispatchers[name]
	r.mu.RUnlock()
	if !ok {
		return fallback
	}
	return d
}

// Register adds or replaces a dispatcher for the given channel name.
func (r *Registry) Register(name string, d Dispatcher) {
	r.mu.Lock()
	r.dispatchers[name] = d
	r.mu.Unlock()
}

// Len returns the number of registered dispatchers.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.dispatchers)
}

// Names returns the names of all registered dispatchers.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.dispatchers))
	for k := range r.dispatchers {
		names = append(names, k)
	}
	return names
}
