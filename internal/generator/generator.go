package generator

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/fr0g-vibe/fr0g-ai-aip/internal/types"
)

// Generator provides methods for creating random and directed identities
type Generator struct{}

// NewGenerator creates a new generator (no longer needs a seeded random number generator)
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateRandomIdentity creates a random identity based on a persona
func (g *Generator) GenerateRandomIdentity(personaID string, name string) *types.Identity {
	identity := &types.Identity{
		PersonaId:      personaID,
		Name:           name,
		Description:    g.generateRandomDescription(),
		Background:     g.generateRandomBackground(),
		IsActive:       true,
		Tags:           g.generateRandomTags(),
		RichAttributes: g.generateRandomRichAttributes(),
	}

	return identity
}

// GenerateDirectedIdentity creates an identity with specific attributes
func (g *Generator) GenerateDirectedIdentity(personaID string, name string,
	demographics *types.Demographics, psychographics *types.Psychographics) *types.Identity {

	identity := &types.Identity{
		PersonaId:   personaID,
		Name:        name,
		Description: g.generateDirectedDescription(demographics, psychographics),
		Background:  g.generateDirectedBackground(demographics, psychographics),
		IsActive:    true,
		Tags:        g.generateDirectedTags(demographics, psychographics),
		RichAttributes: &types.RichAttributes{
			Demographics:         demographics,
			Psychographics:       psychographics,
			LifeHistory:          g.generateDirectedLifeHistory(demographics),
			CulturalReligious:    g.generateDirectedCulturalReligious(demographics),
			PoliticalSocial:      g.generateDirectedPoliticalSocial(psychographics),
			Health:               g.generateDirectedHealth(demographics),
			Preferences:          g.generateDirectedPreferences(psychographics),
			BehavioralTendencies: g.generateDirectedBehavioralTendencies(psychographics),
			CurrentContext:       g.generateDirectedCurrentContext(demographics),
		},
	}

	return identity
}

// GenerateCommunity generates a community of identities with specified demographics
func (g *Generator) GenerateCommunity(personaID string, size int,
	communitySpec *CommunitySpecification) []*types.Identity {

	identities := make([]*types.Identity, size)

	for i := 0; i < size; i++ {
		demographics := g.generateCommunityDemographics(communitySpec)
		psychographics := g.generateCommunityPsychographics(communitySpec)
		name := g.generateName(demographics)

		identities[i] = g.GenerateDirectedIdentity(personaID, name, demographics, psychographics)
	}

	return identities
}

