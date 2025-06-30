package community

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Service provides community generation and management functionality
type Service struct {
	storage storage.Storage
}

// NewService creates a new community service
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// GenerateCommunity creates a new community with generated members
func (s *Service) GenerateCommunity(config types.CommunityGenerationConfig, name, description, communityType string, targetSize int) (*types.Community, error) {
	if targetSize <= 0 {
		return nil, fmt.Errorf("target size must be positive")
	}

	// Create the community structure
	community := &types.Community{
		Id:               generateID(),
		Name:             name,
		Description:      description,
		Type:             communityType,
		Size:             0,
		MemberIds:        make([]string, 0, targetSize),
		MaxMembers:       targetSize * 2, // Allow for growth
		MinMembers:       max(1, targetSize/2),
		GenerationConfig: config,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		IsActive:         true,
		Tags:             []string{},
		Attributes:       make(map[string]interface{}),
	}

	// Generate community members
	members, err := s.generateMembers(config, targetSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate members: %v", err)
	}

	// Add members to community
	for _, member := range members {
		if err := s.storage.CreateIdentity(&member); err != nil {
			return nil, fmt.Errorf("failed to create member identity: %v", err)
		}
		community.MemberIds = append(community.MemberIds, member.Id)
	}

	community.Size = len(community.MemberIds)

	// Calculate community metrics
	s.calculateCommunityMetrics(community, members)

	// Store the community
	if err := s.storage.CreateCommunity(community); err != nil {
		return nil, fmt.Errorf("failed to store community: %v", err)
	}

	return community, nil
}

// generateMembers creates identities based on the generation configuration
func (s *Service) generateMembers(config types.CommunityGenerationConfig, count int) ([]types.Identity, error) {
	// Get available personas
	personas, err := s.storage.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get personas: %v", err)
	}

	if len(personas) == 0 {
		return nil, fmt.Errorf("no personas available for community generation")
	}

	members := make([]types.Identity, 0, count)

	for i := 0; i < count; i++ {
		// Select persona based on weights
		persona := s.selectPersonaByWeight(personas, config.PersonaWeights)

		// Generate identity attributes
		now := time.Now()
		identity := types.Identity{
			Id:          generateID(),
			PersonaId:   persona.Id,
			Name:        s.generateName(),
			Description: fmt.Sprintf("Community member based on %s persona", persona.Name),
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
			Tags:        []string{"community-generated"},
		}

		// Generate rich attributes based on community config
		richAttrs := s.generateRichAttributes(config, i, count)

		// Create RichAttributes with available fields
		identity.RichAttributes = &types.RichAttributes{}
		// Set Demographics
		dem := &types.Demographics{}
		if age, ok := richAttrs["age"].(int); ok {
			dem.Age = int32(age)
		}
		if gender, ok := richAttrs["gender"].(string); ok {
			dem.Gender = gender
		}
		if education, ok := richAttrs["education"].(string); ok {
			dem.Education = education
		}
		if socioeconomic, ok := richAttrs["socioeconomic_status"].(string); ok {
			dem.SocioeconomicStatus = socioeconomic
		}
		if loc, ok := richAttrs["location"].(map[string]interface{}); ok {
			dem.Location = mapToLocation(loc)
		}
		identity.RichAttributes.Demographics = dem
		// Set PoliticalSocial
		if political, ok := richAttrs["political_leaning"].(string); ok {
			ps := &types.PoliticalSocial{PoliticalLeaning: political}
			identity.RichAttributes.PoliticalSocial = ps
		}
		// Set Preferences
		if interests, ok := richAttrs["interests"].([]string); ok {
			prefs := &types.Preferences{Interests: interests}
			identity.RichAttributes.Preferences = prefs
		}
		// Set activity level in custom map
		if activity, ok := richAttrs["activity_level"].(float64); ok {
			if identity.RichAttributes.Custom == nil {
				identity.RichAttributes.Custom = make(map[string]string)
			}
			identity.RichAttributes.Custom["activity_level"] = fmt.Sprintf("%f", activity)
		}
		members = append(members, identity)
	}

	return members, nil
}

