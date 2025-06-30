package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FileStorage implements file-based storage for personas and identities
type FileStorage struct {
	dataDir        string
	personasDir    string
	identitiesDir  string
	communitiesDir string
	mu             sync.RWMutex
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(dataDir string) (*FileStorage, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	personasDir := filepath.Join(dataDir, "personas")
	if err := os.MkdirAll(personasDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create personas directory: %v", err)
	}

	identitiesDir := filepath.Join(dataDir, "identities")
	if err := os.MkdirAll(identitiesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create identities directory: %v", err)
	}

	communitiesDir := filepath.Join(dataDir, "communities")
	if err := os.MkdirAll(communitiesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create communities directory: %v", err)
	}

	return &FileStorage{
		dataDir:        dataDir,
		personasDir:    personasDir,
		identitiesDir:  identitiesDir,
		communitiesDir: communitiesDir,
	}, nil
}

// Persona operations
func (f *FileStorage) Create(p *types.Persona) error {
	f.mu.Lock()
	defer f.mu.Unlock()

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

	p.Id = f.generateID()
	return f.writePersona(*p)
}

func (f *FileStorage) Get(id string) (types.Persona, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.readPersona(id)
}

func (f *FileStorage) List() ([]types.Persona, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	files, err := os.ReadDir(f.personasDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read personas directory: %v", err)
	}

	var personas []types.Persona
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			id := file.Name()[:len(file.Name())-5] // Remove .json extension
			if p, err := f.readPersona(id); err == nil {
				personas = append(personas, p)
			}
		}
	}

	return personas, nil
}

func (f *FileStorage) Update(id string, p types.Persona) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Check if persona exists
	if _, err := f.readPersona(id); err != nil {
		return fmt.Errorf("persona not found: %s", id)
	}

	p.Id = id
	return f.writePersona(p)
}

func (f *FileStorage) Delete(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := filepath.Join(f.personasDir, id+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("persona not found: %s", id)
	}

	return os.Remove(filePath)
}

// Identity operations
func (f *FileStorage) CreateIdentity(i *types.Identity) error {
	f.mu.Lock()
	defer f.mu.Unlock()

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
	if _, err := f.readPersona(i.PersonaId); err != nil {
		return fmt.Errorf("referenced persona not found: %s", i.PersonaId)
	}

	i.Id = f.generateID()
	now := time.Now()
	i.CreatedAt = timestamppb.New(now)
	i.UpdatedAt = timestamppb.New(now)

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

	return f.writeIdentity(*i)
}

func (f *FileStorage) GetIdentity(id string) (types.Identity, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.readIdentity(id)
}

func (f *FileStorage) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	files, err := os.ReadDir(f.identitiesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read identities directory: %v", err)
	}

	var identities []types.Identity
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			id := file.Name()[:len(file.Name())-5] // Remove .json extension
			if i, err := f.readIdentity(id); err == nil {
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
				identities = append(identities, i)
			}
		}
	}

	return identities, nil
}

func (f *FileStorage) UpdateIdentity(id string, i types.Identity) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Check if identity exists
	if _, err := f.readIdentity(id); err != nil {
		return fmt.Errorf("identity not found: %s", id)
	}

	// Verify persona exists
	if _, err := f.readPersona(i.PersonaId); err != nil {
		return fmt.Errorf("referenced persona not found: %s", i.PersonaId)
	}

	i.Id = id
	i.UpdatedAt = timestamppb.New(time.Now())
	return f.writeIdentity(i)
}

func (f *FileStorage) DeleteIdentity(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := filepath.Join(f.identitiesDir, id+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("identity not found: %s", id)
	}

	return os.Remove(filePath)
}

func (f *FileStorage) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	i, err := f.readIdentity(id)
	if err != nil {
		return types.IdentityWithPersona{}, err
	}

	p, err := f.readPersona(i.PersonaId)
	if err != nil {
		return types.IdentityWithPersona{}, fmt.Errorf("referenced persona not found: %s", i.PersonaId)
	}

	return types.IdentityWithPersona{
		Identity: i,
		Persona:  p,
	}, nil
}

// Helper methods
func (f *FileStorage) generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (f *FileStorage) readPersona(id string) (types.Persona, error) {
	filePath := filepath.Join(f.personasDir, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return types.Persona{}, fmt.Errorf("persona not found: %s", id)
		}
		return types.Persona{}, fmt.Errorf("failed to read persona file: %v", err)
	}

	var p types.Persona
	if err := json.Unmarshal(data, &p); err != nil {
		return types.Persona{}, fmt.Errorf("failed to parse persona data: %v", err)
	}

	return p, nil
}

func (f *FileStorage) writePersona(p types.Persona) error {
	filePath := filepath.Join(f.personasDir, p.Id+".json")
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal persona data: %v", err)
	}

	return os.WriteFile(filePath, data, 0644)
}

// Community operations
func (f *FileStorage) CreateCommunity(c *types.Community) error {
	f.mu.Lock()
	defer f.mu.Unlock()

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
		c.Id = f.generateID()
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

	return f.writeCommunity(*c)
}

func (f *FileStorage) GetCommunity(id string) (types.Community, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.readCommunity(id)
}

func (f *FileStorage) ListCommunities(filter *types.CommunityFilter) ([]types.Community, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	files, err := os.ReadDir(f.communitiesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read communities directory: %v", err)
	}

	var communities []types.Community
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			id := file.Name()[:len(file.Name())-5] // Remove .json extension
			if c, err := f.readCommunity(id); err == nil {
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
				communities = append(communities, c)
			}
		}
	}

	return communities, nil
}

func (f *FileStorage) UpdateCommunity(id string, c types.Community) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Check if community exists
	if _, err := f.readCommunity(id); err != nil {
		return fmt.Errorf("community not found: %s", id)
	}

	c.Id = id
	return f.writeCommunity(c)
}

func (f *FileStorage) DeleteCommunity(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := filepath.Join(f.communitiesDir, id+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("community not found: %s", id)
	}

	return os.Remove(filePath)
}

func (f *FileStorage) readCommunity(id string) (types.Community, error) {
	filePath := filepath.Join(f.communitiesDir, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return types.Community{}, fmt.Errorf("community not found: %s", id)
		}
		return types.Community{}, fmt.Errorf("failed to read community file: %v", err)
	}

	var c types.Community
	if err := json.Unmarshal(data, &c); err != nil {
		return types.Community{}, fmt.Errorf("failed to parse community data: %v", err)
	}

	return c, nil
}

func (f *FileStorage) writeCommunity(c types.Community) error {
	filePath := filepath.Join(f.communitiesDir, c.Id+".json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal community data: %v", err)
	}

	return os.WriteFile(filePath, data, 0644)
}

func (f *FileStorage) readIdentity(id string) (types.Identity, error) {
	filePath := filepath.Join(f.identitiesDir, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return types.Identity{}, fmt.Errorf("identity not found: %s", id)
		}
		return types.Identity{}, fmt.Errorf("failed to read identity file: %v", err)
	}

	var i types.Identity
	if err := json.Unmarshal(data, &i); err != nil {
		return types.Identity{}, fmt.Errorf("failed to parse identity data: %v", err)
	}

	return i, nil
}

func (f *FileStorage) writeIdentity(i types.Identity) error {
	filePath := filepath.Join(f.identitiesDir, i.Id+".json")
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal identity data: %v", err)
	}

	return os.WriteFile(filePath, data, 0644)
}