// CommunitySpecification defines the characteristics of a community to generate
type CommunitySpecification struct {
	Location               *types.Location    `json:"location,omitempty"`
	AgeRange               *types.AgeRange    `json:"age_range,omitempty"`
	GenderDistribution     map[string]float64 `json:"gender_distribution,omitempty"` // e.g., {"male": 0.4, "female": 0.6}
	EducationDistribution  map[string]float64 `json:"education_distribution,omitempty"`
	PoliticalDistribution  map[string]float64 `json:"political_distribution,omitempty"`
	UrbanRuralDistribution map[string]float64 `json:"urban_rural_distribution,omitempty"`
	PersonalityProfile     *types.Personality `json:"personality_profile,omitempty"` // Average personality for the community
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

// Helper methods for generating random content
func (g *Generator) generateRandomDescription() string {
	descriptions := []string{
		"A passionate individual with diverse interests and experiences.",
		"Someone who values learning and personal growth.",
		"A community-oriented person with strong connections.",
		"An independent thinker with unique perspectives.",
		"A dedicated professional with a balanced lifestyle.",
	}
	return descriptions[cryptoRandIntn(len(descriptions))]
}

func (g *Generator) generateRandomBackground() string {
	backgrounds := []string{
		"Grew up in a diverse environment that shaped their worldview.",
		"Experienced significant life changes that influenced their perspective.",
		"Developed strong interests through education and personal exploration.",
		"Built a career and life through determination and adaptability.",
		"Formed meaningful relationships that continue to inspire growth.",
	}
	return backgrounds[cryptoRandIntn(len(backgrounds))]
}

func (g *Generator) generateRandomTags() []string {
	allTags := []string{
		"curious", "experienced", "community-minded", "independent", "creative",
		"analytical", "empathetic", "ambitious", "balanced", "adventurous",
		"traditional", "progressive", "practical", "idealistic", "resilient",
	}

	numTags := cryptoRandIntn(4) + 1 // 1-4 tags
	tags := make([]string, numTags)
	used := make(map[string]bool)

	for i := 0; i < numTags; i++ {
		for {
			tag := allTags[cryptoRandIntn(len(allTags))]
			if !used[tag] {
				tags[i] = tag
				used[tag] = true
				break
			}
		}
	}

	return tags
}

func (g *Generator) generateRandomRichAttributes() *types.RichAttributes {
	return &types.RichAttributes{
		Demographics:         g.generateRandomDemographics(),
		Psychographics:       g.generateRandomPsychographics(),
		LifeHistory:          g.generateRandomLifeHistory(),
		CulturalReligious:    g.generateRandomCulturalReligious(),
		PoliticalSocial:      g.generateRandomPoliticalSocial(),
		Health:               g.generateRandomHealth(),
		Preferences:          g.generateRandomPreferences(),
		BehavioralTendencies: g.generateRandomBehavioralTendencies(),
		CurrentContext:       g.generateRandomCurrentContext(),
	}
}

// Additional helper methods would be implemented here...
// For brevity, I'll show a few key ones:

func (g *Generator) generateRandomDemographics() *types.Demographics {
	ages := []int{18, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70}
	genders := []string{"male", "female", "non-binary", "prefer not to say"}
	ethnicities := []string{"White", "Black", "Hispanic", "Asian", "Mixed", "Other"}
	education := []string{"high_school", "bachelors", "masters", "phd"}

	return &types.Demographics{
		Age:       int32(ages[cryptoRandIntn(len(ages))]),
		Gender:    genders[cryptoRandIntn(len(genders))],
		Ethnicity: ethnicities[cryptoRandIntn(len(ethnicities))],
		Education: education[cryptoRandIntn(len(education))],
		Location: &types.Location{
			Country:    "United States",
			City:       "New York",
			UrbanRural: "urban",
		},
	}
}

func (g *Generator) generateRandomPsychographics() *types.Psychographics {
	return &types.Psychographics{
		Personality: &types.Personality{
			Openness:          cryptoRandFloat64(),
			Conscientiousness: cryptoRandFloat64(),
			Extraversion:      cryptoRandFloat64(),
			Agreeableness:     cryptoRandFloat64(),
			Neuroticism:       cryptoRandFloat64(),
		},
		Values:        []string{"honesty", "compassion", "growth"},
		RiskTolerance: "medium",
	}
}

// Placeholder methods for other random generators
func (g *Generator) generateRandomLifeHistory() *types.LifeHistory {
	return &types.LifeHistory{}
}

func (g *Generator) generateRandomCulturalReligious() *types.CulturalReligious {
	return &types.CulturalReligious{}
}

func (g *Generator) generateRandomPoliticalSocial() *types.PoliticalSocial {
	return &types.PoliticalSocial{}
}

func (g *Generator) generateRandomHealth() *types.Health {
	return &types.Health{}
}

func (g *Generator) generateRandomPreferences() *types.Preferences {
	return &types.Preferences{}
}

func (g *Generator) generateRandomBehavioralTendencies() *types.BehavioralTendencies {
	return &types.BehavioralTendencies{}
}

func (g *Generator) generateRandomCurrentContext() *types.CurrentContext {
	return &types.CurrentContext{}
}

// Directed generation methods
func (g *Generator) generateDirectedDescription(demographics *types.Demographics, psychographics *types.Psychographics) string {
	return "A directed identity with specific characteristics."
}

func (g *Generator) generateDirectedBackground(demographics *types.Demographics, psychographics *types.Psychographics) string {
	return "Background shaped by directed attributes."
}

func (g *Generator) generateDirectedTags(demographics *types.Demographics, psychographics *types.Psychographics) []string {
	return []string{"directed", "specific"}
}

func (g *Generator) generateDirectedLifeHistory(demographics *types.Demographics) *types.LifeHistory {
	return &types.LifeHistory{}
}

func (g *Generator) generateDirectedCulturalReligious(demographics *types.Demographics) *types.CulturalReligious {
	return &types.CulturalReligious{}
}

func (g *Generator) generateDirectedPoliticalSocial(psychographics *types.Psychographics) *types.PoliticalSocial {
	return &types.PoliticalSocial{}
}

func (g *Generator) generateDirectedHealth(demographics *types.Demographics) *types.Health {
	return &types.Health{}
}

func (g *Generator) generateDirectedPreferences(psychographics *types.Psychographics) *types.Preferences {
	return &types.Preferences{}
}

func (g *Generator) generateDirectedBehavioralTendencies(psychographics *types.Psychographics) *types.BehavioralTendencies {
	return &types.BehavioralTendencies{}
}

func (g *Generator) generateDirectedCurrentContext(demographics *types.Demographics) *types.CurrentContext {
	return &types.CurrentContext{}
}

// Community generation methods
func (g *Generator) generateCommunityDemographics(spec *CommunitySpecification) *types.Demographics {
	demographics := &types.Demographics{}

	if spec.AgeRange != nil {
		demographics.Age = int32(cryptoRandIntn(int(spec.AgeRange.Max-spec.AgeRange.Min+1)) + int(spec.AgeRange.Min))
	}

	if spec.Location != nil {
		demographics.Location = spec.Location
	}

	// Apply gender distribution if specified
	if spec.GenderDistribution != nil {
		demographics.Gender = g.selectFromDistribution(spec.GenderDistribution)
	}

	return demographics
}

func (g *Generator) generateCommunityPsychographics(spec *CommunitySpecification) *types.Psychographics {
	psychographics := &types.Psychographics{}

	if spec.PersonalityProfile != nil {
		// Add some variation around the community average
		psychographics.Personality = &types.Personality{
			Openness:          g.addVariation(spec.PersonalityProfile.Openness, 0.2),
			Conscientiousness: g.addVariation(spec.PersonalityProfile.Conscientiousness, 0.2),
			Extraversion:      g.addVariation(spec.PersonalityProfile.Extraversion, 0.2),
			Agreeableness:     g.addVariation(spec.PersonalityProfile.Agreeableness, 0.2),
			Neuroticism:       g.addVariation(spec.PersonalityProfile.Neuroticism, 0.2),
		}
	}

	return psychographics
}

func (g *Generator) generateName(demographics *types.Demographics) string {
	// Simple name generation - in a real implementation, you'd use a proper name database
	firstNames := []string{"Alex", "Jordan", "Casey", "Taylor", "Morgan", "Riley", "Quinn", "Avery"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis"}

	return firstNames[cryptoRandIntn(len(firstNames))] + " " + lastNames[cryptoRandIntn(len(lastNames))]
}

// Utility methods
func (g *Generator) selectFromDistribution(distribution map[string]float64) string {
	r := cryptoRandFloat64()
	cumulative := 0.0

	for key, prob := range distribution {
		cumulative += prob
		if r <= cumulative {
			return key
		}
	}

	// Fallback to first key
	for key := range distribution {
		return key
	}
	return ""
}

func (g *Generator) addVariation(base float64, variation float64) float64 {
	change := (cryptoRandFloat64() - 0.5) * 2 * variation
	result := base + change
	if result < 0 {
		return 0
	}
	if result > 1 {
		return 1
	}
	return result
}