// selectPersonaByWeight selects a persona based on configured weights
func (s *Service) selectPersonaByWeight(personas []types.Persona, weights map[string]float64) types.Persona {
	if len(weights) == 0 {
		// Equal probability if no weights specified
		return personas[cryptoRandIntn(len(personas))]
	}

	// Calculate total weight for available personas
	totalWeight := 0.0
	availablePersonas := make([]types.Persona, 0)
	personaWeights := make([]float64, 0)

	for _, persona := range personas {
		weight := weights[persona.Id]
		if weight <= 0 {
			weight = 1.0 // Default weight
		}
		availablePersonas = append(availablePersonas, persona)
		personaWeights = append(personaWeights, weight)
		totalWeight += weight
	}

	// Select based on weighted random
	target := cryptoRandFloat64() * totalWeight
	current := 0.0

	for i, weight := range personaWeights {
		current += weight
		if current >= target {
			return availablePersonas[i]
		}
	}

	// Fallback
	return availablePersonas[len(availablePersonas)-1]
}

// generateRichAttributes creates realistic attributes for a community member
func (s *Service) generateRichAttributes(config types.CommunityGenerationConfig, memberIndex, totalMembers int) map[string]interface{} {
	attrs := make(map[string]interface{})

	// Generate age based on distribution
	age := s.generateAge(config.AgeDistribution)
	attrs["age"] = age

	// Generate location based on constraints
	location := s.generateLocation(config.LocationConstraint)
	attrs["location"] = location

	// Generate political leaning with specified spread
	politicalLeaning := s.generatePoliticalLeaning(config.PoliticalSpread)
	attrs["political_leaning"] = politicalLeaning

	// Generate socioeconomic status
	socioeconomicStatus := s.generateSocioeconomicStatus(config.SocioeconomicRange)
	attrs["socioeconomic_status"] = socioeconomicStatus

	// Generate interests with diversity
	interests := s.generateInterests(config.InterestSpread)
	attrs["interests"] = interests

	// Generate activity level
	activityLevel := s.generateActivityLevel(config.ActivityLevel)
	attrs["activity_level"] = activityLevel

	// Generate gender (simple binary for now, could be expanded)
	gender := s.generateGender()
	attrs["gender"] = gender

	// Generate education level
	education := s.generateEducationLevel(age)
	attrs["education"] = education

	return attrs
}

// generateAge creates an age based on the distribution parameters
func (s *Service) generateAge(dist types.AgeDistribution) int {
	// Use normal distribution with constraints
	age := cryptoRandNormFloat64()*dist.StdDev + dist.Mean

	// Apply skewness (simple implementation)
	if dist.Skewness != 0 {
		skewAdjustment := dist.Skewness * (age - dist.Mean) * 0.5
		age += skewAdjustment
	}

	// Constrain to bounds
	ageInt := int(math.Round(age))
	if ageInt < dist.MinAge {
		ageInt = dist.MinAge
	}
	if ageInt > dist.MaxAge {
		ageInt = dist.MaxAge
	}

	return ageInt
}

// generateLocation creates a location based on constraints
func (s *Service) generateLocation(constraint types.LocationConstraint) map[string]interface{} {
	location := make(map[string]interface{})

	switch constraint.Type {
	case "city":
		if len(constraint.Locations) > 0 {
			city := constraint.Locations[cryptoRandIntn(len(constraint.Locations))]
			location["city"] = city
			location["type"] = "city"
		} else {
			location["city"] = s.generateRandomCity()
			location["type"] = "city"
		}
	case "region":
		if len(constraint.Locations) > 0 {
			region := constraint.Locations[cryptoRandIntn(len(constraint.Locations))]
			location["region"] = region
			location["type"] = "region"
		}
	case "country":
		if len(constraint.Locations) > 0 {
			country := constraint.Locations[cryptoRandIntn(len(constraint.Locations))]
			location["country"] = country
			location["type"] = "country"
		}
	default:
		location["type"] = "global"
	}

	// Add urban/rural designation
	if constraint.Urban != nil {
		location["urban"] = *constraint.Urban
	} else {
		location["urban"] = cryptoRandFloat64() > 0.3 // 70% urban by default
	}

	if constraint.Timezone != "" {
		location["timezone"] = constraint.Timezone
	}

	return location
}

