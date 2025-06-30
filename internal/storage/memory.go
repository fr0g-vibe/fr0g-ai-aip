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
	personas    map[string]types.Persona
	identities  map[string]types.Identity
	communities map[string]types.Community
	mu          sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		personas:    make(map[string]types.Persona),
		identities:  make(map[string]types.Identity),
		communities: make(map[string]types.Community),
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

	p.Id = generateID()
	m.personas[p.Id] = *p
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

	p.Id = id
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
	if i.PersonaId == "" {
		return fmt.Errorf("persona ID is required")
	}
	if i.Name == "" {
		return fmt.Errorf("identity name is required")
	}

	// Verify persona exists
	if _, exists := m.personas[i.PersonaId]; !exists {
		return fmt.Errorf("referenced persona not found: %s", i.PersonaId)
	}

	i.Id = generateID()
	now := time.Now()
	i.CreatedAt = now
	i.UpdatedAt = now

	// Set default values
	if i.RichAttributes == nil {
		i.RichAttributes = &types.RichAttributes{}
	}
	if i.Tags == nil {
		i.Tags = []string{}
	}
	if !i.IsActive {
		i.IsActive = true
	}

	m.identities[i.Id] = *i
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
			if filter.PersonaID != "" && i.PersonaId != filter.PersonaID {
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
	if _, exists := m.personas[i.PersonaId]; !exists {
		return fmt.Errorf("referenced persona not found: %s", i.PersonaId)
	}

	i.Id = id
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

	p, exists := m.personas[i.PersonaId]
	if !exists {
		return types.IdentityWithPersona{}, fmt.Errorf("referenced persona not found: %s", i.PersonaId)
	}

	return types.IdentityWithPersona{
		Identity: i,
		Persona:  p,
	}, nil
}

// Community operations
func (m *MemoryStorage) CreateCommunity(c *types.Community) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c == nil {
		return fmt.Errorf("community cannot be nil")
	}
	if c.Name == "" {
		return fmt.Errorf("community name is required")
	}
	if c.Type == "" {
		return fmt.Errorf("community type is required")
	}

	if c.Id == "" {
		c.Id = generateID()
	}
	
	// Initialize empty slices if nil
	if c.MemberIds == nil {
		c.MemberIds = []string{}
	}
	if c.Tags == nil {
		c.Tags = []string{}
	}
	if c.Attributes == nil {
		c.Attributes = make(map[string]interface{})
	}

	m.communities[c.Id] = *c
	return nil
}

func (m *MemoryStorage) GetCommunity(id string) (types.Community, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, exists := m.communities[id]
	if !exists {
		return types.Community{}, fmt.Errorf("community not found: %s", id)
	}
	return c, nil
}

func (m *MemoryStorage) ListCommunities(filter *types.CommunityFilter) ([]types.Community, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []types.Community

	for _, c := range m.communities {
		// Apply filters
		if filter != nil {
			if filter.Type != "" && c.Type != filter.Type {
				continue
			}
			if filter.IsActive != nil && c.IsActive != *filter.IsActive {
				continue
			}
			if filter.MinSize != nil && c.Size < *filter.MinSize {
				continue
			}
			if filter.MaxSize != nil && c.Size > *filter.MaxSize {
				continue
			}
			if filter.MinDiversity != nil && c.Diversity < *filter.MinDiversity {
				continue
			}
			if filter.MaxDiversity != nil && c.Diversity > *filter.MaxDiversity {
				continue
			}
			if len(filter.Tags) > 0 {
				hasTag := false
				for _, tag := range filter.Tags {
					for _, communityTag := range c.Tags {
						if communityTag == tag {
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
				nameMatch := strings.Contains(strings.ToLower(c.Name), searchLower)
				descMatch := strings.Contains(strings.ToLower(c.Description), searchLower)
				if !nameMatch && !descMatch {
					continue
				}
			}
		}
		result = append(result, c)
	}

	return result, nil
}

func (m *MemoryStorage) UpdateCommunity(id string, c types.Community) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.communities[id]; !exists {
		return fmt.Errorf("community not found: %s", id)
	}

	c.Id = id
	m.communities[id] = c
	return nil
}

func (m *MemoryStorage) DeleteCommunity(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.communities[id]; !exists {
		return fmt.Errorf("community not found: %s", id)
	}
	delete(m.communities, id)
	return nil
}
