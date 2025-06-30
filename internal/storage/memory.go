package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// MemoryStorage implements in-memory storage for personas and identities
type MemoryStorage struct {
	personas   map[string]types.Persona
	identities map[string]types.Identity
	mu         sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		personas:   make(map[string]types.Persona),
		identities: make(map[string]types.Identity),
	}
}

// generateID creates a random ID for a persona or identity
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Persona operations
func (m *MemoryStorage) Create(p *types.Persona) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p == nil {
		return fmt.Errorf("persona cannot be nil")
	}
	if p.Name == "" {
		return fmt.Errorf("persona name is required")
	}
	if p.Topic == "" {
		return fmt.Errorf("persona topic is required")
	}
	if p.Prompt == "" {
		return fmt.Errorf("persona prompt is required")
	}

	p.ID = generateID()
	m.personas[p.ID] = *p
	return nil
}

func (m *MemoryStorage) Get(id string) (types.Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.personas[id]
	if !exists {
		return types.Persona{}, fmt.Errorf("persona not found: %s", id)
	}
	return p, nil
}

func (m *MemoryStorage) List() ([]types.Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]types.Persona, 0, len(m.personas))
	for _, p := range m.personas {
		result = append(result, p)
	}
	return result, nil
}

func (m *MemoryStorage) Update(id string, p types.Persona) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.personas[id]; !exists {
		return fmt.Errorf("persona not found: %s", id)
	}

	p.ID = id
	m.personas[id] = p
	return nil
}

func (m *MemoryStorage) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.personas[id]; !exists {
		return fmt.Errorf("persona not found: %s", id)
	}
	delete(m.personas, id)
	return nil
}

// Identity operations
func (m *MemoryStorage) CreateIdentity(i *types.Identity) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if i == nil {
		return fmt.Errorf("identity cannot be nil")
	}
	if i.PersonaID == "" {
		return fmt.Errorf("persona ID is required")
	}
	if i.Name == "" {
		return fmt.Errorf("identity name is required")
	}

	// Verify persona exists
	if _, exists := m.personas[i.PersonaID]; !exists {
		return fmt.Errorf("referenced persona not found: %s", i.PersonaID)
	}

	i.ID = generateID()
	now := time.Now()
	i.CreatedAt = now
	i.UpdatedAt = now

	// Set default values
	if i.Attributes == nil {
		i.Attributes = make(map[string]string)
	}
	if i.Preferences == nil {
		i.Preferences = make(map[string]string)
	}
	if i.Tags == nil {
		i.Tags = []string{}
	}
	if !i.IsActive {
		i.IsActive = true
	}

	m.identities[i.ID] = *i
	return nil
}

func (m *MemoryStorage) GetIdentity(id string) (types.Identity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	i, exists := m.identities[id]
	if !exists {
		return types.Identity{}, fmt.Errorf("identity not found: %s", id)
	}
	return i, nil
}

func (m *MemoryStorage) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []types.Identity

	for _, i := range m.identities {
		// Apply filters
		if filter != nil {
			if filter.PersonaID != "" && i.PersonaID != filter.PersonaID {
				continue
			}
			if filter.IsActive != nil && i.IsActive != *filter.IsActive {
				continue
			}
			if len(filter.Tags) > 0 {
				hasTag := false
				for _, tag := range filter.Tags {
					for _, identityTag := range i.Tags {
						if identityTag == tag {
							hasTag = true
							break
						}
					}
					if hasTag {
						break
					}
				}
				if !hasTag {
					continue
				}
			}
			if filter.Search != "" {
				searchLower := strings.ToLower(filter.Search)
				nameMatch := strings.Contains(strings.ToLower(i.Name), searchLower)
				descMatch := strings.Contains(strings.ToLower(i.Description), searchLower)
				if !nameMatch && !descMatch {
					continue
				}
			}
		}
		result = append(result, i)
	}

	return result, nil
}

func (m *MemoryStorage) UpdateIdentity(id string, i types.Identity) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.identities[id]; !exists {
		return fmt.Errorf("identity not found: %s", id)
	}

	// Verify persona exists
	if _, exists := m.personas[i.PersonaID]; !exists {
		return fmt.Errorf("referenced persona not found: %s", i.PersonaID)
	}

	i.ID = id
	i.UpdatedAt = time.Now()
	m.identities[id] = i
	return nil
}

func (m *MemoryStorage) DeleteIdentity(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.identities[id]; !exists {
		return fmt.Errorf("identity not found: %s", id)
	}
	delete(m.identities, id)
	return nil
}

func (m *MemoryStorage) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	i, exists := m.identities[id]
	if !exists {
		return types.IdentityWithPersona{}, fmt.Errorf("identity not found: %s", id)
	}

	p, exists := m.personas[i.PersonaID]
	if !exists {
		return types.IdentityWithPersona{}, fmt.Errorf("referenced persona not found: %s", i.PersonaID)
	}

	return types.IdentityWithPersona{
		Identity: i,
		Persona:  p,
	}, nil
}