// generatePoliticalLeaning creates political orientation with specified spread
func (s *Service) generatePoliticalLeaning(spread float64) string {
	// Generate value from -1 (very liberal) to 1 (very conservative)
	center := 0.0 // Neutral center
	value := cryptoRandNormFloat64()*spread + center

	// Constrain to [-1, 1]
	if value < -1 {
		value = -1
	}
	if value > 1 {
		value = 1
	}

	// Convert to categorical
	if value < -0.6 {
		return "very_liberal"
	} else if value < -0.2 {
		return "liberal"
	} else if value < 0.2 {
		return "moderate"
	} else if value < 0.6 {
		return "conservative"
	} else {
		return "very_conservative"
	}
}

// generateSocioeconomicStatus creates economic status with specified range
func (s *Service) generateSocioeconomicStatus(spread float64) string {
	value := cryptoRandFloat64() * spread

	if value < 0.2 {
		return "low_income"
	} else if value < 0.4 {
		return "lower_middle"
	} else if value < 0.6 {
		return "middle"
	} else if value < 0.8 {
		return "upper_middle"
	} else {
		return "high_income"
	}
}

// generateInterests creates a list of interests with specified diversity
func (s *Service) generateInterests(diversity float64) []string {
	allInterests := []string{
		"technology", "sports", "music", "art", "cooking", "travel", "reading",
		"gaming", "fitness", "photography", "gardening", "movies", "politics",
		"science", "history", "fashion", "cars", "pets", "crafts", "business",
	}

	// Number of interests based on diversity (more diversity = more varied interests)
	numInterests := int(diversity*10) + 2 // 2-12 interests
	if numInterests > len(allInterests) {
		numInterests = len(allInterests)
	}

	// Randomly select interests
	selected := make([]string, 0, numInterests)
	indices := cryptoRandPerm(len(allInterests))

	for i := 0; i < numInterests; i++ {
		selected = append(selected, allInterests[indices[i]])
	}

	return selected
}

// generateActivityLevel creates activity level based on community config
func (s *Service) generateActivityLevel(baseLevel float64) float64 {
	// Add some randomness around the base level
	variation := cryptoRandNormFloat64() * 0.2 // 20% standard deviation
	level := baseLevel + variation

	// Constrain to [0, 1]
	if level < 0 {
		level = 0
	}
	if level > 1 {
		level = 1
	}

	return level
}

// generateGender creates gender (simplified binary model)
func (s *Service) generateGender() string {
	if cryptoRandFloat64() < 0.5 {
		return "male"
	}
	return "female"
}

// generateEducationLevel creates education level based on age
func (s *Service) generateEducationLevel(age int) string {
	// Younger people more likely to have higher education
	baseProb := 0.3
	if age < 30 {
		baseProb = 0.6
	} else if age < 50 {
		baseProb = 0.4
	}

	value := cryptoRandFloat64()
	if value < baseProb*0.3 {
		return "graduate"
	} else if value < baseProb*0.7 {
		return "bachelor"
	} else if value < baseProb {
		return "associate"
	} else if value < 0.8 {
		return "high_school"
	} else {
		return "some_high_school"
	}
}

// generateName creates a random name for community members
func (s *Service) generateName() string {
	firstNames := []string{
		"Alex", "Jordan", "Taylor", "Casey", "Morgan", "Riley", "Avery", "Quinn",
		"Sam", "Blake", "Cameron", "Drew", "Emery", "Finley", "Harper", "Hayden",
		"Jamie", "Kendall", "Logan", "Parker", "Peyton", "Reese", "Sage", "Skyler",
	}

	lastNames := []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
		"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
		"Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White",
	}

	firstName := firstNames[cryptoRandIntn(len(firstNames))]
	lastName := lastNames[cryptoRandIntn(len(lastNames))]

	return fmt.Sprintf("%s %s", firstName, lastName)
}

