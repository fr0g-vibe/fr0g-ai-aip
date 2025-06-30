# Community Generation Guide

## Overview

The fr0g-ai-aip system includes advanced community generation capabilities that allow you to create realistic populations of AI personas with diverse demographics, political leanings, interests, and social dynamics.

## Core Concepts

### Personas
Base AI experts with specific knowledge domains and prompts.

### Identities
Instances of personas with rich demographic attributes including:
- Age, gender, education level
- Political leaning and socioeconomic status
- Geographic location and timezone
- Interests and activity levels
- Custom attributes and tags

### Communities
Collections of identities with shared characteristics and configurable diversity metrics.

## Community Generation Configuration

### Age Distribution
```json
{
  "age_distribution": {
    "mean": 35.0,
    "std_dev": 12.0,
    "min_age": 18,
    "max_age": 75,
    "skewness": -0.2
  }
}
```

### Location Constraints
```json
{
  "location_constraint": {
    "type": "city",
    "locations": ["Seattle", "San Francisco", "Austin"],
    "urban": true,
    "timezone": "America/Los_Angeles"
  }
}
```

### Diversity Settings
```json
{
  "political_spread": 0.7,
  "interest_spread": 0.8,
  "socioeconomic_range": 0.6,
  "activity_level": 0.75
}
```

## API Examples

### Generate a Tech Community
```bash
curl -X POST http://localhost:8080/communities/generate \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Silicon Valley Startup Community",
    "type": "professional",
    "description": "Tech entrepreneurs and engineers",
    "target_size": 50,
    "generation_config": {
      "persona_weights": {
        "tech-expert-id": 0.4,
        "business-expert-id": 0.3,
        "designer-expert-id": 0.3
      },
      "age_distribution": {
        "mean": 32,
        "std_dev": 8,
        "min_age": 22,
        "max_age": 55,
        "skewness": 0.3
      },
      "location_constraint": {
        "type": "city",
        "locations": ["San Francisco", "Palo Alto", "Mountain View"],
        "urban": true
      },
      "political_spread": 0.5,
      "interest_spread": 0.9,
      "socioeconomic_range": 0.8,
      "activity_level": 0.8
    }
  }'
```

### Generate a Political Discussion Group
```bash
curl -X POST http://localhost:8080/communities/generate \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Political Discussion Forum",
    "type": "political",
    "description": "Diverse political perspectives",
    "target_size": 30,
    "generation_config": {
      "age_distribution": {
        "mean": 45,
        "std_dev": 15,
        "min_age": 25,
        "max_age": 70
      },
      "political_spread": 1.0,
      "interest_spread": 0.6,
      "socioeconomic_range": 0.9,
      "activity_level": 0.9
    }
  }'
```

## Community Analytics

### Diversity Metrics
- **Shannon Diversity Index**: Measures variety across demographic dimensions
- **Political Distribution**: Breakdown of political leanings
- **Age Demographics**: Mean, median, and distribution statistics
- **Geographic Spread**: Location diversity analysis

### Cohesion Scoring
- **Similarity Measures**: Age, political, and interest similarity
- **Network Density**: Interconnectedness of community members
- **Clustering Factor**: Tendency to form subgroups

### Example Analytics Response
```json
{
  "community_id": "comm_123",
  "member_count": 50,
  "active_members": 42,
  "average_age": 32.4,
  "gender_ratio": {
    "male": 0.52,
    "female": 0.48
  },
  "political_spread": {
    "very_liberal": 0.12,
    "liberal": 0.28,
    "moderate": 0.35,
    "conservative": 0.20,
    "very_conservative": 0.05
  },
  "diversity_index": 0.78,
  "cohesion_score": 0.65,
  "engagement_score": 0.82
}
```

## CLI Usage

### Generate Community
```bash
./bin/fr0g-ai-aip generate-community \
  -name "Research Community" \
  -type "academic" \
  -size 25 \
  -description "University researchers and academics"
```

### List Communities
```bash
./bin/fr0g-ai-aip list-communities
```

### Get Community Statistics
```bash
./bin/fr0g-ai-aip community-stats <community-id>
```

### Add Member to Community
```bash
./bin/fr0g-ai-aip add-member <community-id> <identity-id>
```

### Remove Member from Community
```bash
./bin/fr0g-ai-aip remove-member <community-id> <identity-id>
```

## Use Cases

### Social Research
Create diverse populations to study:
- Political polarization effects
- Demographic interaction patterns
- Social network formation
- Opinion dynamics and influence

### Product Testing
Generate realistic user bases for:
- A/B testing with demographic controls
- Accessibility testing across age groups
- Cultural sensitivity validation
- Market segmentation analysis

### Educational Simulations
Build learning environments with:
- Historical community recreations
- Diverse classroom populations
- Debate and discussion groups
- Cross-cultural interaction scenarios

### Content Moderation
Develop moderation policies using:
- Diverse perspective simulation
- Edge case identification
- Bias detection and mitigation
- Community guideline testing

## Advanced Features

### Custom Attribute Generation
Extend identity attributes with custom fields:
```json
{
  "custom_attributes": {
    "profession": "software_engineer",
    "years_experience": 8,
    "programming_languages": ["Go", "Python", "JavaScript"],
    "remote_work_preference": true
  }
}
```

### Relationship Modeling
Define connections between community members:
```json
{
  "relationships": {
    "network_density": 0.3,
    "clustering_factor": 0.7,
    "influence_patterns": "small_world"
  }
}
```

### Temporal Dynamics
Model community evolution over time:
```json
{
  "temporal_config": {
    "growth_rate": 0.1,
    "churn_rate": 0.05,
    "activity_cycles": "weekly"
  }
}
```

## Best Practices

1. **Start Small**: Begin with communities of 10-25 members for testing
2. **Validate Distributions**: Check generated demographics match expectations
3. **Monitor Diversity**: Ensure adequate representation across dimensions
4. **Test Edge Cases**: Generate communities with extreme configurations
5. **Document Configurations**: Save generation configs for reproducibility

## Troubleshooting

### Common Issues

**Low Diversity Scores**
- Increase `political_spread` and `interest_spread` values
- Expand age distribution parameters
- Add more location options

**Unrealistic Demographics**
- Adjust age distribution mean and standard deviation
- Constrain socioeconomic range appropriately
- Validate persona weight distributions

**Performance Issues**
- Reduce community size for initial testing
- Use memory storage for faster generation
- Optimize persona selection weights

### Error Messages

- `"no personas available for community generation"`: Create personas first
- `"target size must be positive"`: Specify valid community size
- `"referenced persona not found"`: Verify persona IDs in weights config
