package community

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service provides community generation and management functionality
type Service struct {
	storage storage.Storage
	rand    *rand.Rand
}

// NewService creates a new community service
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
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
		
		// Set age if the field exists
		if age, ok := richAttrs["age"].(int); ok {
			if hasField(identity.RichAttributes, "Age") {
				setField(identity.RichAttributes, "Age", int32(age))
			}
		}
		
		// Set other attributes using a generic approach
		if gender, ok := richAttrs["gender"].(string); ok {
			setStringField(identity.RichAttributes, "Gender", gender)
		}
		if political, ok := richAttrs["political_leaning"].(string); ok {
			setStringField(identity.RichAttributes, "PoliticalLeaning", political)
		}
		if socioeconomic, ok := richAttrs["socioeconomic_status"].(string); ok {
			setStringField(identity.RichAttributes, "SocioeconomicStatus", socioeconomic)
		}
		if education, ok := richAttrs["education"].(string); ok {
			setStringField(identity.RichAttributes, "Education", education)
		}
		if activity, ok := richAttrs["activity_level"].(float64); ok {
			setFloatField(identity.RichAttributes, "ActivityLevel", activity)
		}
		
		// Handle location if supported
		if loc, ok := richAttrs["location"].(map[string]interface{}); ok {
			setLocationField(identity.RichAttributes, loc)
		}
		
		// Handle interests if supported
		if interests, ok := richAttrs["interests"].([]string); ok {
			setStringSliceField(identity.RichAttributes, "Interests", interests)
		}

		members = append(members, identity)
	}

	return members, nil
}