// generateRandomCity creates a random city name
func (s *Service) generateRandomCity() string {
	cities := []string{
		"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia",
		"San Antonio", "San Diego", "Dallas", "San Jose", "Austin", "Jacksonville",
		"Fort Worth", "Columbus", "Charlotte", "San Francisco", "Indianapolis", "Seattle",
		"Denver", "Washington", "Boston", "El Paso", "Nashville", "Detroit", "Portland",
	}
	return cities[cryptoRandIntn(len(cities))]
}

// calculateCommunityMetrics computes diversity and cohesion scores
func (s *Service) calculateCommunityMetrics(community *types.Community, members []types.Identity) {
	if len(members) == 0 {
		community.Diversity = 0
		community.Cohesion = 0
		return
	}

	// Calculate diversity based on attribute variation
	diversity := s.calculateDiversityIndex(members)
	community.Diversity = diversity

	// Calculate cohesion based on similarity and network density
	cohesion := s.calculateCohesionScore(members, community.GenerationConfig)
	community.Cohesion = cohesion

	// Set community attributes
	community.Attributes["average_age"] = s.calculateAverageAge(members)
	community.Attributes["political_distribution"] = s.calculatePoliticalDistribution(members)
	community.Attributes["location_spread"] = s.calculateLocationSpread(members)
}

// calculateDiversityIndex computes a diversity score based on member attributes
func (s *Service) calculateDiversityIndex(members []types.Identity) float64 {
	if len(members) <= 1 {
		return 0
	}

	// Calculate diversity across multiple dimensions
	ageDiversity := s.calculateAttributeDiversity(members, "age")
	politicalDiversity := s.calculateAttributeDiversity(members, "political_leaning")
	locationDiversity := s.calculateAttributeDiversity(members, "location")
	interestDiversity := s.calculateInterestDiversity(members)

	// Weighted average of diversity measures
	return (ageDiversity*0.25 + politicalDiversity*0.25 + locationDiversity*0.25 + interestDiversity*0.25)
}

// calculateAttributeDiversity computes diversity for a specific attribute
func (s *Service) calculateAttributeDiversity(members []types.Identity, attribute string) float64 {
	values := make(map[interface{}]int)
	for _, member := range members {
		if member.RichAttributes == nil {
			continue
		}
		var val interface{}
		switch attribute {
		case "age":
			if member.RichAttributes.Demographics != nil {
				val = member.RichAttributes.Demographics.Age
			}
		case "political_leaning":
			if member.RichAttributes.PoliticalSocial != nil {
				val = member.RichAttributes.PoliticalSocial.PoliticalLeaning
			}
		case "gender":
			if member.RichAttributes.Demographics != nil {
				val = member.RichAttributes.Demographics.Gender
			}
		case "education":
			if member.RichAttributes.Demographics != nil {
				val = member.RichAttributes.Demographics.Education
			}
		case "socioeconomic_status":
			if member.RichAttributes.Demographics != nil {
				val = member.RichAttributes.Demographics.SocioeconomicStatus
			}
		case "location":
			if member.RichAttributes.Demographics != nil && member.RichAttributes.Demographics.Location != nil {
				val = member.RichAttributes.Demographics.Location.City
			}
		default:
			continue
		}
		if val != nil && val != "" && val != int32(0) && val != 0 {
			values[val]++
		}
	}
	if len(values) <= 1 {
		return 0
	}
	// Shannon diversity index
	total := float64(len(members))
	diversity := 0.0

	for _, count := range values {
		if count > 0 {
			p := float64(count) / total
			diversity -= p * math.Log2(p)
		}
	}

	// Normalize by maximum possible diversity
	maxDiversity := math.Log2(float64(len(values)))
	if maxDiversity > 0 {
		return diversity / maxDiversity
	}

	return 0
}

