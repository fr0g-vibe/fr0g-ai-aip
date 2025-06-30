package types

import "time"

// Community represents a generated community of identities
type Community struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "geographic", "demographic", "interest", "political", "professional"
	
	// Community characteristics
	Size         int               `json:"size"`
	Diversity    float64          `json:"diversity"`    // 0.0-1.0, measure of internal diversity
	Cohesion     float64          `json:"cohesion"`     // 0.0-1.0, measure of group cohesion
	Attributes   map[string]interface{} `json:"attributes"` // Community-specific attributes
	
	// Member management
	MemberIds    []string         `json:"member_ids"`
	MaxMembers   int              `json:"max_members"`
	MinMembers   int              `json:"min_members"`
	
	// Generation parameters
	GenerationConfig CommunityGenerationConfig `json:"generation_config"`
	
	// Metadata
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Tags         []string         `json:"tags,omitempty"`
	IsActive     bool             `json:"is_active"`
}

// CommunityGenerationConfig defines parameters for community generation
type CommunityGenerationConfig struct {
	// Base persona distribution
	PersonaWeights map[string]float64 `json:"persona_weights"` // persona_id -> weight
	
	// Demographic constraints
	AgeDistribution    AgeDistribution    `json:"age_distribution"`
	LocationConstraint LocationConstraint `json:"location_constraint"`
	
	// Diversity settings
	PoliticalSpread    float64 `json:"political_spread"`    // 0.0-1.0, how politically diverse
	InterestSpread     float64 `json:"interest_spread"`     // 0.0-1.0, how diverse interests are
	SocioeconomicRange float64 `json:"socioeconomic_range"` // 0.0-1.0, income/class diversity
	
	// Relationship patterns
	NetworkDensity     float64 `json:"network_density"`     // 0.0-1.0, how interconnected
	ClusteringFactor   float64 `json:"clustering_factor"`   // 0.0-1.0, tendency to form subgroups
	
	// Behavioral parameters
	ActivityLevel      float64 `json:"activity_level"`      // 0.0-1.0, how active members are
	EngagementStyle    string  `json:"engagement_style"`    // "collaborative", "competitive", "passive"
}

// AgeDistribution defines age distribution parameters
type AgeDistribution struct {
	Mean     float64 `json:"mean"`
	StdDev   float64 `json:"std_dev"`
	MinAge   int     `json:"min_age"`
	MaxAge   int     `json:"max_age"`
	Skewness float64 `json:"skewness"` // -1.0 to 1.0, negative = younger skew
}

// LocationConstraint defines geographic constraints
type LocationConstraint struct {
	Type        string   `json:"type"`        // "city", "region", "country", "global"
	Locations   []string `json:"locations"`   // specific locations to include
	Radius      float64  `json:"radius"`      // km radius for geographic clustering
	Urban       *bool    `json:"urban"`       // true=urban, false=rural, nil=mixed
	Timezone    string   `json:"timezone"`    // preferred timezone
}

// CommunityFilter defines filtering options for community queries
type CommunityFilter struct {
	Type         string   `json:"type,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
	MinSize      *int     `json:"min_size,omitempty"`
	MaxSize      *int     `json:"max_size,omitempty"`
	MinDiversity *float64 `json:"min_diversity,omitempty"`
	MaxDiversity *float64 `json:"max_diversity,omitempty"`
	Search       string   `json:"search,omitempty"`
}

// CommunityMember represents a member within a community context
type CommunityMember struct {
	Identity     Identity               `json:"identity"`
	Role         string                 `json:"role"`         // "leader", "active", "passive", "newcomer"
	Influence    float64               `json:"influence"`    // 0.0-1.0, influence within community
	Connections  []string              `json:"connections"`  // IDs of other members they're connected to
	JoinedAt     time.Time             `json:"joined_at"`
	LastActive   time.Time             `json:"last_active"`
	Attributes   map[string]interface{} `json:"attributes"`   // Member-specific community attributes
}

// CommunityStats provides analytics about a community
type CommunityStats struct {
	CommunityId      string            `json:"community_id"`
	MemberCount      int               `json:"member_count"`
	ActiveMembers    int               `json:"active_members"`
	AverageAge       float64           `json:"average_age"`
	GenderRatio      map[string]float64 `json:"gender_ratio"`
	LocationSpread   map[string]int    `json:"location_spread"`
	PoliticalSpread  map[string]float64 `json:"political_spread"`
	EngagementScore  float64           `json:"engagement_score"`
	DiversityIndex   float64           `json:"diversity_index"`
	CohesionScore    float64           `json:"cohesion_score"`
	GeneratedAt      time.Time         `json:"generated_at"`
}

// CommunityInteraction represents interactions between community members
type CommunityInteraction struct {
	Id           string                 `json:"id"`
	CommunityId  string                 `json:"community_id"`
	Type         string                 `json:"type"` // "discussion", "conflict", "collaboration", "event"
	Participants []string               `json:"participants"` // member IDs
	Topic        string                 `json:"topic"`
	Sentiment    float64               `json:"sentiment"`    // -1.0 to 1.0
	Intensity    float64               `json:"intensity"`    // 0.0-1.0
	Outcome      string                `json:"outcome"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time             `json:"created_at"`
}