// selectPersonaByWeight selects a persona based on configured weights
func (s *Service) selectPersonaByWeight(personas []types.Persona, weights map[string]float64) types.Persona {
	if len(weights) == 0 {
		// Equal probability if no weights specified
		return personas[s.rand.Intn(len(personas))]
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
	target := s.rand.Float64() * totalWeight
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
	age := s.rand.NormFloat64()*dist.StdDev + dist.Mean
	
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
			city := constraint.Locations[s.rand.Intn(len(constraint.Locations))]
			location["city"] = city
			location["type"] = "city"
		} else {
			location["city"] = s.generateRandomCity()
			location["type"] = "city"
		}
	case "region":
		if len(constraint.Locations) > 0 {
			region := constraint.Locations[s.rand.Intn(len(constraint.Locations))]
			location["region"] = region
			location["type"] = "region"
		}
	case "country":
		if len(constraint.Locations) > 0 {
			country := constraint.Locations[s.rand.Intn(len(constraint.Locations))]
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
		location["urban"] = s.rand.Float64() > 0.3 // 70% urban by default
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
	value := s.rand.NormFloat64()*spread + center

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
	value := s.rand.Float64() * spread

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
	indices := s.rand.Perm(len(allInterests))
	
	for i := 0; i < numInterests; i++ {
		selected = append(selected, allInterests[indices[i]])
	}

	return selected
}

// generateActivityLevel creates activity level based on community config
func (s *Service) generateActivityLevel(baseLevel float64) float64 {
	// Add some randomness around the base level
	variation := s.rand.NormFloat64() * 0.2 // 20% standard deviation
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
	if s.rand.Float64() < 0.5 {
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

	value := s.rand.Float64()
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

	firstName := firstNames[s.rand.Intn(len(firstNames))]
	lastName := lastNames[s.rand.Intn(len(lastNames))]
	
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
	return cities[s.rand.Intn(len(cities))]
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
			val = getFieldValue(member.RichAttributes, "Age")
		case "political_leaning":
			val = getFieldValue(member.RichAttributes, "PoliticalLeaning")
		case "gender":
			val = getFieldValue(member.RichAttributes, "Gender")
		case "education":
			val = getFieldValue(member.RichAttributes, "Education")
		case "socioeconomic_status":
			val = getFieldValue(member.RichAttributes, "SocioeconomicStatus")
		case "location":
			location := getFieldValue(member.RichAttributes, "Location")
			if location != nil {
				val = getFieldValue(location, "City")
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
		if member.RichAttributes != nil {
			interests := getFieldValue(member.RichAttributes, "Interests")
			if interestSlice, ok := interests.([]string); ok && len(interestSlice) > 0 {
				for _, interest := range interestSlice {
					allInterests[interest]++
					totalInterests++
				}
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
	if member1.RichAttributes != nil && member2.RichAttributes != nil {
		age1Val := getFieldValue(member1.RichAttributes, "Age")
		age2Val := getFieldValue(member2.RichAttributes, "Age")
		if age1Val != nil && age2Val != nil {
			age1 := convertToInt(age1Val)
			age2 := convertToInt(age2Val)
			ageDiff := math.Abs(float64(age1 - age2))
			ageSim := math.Max(0, 1.0-ageDiff/50.0) // Normalize by 50-year span
			similarities = append(similarities, ageSim)
		}
	}

	// Political similarity
	if member1.RichAttributes != nil && member2.RichAttributes != nil {
		pol1Val := getFieldValue(member1.RichAttributes, "PoliticalLeaning")
		pol2Val := getFieldValue(member2.RichAttributes, "PoliticalLeaning")
		if pol1Str, ok1 := pol1Val.(string); ok1 && pol1Str != "" {
			if pol2Str, ok2 := pol2Val.(string); ok2 && pol2Str != "" {
				polSim := s.calculatePoliticalSimilarity(pol1Str, pol2Str)
				similarities = append(similarities, polSim)
			}
		}
	}

	// Interest similarity
	if member1.RichAttributes != nil && member2.RichAttributes != nil {
		int1Val := getFieldValue(member1.RichAttributes, "Interests")
		int2Val := getFieldValue(member2.RichAttributes, "Interests")
		if int1Slice, ok1 := int1Val.([]string); ok1 && len(int1Slice) > 0 {
			if int2Slice, ok2 := int2Val.([]string); ok2 && len(int2Slice) > 0 {
				intSim := s.calculateInterestSimilarity(int1Slice, int2Slice)
				similarities = append(similarities, intSim)
			}
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
		if member.RichAttributes != nil {
			ageVal := getFieldValue(member.RichAttributes, "Age")
			if ageVal != nil {
				age := convertToInt(ageVal)
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
		if member.RichAttributes != nil {
			polVal := getFieldValue(member.RichAttributes, "PoliticalLeaning")
			if polStr, ok := polVal.(string); ok && polStr != "" {
				distribution[polStr]++
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
		if member.RichAttributes != nil {
			locationVal := getFieldValue(member.RichAttributes, "Location")
			if locationVal != nil {
				cityVal := getFieldValue(locationVal, "City")
				if cityStr, ok := cityVal.(string); ok && cityStr != "" {
					locations[cityStr]++
				} else {
					typeVal := getFieldValue(locationVal, "Type")
					if typeStr, ok := typeVal.(string); ok && typeStr != "" {
						locations[typeStr]++
					}
				}
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
			activityVal := getFieldValue(member.RichAttributes, "ActivityLevel")
			if activityFloat, ok := activityVal.(float64); ok && activityFloat > 0.5 {
				activeCount++
			}
		}
	}
	stats.ActiveMembers = activeCount

	// Calculate gender ratio
	genderCount := make(map[string]int)
	for _, member := range members {
		if member.RichAttributes != nil {
			genderVal := getFieldValue(member.RichAttributes, "Gender")
			if genderStr, ok := genderVal.(string); ok && genderStr != "" {
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
			activityVal := getFieldValue(member.RichAttributes, "ActivityLevel")
			if activityFloat, ok := activityVal.(float64); ok {
				totalActivity += activityFloat
				activityCount++
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

// Helper functions to safely set fields on protobuf structs using reflection
func hasField(v interface{}, fieldName string) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}
	return val.FieldByName(fieldName).IsValid()
}

func setField(v interface{}, fieldName string, value interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	field := val.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		fieldValue := reflect.ValueOf(value)
		if field.Type() == fieldValue.Type() {
			field.Set(fieldValue)
		}
	}
}

func setStringField(v interface{}, fieldName string, value string) {
	setField(v, fieldName, value)
}

func setFloatField(v interface{}, fieldName string, value float64) {
	setField(v, fieldName, value)
}

func setStringSliceField(v interface{}, fieldName string, value []string) {
	setField(v, fieldName, value)
}

func setLocationField(richAttrs *types.RichAttributes, loc map[string]interface{}) {
	// Try to set location using reflection
	val := reflect.ValueOf(richAttrs)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	locationField := val.FieldByName("Location")
	if locationField.IsValid() && locationField.CanSet() {
		// Create a new Location struct if the field exists
		locationType := locationField.Type()
		if locationType.Kind() == reflect.Ptr {
			locationType = locationType.Elem()
		}
		
		newLocation := reflect.New(locationType)
		locationValue := newLocation.Elem()
		
		// Set location fields if they exist
		if city, ok := loc["city"].(string); ok {
			cityField := locationValue.FieldByName("City")
			if cityField.IsValid() && cityField.CanSet() {
				cityField.SetString(city)
			}
		}
		if locType, ok := loc["type"].(string); ok {
			typeField := locationValue.FieldByName("Type")
			if typeField.IsValid() && typeField.CanSet() {
				typeField.SetString(locType)
			}
		}
		if urban, ok := loc["urban"].(bool); ok {
			urbanField := locationValue.FieldByName("Urban")
			if urbanField.IsValid() && urbanField.CanSet() {
				urbanField.SetBool(urban)
			}
		}
		if timezone, ok := loc["timezone"].(string); ok {
			timezoneField := locationValue.FieldByName("Timezone")
			if timezoneField.IsValid() && timezoneField.CanSet() {
				timezoneField.SetString(timezone)
			}
		}
		
		locationField.Set(newLocation)
	}
}

func getFieldValue(v interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}
	return field.Interface()
}

func convertToInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	case float32:
		return int(val)
	default:
		return 0
	}
}
