package db

import (
	"context"
	"fmt"
	"sync"
)

// ProviderFunc creates a database connection and returns it as interface{}.
type ProviderFunc func(ctx context.Context, cfg map[string]string) (interface{}, error)

// CloseFunc closes a previously opened connection.
type CloseFunc func(ctx context.Context, conn interface{}) error

// ProviderRegistration bundles the open and close functions for a database kind.
type ProviderRegistration struct {
	Open  ProviderFunc
	Close CloseFunc
}

type instance struct {
	kind string
	conn interface{}
}

// ConnectionManager manages named database connections.
type ConnectionManager struct {
	mu        sync.RWMutex
	providers map[string]ProviderRegistration
	instances map[string]instance
}

// NewConnectionManager creates a new ConnectionManager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		providers: make(map[string]ProviderRegistration),
		instances: make(map[string]instance),
	}
}

// RegisterProvider registers an open/close pair for a database kind.
func (m *ConnectionManager) RegisterProvider(kind string, reg ProviderRegistration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers[kind] = reg
}

// Connect opens a named connection using the appropriate provider.
func (m *ConnectionManager) Connect(ctx context.Context, name, kind string, cfg map[string]string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	provider, ok := m.providers[kind]
	if !ok {
		return fmt.Errorf("no provider registered for kind %q", kind)
	}

	conn, err := provider.Open(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect %q: %w", name, err)
	}

	m.instances[name] = instance{kind: kind, conn: conn}
	return nil
}

// Get retrieves a named connection.
func (m *ConnectionManager) Get(name string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inst, ok := m.instances[name]
	if !ok {
		return nil, fmt.Errorf("connection %q not found", name)
	}
	return inst.conn, nil
}

// CloseAll terminates all open connections.
func (m *ConnectionManager) CloseAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, inst := range m.instances {
		provider, ok := m.providers[inst.kind]
		if !ok {
			continue
		}
		if err := provider.Close(ctx, inst.conn); err != nil {
			errs = append(errs, fmt.Errorf("failed to close %q: %w", name, err))
		}
	}
	m.instances = make(map[string]instance)

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}
	return nil
}