// calculateInterestDiversity computes diversity of interests across members
func (s *Service) calculateInterestDiversity(members []types.Identity) float64 {
	allInterests := make(map[string]int)
	totalInterests := 0
	for _, member := range members {
		if member.RichAttributes != nil && member.RichAttributes.Preferences != nil && len(member.RichAttributes.Preferences.Interests) > 0 {
			for _, interest := range member.RichAttributes.Preferences.Interests {
				allInterests[interest]++
				totalInterests++
			}
		}
	}
	if len(allInterests) <= 1 || totalInterests == 0 {
		return 0
	}
	// Shannon diversity for interests
	diversity := 0.0
	total := float64(totalInterests)

	for _, count := range allInterests {
		if count > 0 {
			p := float64(count) / total
			diversity -= p * math.Log2(p)
		}
	}

	// Normalize
	maxDiversity := math.Log2(float64(len(allInterests)))
	if maxDiversity > 0 {
		return diversity / maxDiversity
	}

	return 0
}

// calculateCohesionScore computes how cohesive the community is
func (s *Service) calculateCohesionScore(members []types.Identity, config types.CommunityGenerationConfig) float64 {
	if len(members) <= 1 {
		return 1.0
	}

	// Cohesion based on similarity in key attributes
	similarities := make([]float64, 0)

	for i := 0; i < len(members); i++ {
		for j := i + 1; j < len(members); j++ {
			similarity := s.calculateMemberSimilarity(members[i], members[j])
			similarities = append(similarities, similarity)
		}
	}

	if len(similarities) == 0 {
		return 0
	}

	// Average similarity
	total := 0.0
	for _, sim := range similarities {
		total += sim
	}

	return total / float64(len(similarities))
}

// calculateMemberSimilarity computes similarity between two members
func (s *Service) calculateMemberSimilarity(member1, member2 types.Identity) float64 {
	similarities := make([]float64, 0)

	// Age similarity
	if member1.RichAttributes != nil && member2.RichAttributes != nil &&
		member1.RichAttributes.Demographics != nil && member2.RichAttributes.Demographics != nil {
		age1 := int(member1.RichAttributes.Demographics.Age)
		age2 := int(member2.RichAttributes.Demographics.Age)
		ageDiff := math.Abs(float64(age1 - age2))
		ageSim := math.Max(0, 1.0-ageDiff/50.0) // Normalize by 50-year span
		similarities = append(similarities, ageSim)
	}

	// Political similarity
	if member1.RichAttributes != nil && member2.RichAttributes != nil &&
		member1.RichAttributes.PoliticalSocial != nil && member2.RichAttributes.PoliticalSocial != nil {
		pol1Str := member1.RichAttributes.PoliticalSocial.PoliticalLeaning
		pol2Str := member2.RichAttributes.PoliticalSocial.PoliticalLeaning
		if pol1Str != "" && pol2Str != "" {
			polSim := s.calculatePoliticalSimilarity(pol1Str, pol2Str)
			similarities = append(similarities, polSim)
		}
	}

	// Interest similarity
	if member1.RichAttributes != nil && member2.RichAttributes != nil &&
		member1.RichAttributes.Preferences != nil && member2.RichAttributes.Preferences != nil {
		int1Slice := member1.RichAttributes.Preferences.Interests
		int2Slice := member2.RichAttributes.Preferences.Interests
		if len(int1Slice) > 0 && len(int2Slice) > 0 {
			intSim := s.calculateInterestSimilarity(int1Slice, int2Slice)
			similarities = append(similarities, intSim)
		}
	}

	if len(similarities) == 0 {
		return 0
	}

	// Average similarity across attributes
	total := 0.0
	for _, sim := range similarities {
		total += sim
	}

	return total / float64(len(similarities))
}

// calculatePoliticalSimilarity computes similarity between political leanings
func (s *Service) calculatePoliticalSimilarity(pol1, pol2 string) float64 {
	politicalOrder := map[string]int{
		"very_liberal":      0,
		"liberal":           1,
		"moderate":          2,
		"conservative":      3,
		"very_conservative": 4,
	}

	val1, ok1 := politicalOrder[pol1]
	val2, ok2 := politicalOrder[pol2]

	if !ok1 || !ok2 {
		return 0
	}

	diff := math.Abs(float64(val1 - val2))
	return math.Max(0, 1.0-diff/4.0) // Normalize by maximum difference
}

// calculateInterestSimilarity computes Jaccard similarity between interest lists
func (s *Service) calculateInterestSimilarity(interests1, interests2 []string) float64 {
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, interest := range interests1 {
		set1[interest] = true
	}
	for _, interest := range interests2 {
		set2[interest] = true
	}

	intersection := 0
	union := len(set1)

	for interest := range set2 {
		if set1[interest] {
			intersection++
		} else {
			union++
		}
	}

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// Helper functions for community attributes
func (s *Service) calculateAverageAge(members []types.Identity) float64 {
	total := 0.0
	count := 0
	for _, member := range members {
		if member.RichAttributes != nil && member.RichAttributes.Demographics != nil {
			age := int(member.RichAttributes.Demographics.Age)
			if age > 0 {
				total += float64(age)
				count++
			}
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func (s *Service) calculatePoliticalDistribution(members []types.Identity) map[string]float64 {
	distribution := make(map[string]int)
	total := 0
	for _, member := range members {
		if member.RichAttributes != nil && member.RichAttributes.PoliticalSocial != nil {
			pol := member.RichAttributes.PoliticalSocial.PoliticalLeaning
			if pol != "" {
				distribution[pol]++
				total++
			}
		}
	}
	result := make(map[string]float64)
	for pol, count := range distribution {
		result[pol] = float64(count) / float64(total)
	}
	return result
}

func (s *Service) calculateLocationSpread(members []types.Identity) map[string]int {
	locations := make(map[string]int)
	for _, member := range members {
		if member.RichAttributes != nil && member.RichAttributes.Demographics != nil && member.RichAttributes.Demographics.Location != nil {
			city := member.RichAttributes.Demographics.Location.City
			if city != "" {
				locations[city]++
			} else if member.RichAttributes.Demographics.Location.UrbanRural != "" {
				locations[member.RichAttributes.Demographics.Location.UrbanRural]++
			}
		}
	}
	return locations
}

// GetCommunity retrieves a community by ID
func (s *Service) GetCommunity(id string) (types.Community, error) {
	return s.storage.GetCommunity(id)
}

// ListCommunities returns communities with optional filtering
func (s *Service) ListCommunities(filter *types.CommunityFilter) ([]types.Community, error) {
	return s.storage.ListCommunities(filter)
}

// UpdateCommunity updates an existing community
func (s *Service) UpdateCommunity(id string, community types.Community) error {
	community.UpdatedAt = time.Now()
	return s.storage.UpdateCommunity(id, community)
}

// DeleteCommunity removes a community
func (s *Service) DeleteCommunity(id string) error {
	return s.storage.DeleteCommunity(id)
}

// GetCommunityStats generates analytics for a community
func (s *Service) GetCommunityStats(communityId string) (*types.CommunityStats, error) {
	community, err := s.storage.GetCommunity(communityId)
	if err != nil {
		return nil, err
	}

	// Get all community members
	members := make([]types.Identity, 0, len(community.MemberIds))
	for _, memberId := range community.MemberIds {
		member, err := s.storage.GetIdentity(memberId)
		if err != nil {
			continue // Skip missing members
		}
		members = append(members, member)
	}

	stats := &types.CommunityStats{
		CommunityId:     communityId,
		MemberCount:     len(members),
		AverageAge:      s.calculateAverageAge(members),
		LocationSpread:  s.calculateLocationSpread(members),
		PoliticalSpread: s.calculatePoliticalDistribution(members),
		DiversityIndex:  s.calculateDiversityIndex(members),
		CohesionScore:   s.calculateCohesionScore(members, community.GenerationConfig),
		GeneratedAt:     time.Now(),
	}

	// Calculate active members (activity_level > 0.5)
	activeCount := 0
	for _, member := range members {
		if member.RichAttributes != nil {
			activityVal := member.RichAttributes.Custom["activity_level"]
			if activityVal != "" {
				if activityFloat, err := strconv.ParseFloat(activityVal, 64); err == nil && activityFloat > 0.5 {
					activeCount++
				}
			}
		}
	}
	stats.ActiveMembers = activeCount

	// Calculate gender ratio
	genderCount := make(map[string]int)
	for _, member := range members {
		if member.RichAttributes != nil && member.RichAttributes.Demographics != nil {
			genderStr := member.RichAttributes.Demographics.Gender
			if genderStr != "" {
				genderCount[genderStr]++
			}
		}
	}

	genderRatio := make(map[string]float64)
	for gender, count := range genderCount {
		genderRatio[gender] = float64(count) / float64(len(members))
	}
	stats.GenderRatio = genderRatio

	// Calculate engagement score (average activity level)
	totalActivity := 0.0
	activityCount := 0
	for _, member := range members {
		if member.RichAttributes != nil {
			activityVal := member.RichAttributes.Custom["activity_level"]
			if activityVal != "" {
				if activityFloat, err := strconv.ParseFloat(activityVal, 64); err == nil {
					totalActivity += activityFloat
					activityCount++
				}
			}
		}
	}
	if activityCount > 0 {
		stats.EngagementScore = totalActivity / float64(activityCount)
	}

	return stats, nil
}

// AddMemberToCommunity adds an existing identity to a community
func (s *Service) AddMemberToCommunity(communityId, identityId string) error {
	community, err := s.storage.GetCommunity(communityId)
	if err != nil {
		return err
	}

	// Check if identity exists
	_, err = s.storage.GetIdentity(identityId)
	if err != nil {
		return fmt.Errorf("identity not found: %v", err)
	}

	// Check if already a member
	for _, memberId := range community.MemberIds {
		if memberId == identityId {
			return fmt.Errorf("identity is already a member of this community")
		}
	}

	// Check size limits
	if len(community.MemberIds) >= community.MaxMembers {
		return fmt.Errorf("community has reached maximum size")
	}

	// Add member
	community.MemberIds = append(community.MemberIds, identityId)
	community.Size = len(community.MemberIds)
	community.UpdatedAt = time.Now()

	return s.storage.UpdateCommunity(communityId, community)
}

// RemoveMemberFromCommunity removes a member from a community
func (s *Service) RemoveMemberFromCommunity(communityId, identityId string) error {
	community, err := s.storage.GetCommunity(communityId)
	if err != nil {
		return err
	}

	// Find and remove member
	newMemberIds := make([]string, 0, len(community.MemberIds))
	found := false

	for _, memberId := range community.MemberIds {
		if memberId != identityId {
			newMemberIds = append(newMemberIds, memberId)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("identity is not a member of this community")
	}

	// Check minimum size
	if len(newMemberIds) < community.MinMembers {
		return fmt.Errorf("removing member would violate minimum community size")
	}

	community.MemberIds = newMemberIds
	community.Size = len(community.MemberIds)
	community.UpdatedAt = time.Now()

	return s.storage.UpdateCommunity(communityId, community)
}

// Utility functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper to convert map[string]interface{} to *types.Location
func mapToLocation(loc map[string]interface{}) *types.Location {
	l := &types.Location{}
	if city, ok := loc["city"].(string); ok {
		l.City = city
	}
	if locType, ok := loc["type"].(string); ok {
		l.UrbanRural = locType // Map "type" to UrbanRural if present
	}
	if timezone, ok := loc["timezone"].(string); ok {
		l.Timezone = timezone
	}
	if country, ok := loc["country"].(string); ok {
		l.Country = country
	}
	if region, ok := loc["region"].(string); ok {
		l.Region = region
	}
	return l
}

// Helper functions for cryptographically secure random numbers
func cryptoRandIntn(max int) int {
	if max <= 0 {
		return 0
	}
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	num := binary.BigEndian.Uint64(b[:])
	return int(num % uint64(max))
}

func cryptoRandFloat64() float64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	num := binary.BigEndian.Uint64(b[:])
	return float64(num) / (1 << 64)
}

// Box-Muller transform for normal distribution
func cryptoRandNormFloat64() float64 {
	u1 := cryptoRandFloat64()
	u2 := cryptoRandFloat64()
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return z0
}

// Fisher-Yates shuffle for permutation
func cryptoRandPerm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		m[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := cryptoRandIntn(i + 1)
		m[i], m[j] = m[j], m[i]
	}
	return m
}
