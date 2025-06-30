# API Documentation

## Overview

The fr0g-ai-aip system provides comprehensive APIs for managing AI personas, identities, and communities. This document covers all available endpoints, data models, and usage examples.

## Base URLs

- **HTTP REST API**: `http://localhost:8080`
- **gRPC API**: `localhost:9090`
- **Health Check**: `http://localhost:8080/health`

## Authentication

Currently, the API supports optional API key authentication via the `X-API-Key` header.

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/personas
```

## Data Models

### Persona

A Persona represents an AI subject matter expert with specific knowledge and capabilities.

```json
{
  "id": "string",
  "name": "string",
  "topic": "string", 
  "prompt": "string",
  "context": {
    "key": "value"
  },
  "rag": ["string"]
}
```

**Fields:**
- `id`: Unique identifier (auto-generated)
- `name`: Display name for the persona (required, 1-100 chars)
- `topic`: Subject area or domain (required, 1-100 chars)
- `prompt`: System prompt for the AI (required, 1-10000 chars)
- `context`: Key-value pairs for additional context (optional)
- `rag`: Array of RAG document references (optional)

### Identity

An Identity represents an instance of a persona with rich demographic and behavioral attributes.

```json
{
  "id": "string",
  "persona_id": "string",
  "name": "string",
  "description": "string",
  "background": "string",
  "rich_attributes": {
    "age": 32,
    "gender": "female",
    "political_leaning": "moderate",
    "education": "bachelor",
    "socioeconomic_status": "middle",
    "activity_level": 0.75,
    "location": {
      "city": "Seattle",
      "type": "city",
      "urban": true,
      "timezone": "America/Los_Angeles"
    },
    "interests": ["technology", "music", "travel"]
  },
  "tags": ["string"],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Community

A Community represents a collection of identities with shared characteristics and configurable diversity.

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "type": "professional",
  "size": 50,
  "diversity": 0.78,
  "cohesion": 0.65,
  "member_ids": ["string"],
  "max_members": 100,
  "min_members": 10,
  "generation_config": {
    "persona_weights": {
      "persona_id": 0.4
    },
    "age_distribution": {
      "mean": 35.0,
      "std_dev": 12.0,
      "min_age": 18,
      "max_age": 75,
      "skewness": -0.2
    },
    "political_spread": 0.7,
    "interest_spread": 0.8,
    "socioeconomic_range": 0.6,
    "activity_level": 0.75
  },
  "attributes": {},
  "tags": ["string"],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Persona Endpoints

### Create Persona

**POST** `/personas`

Creates a new AI persona.

**Request Body:**
```json
{
  "name": "Security Expert",
  "topic": "Cybersecurity",
  "prompt": "You are a cybersecurity expert with extensive knowledge of threat analysis, security best practices, and incident response.",
  "context": {
    "domain": "enterprise security",
    "experience": "15 years"
  },
  "rag": [
    "security-frameworks.md",
    "incident-response-playbook.pdf"
  ]
}
```

**Response:** `201 Created`
```json
{
  "id": "abc123",
  "name": "Security Expert",
  "topic": "Cybersecurity",
  "prompt": "You are a cybersecurity expert...",
  "context": {
    "domain": "enterprise security",
    "experience": "15 years"
  },
  "rag": [
    "security-frameworks.md",
    "incident-response-playbook.pdf"
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input data
- `422 Unprocessable Entity`: Validation errors

### Get Persona

**GET** `/personas/{id}`

Retrieves a specific persona by ID.

**Response:** `200 OK`
```json
{
  "id": "abc123",
  "name": "Security Expert",
  "topic": "Cybersecurity",
  "prompt": "You are a cybersecurity expert...",
  "context": {
    "domain": "enterprise security"
  },
  "rag": ["security-frameworks.md"]
}
```

**Error Responses:**
- `404 Not Found`: Persona does not exist

### List Personas

**GET** `/personas`

Retrieves all personas.

**Response:** `200 OK`
```json
[
  {
    "id": "abc123",
    "name": "Security Expert",
    "topic": "Cybersecurity",
    "prompt": "You are a cybersecurity expert..."
  },
  {
    "id": "def456", 
    "name": "Go Developer",
    "topic": "Golang Programming",
    "prompt": "You are an expert Go programmer..."
  }
]
```

### Update Persona

**PUT** `/personas/{id}`

Updates an existing persona.

**Request Body:**
```json
{
  "name": "Updated Security Expert",
  "topic": "Advanced Cybersecurity",
  "prompt": "You are an advanced cybersecurity expert...",
  "context": {
    "domain": "enterprise security",
    "experience": "20 years",
    "certifications": "CISSP, CISM"
  }
}
```

**Response:** `200 OK`
```json
{
  "id": "abc123",
  "name": "Updated Security Expert",
  "topic": "Advanced Cybersecurity",
  "prompt": "You are an advanced cybersecurity expert...",
  "context": {
    "domain": "enterprise security",
    "experience": "20 years",
    "certifications": "CISSP, CISM"
  }
}
```

**Error Responses:**
- `404 Not Found`: Persona does not exist
- `400 Bad Request`: Invalid input data

### Delete Persona

**DELETE** `/personas/{id}`

Deletes a persona.

**Response:** `204 No Content`

**Error Responses:**
- `404 Not Found`: Persona does not exist

## Identity Endpoints

### Create Identity

**POST** `/identities`

Creates a new identity based on a persona.

**Request Body:**
```json
{
  "persona_id": "abc123",
  "name": "Alice Johnson",
  "description": "Senior cybersecurity analyst",
  "background": "10 years experience in enterprise security",
  "rich_attributes": {
    "age": 32,
    "gender": "female",
    "political_leaning": "moderate",
    "education": "master",
    "socioeconomic_status": "upper_middle",
    "activity_level": 0.8,
    "location": {
      "city": "Seattle",
      "type": "city",
      "urban": true
    },
    "interests": ["cybersecurity", "technology", "hiking"]
  },
  "tags": ["security", "analyst", "senior"]
}
```

**Response:** `201 Created`
```json
{
  "id": "identity123",
  "persona_id": "abc123",
  "name": "Alice Johnson",
  "description": "Senior cybersecurity analyst",
  "background": "10 years experience in enterprise security",
  "rich_attributes": {
    "age": 32,
    "gender": "female",
    "political_leaning": "moderate",
    "education": "master",
    "socioeconomic_status": "upper_middle",
    "activity_level": 0.8,
    "location": {
      "city": "Seattle",
      "type": "city", 
      "urban": true
    },
    "interests": ["cybersecurity", "technology", "hiking"]
  },
  "tags": ["security", "analyst", "senior"],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Get Identity

**GET** `/identities/{id}`

Retrieves a specific identity.

**Response:** `200 OK`
```json
{
  "id": "identity123",
  "persona_id": "abc123",
  "name": "Alice Johnson",
  "description": "Senior cybersecurity analyst",
  "rich_attributes": {
    "age": 32,
    "gender": "female"
  },
  "tags": ["security", "analyst"],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Get Identity with Persona

**GET** `/identities/{id}/with-persona`

Retrieves an identity along with its associated persona data.

**Response:** `200 OK`
```json
{
  "identity": {
    "id": "identity123",
    "persona_id": "abc123",
    "name": "Alice Johnson",
    "description": "Senior cybersecurity analyst"
  },
  "persona": {
    "id": "abc123",
    "name": "Security Expert",
    "topic": "Cybersecurity",
    "prompt": "You are a cybersecurity expert..."
  }
}
```

### List Identities

**GET** `/identities`

Retrieves identities with optional filtering.

**Query Parameters:**
- `persona_id`: Filter by persona ID
- `tags`: Filter by tags (comma-separated)
- `is_active`: Filter by active status (true/false)
- `search`: Search in name and description

**Example:**
```bash
GET /identities?persona_id=abc123&tags=security,analyst&is_active=true
```

**Response:** `200 OK`
```json
[
  {
    "id": "identity123",
    "persona_id": "abc123",
    "name": "Alice Johnson",
    "description": "Senior cybersecurity analyst",
    "tags": ["security", "analyst"],
    "is_active": true
  }
]
```

### Update Identity

**PUT** `/identities/{id}`

Updates an existing identity.

**Request Body:**
```json
{
  "persona_id": "abc123",
  "name": "Alice Johnson-Smith",
  "description": "Lead cybersecurity analyst",
  "rich_attributes": {
    "age": 33,
    "activity_level": 0.9
  },
  "tags": ["security", "analyst", "lead"]
}
```

**Response:** `200 OK`

### Delete Identity

**DELETE** `/identities/{id}`

Deletes an identity.

**Response:** `204 No Content`

## Community Endpoints

### Generate Community

**POST** `/communities/generate`

Generates a new community with specified characteristics.

**Request Body:**
```json
{
  "name": "Tech Startup Community",
  "type": "professional",
  "description": "Software developers and entrepreneurs in tech startups",
  "target_size": 25,
  "generation_config": {
    "persona_weights": {
      "tech-expert-id": 0.4,
      "business-expert-id": 0.3,
      "designer-expert-id": 0.3
    },
    "age_distribution": {
      "mean": 30,
      "std_dev": 8,
      "min_age": 22,
      "max_age": 55,
      "skewness": 0.3
    },
    "location_constraint": {
      "type": "city",
      "locations": ["San Francisco", "Seattle", "Austin"],
      "urban": true
    },
    "political_spread": 0.6,
    "interest_spread": 0.8,
    "socioeconomic_range": 0.7,
    "activity_level": 0.8
  }
}
```

**Response:** `201 Created`
```json
{
  "id": "community123",
  "name": "Tech Startup Community",
  "type": "professional",
  "description": "Software developers and entrepreneurs in tech startups",
  "size": 25,
  "diversity": 0.78,
  "cohesion": 0.65,
  "member_ids": ["identity1", "identity2", "..."],
  "generation_config": {
    "persona_weights": {
      "tech-expert-id": 0.4
    },
    "age_distribution": {
      "mean": 30,
      "std_dev": 8
    }
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Get Community

**GET** `/communities/{id}`

Retrieves a specific community.

**Response:** `200 OK`
```json
{
  "id": "community123",
  "name": "Tech Startup Community",
  "type": "professional",
  "size": 25,
  "diversity": 0.78,
  "cohesion": 0.65,
  "member_ids": ["identity1", "identity2"],
  "attributes": {
    "average_age": 32.4,
    "political_distribution": {
      "liberal": 0.4,
      "moderate": 0.4,
      "conservative": 0.2
    }
  }
}
```

### List Communities

**GET** `/communities`

Retrieves communities with optional filtering.

**Query Parameters:**
- `type`: Filter by community type
- `tags`: Filter by tags
- `is_active`: Filter by active status
- `min_size`: Minimum community size
- `max_size`: Maximum community size
- `min_diversity`: Minimum diversity score
- `max_diversity`: Maximum diversity score
- `search`: Search in name and description

### Get Community Statistics

**GET** `/communities/{id}/stats`

Retrieves detailed analytics for a community.

**Response:** `200 OK`
```json
{
  "community_id": "community123",
  "member_count": 25,
  "active_members": 22,
  "average_age": 32.4,
  "gender_ratio": {
    "male": 0.52,
    "female": 0.48
  },
  "political_spread": {
    "very_liberal": 0.08,
    "liberal": 0.32,
    "moderate": 0.40,
    "conservative": 0.16,
    "very_conservative": 0.04
  },
  "location_spread": {
    "San Francisco": 12,
    "Seattle": 8,
    "Austin": 5
  },
  "diversity_index": 0.78,
  "cohesion_score": 0.65,
  "engagement_score": 0.82,
  "generated_at": "2024-01-01T00:00:00Z"
}
```

### Add Member to Community

**POST** `/communities/{id}/members`

Adds an existing identity to a community.

**Request Body:**
```json
{
  "identity_id": "identity456"
}
```

**Response:** `200 OK`

### Remove Member from Community

**DELETE** `/communities/{id}/members/{identity_id}`

Removes a member from a community.

**Response:** `204 No Content`

### Update Community

**PUT** `/communities/{id}`

Updates community metadata (not members).

**Request Body:**
```json
{
  "name": "Updated Tech Community",
  "description": "Updated description",
  "tags": ["tech", "startup", "updated"]
}
```

**Response:** `200 OK`

### Delete Community

**DELETE** `/communities/{id}`

Deletes a community (does not delete member identities).

**Response:** `204 No Content`

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "details": {
    "field": "Additional error details"
  }
}
```

**Common Error Codes:**
- `400 Bad Request`: Invalid request format or parameters
- `401 Unauthorized`: Missing or invalid authentication
- `404 Not Found`: Resource does not exist
- `422 Unprocessable Entity`: Validation errors
- `500 Internal Server Error`: Server error

## Rate Limiting

The API implements rate limiting:
- **Default**: 100 requests per minute per IP
- **Authenticated**: 1000 requests per minute per API key

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## Pagination

List endpoints support pagination:

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

**Response Headers:**
```
X-Total-Count: 150
X-Page: 1
X-Per-Page: 20
X-Total-Pages: 8
```

## Webhooks

The system supports webhooks for real-time notifications:

**Supported Events:**
- `persona.created`
- `persona.updated`
- `persona.deleted`
- `identity.created`
- `identity.updated`
- `identity.deleted`
- `community.created`
- `community.updated`
- `community.deleted`

**Webhook Payload:**
```json
{
  "event": "persona.created",
  "timestamp": "2024-01-01T00:00:00Z",
  "data": {
    "id": "abc123",
    "name": "Security Expert"
  }
}
```

## SDK Examples

### Go SDK

```go
package main

import (
    "github.com/fr0g-vibe/fr0g-ai-aip/pkg/client"
)

func main() {
    client := client.NewRESTClient("http://localhost:8080")
    
    // Create persona
    persona := &types.Persona{
        Name:   "Go Expert",
        Topic:  "Golang Programming", 
        Prompt: "You are an expert Go programmer...",
    }
    
    err := client.Create(persona)
    if err != nil {
        log.Fatal(err)
    }
    
    // List personas
    personas, err := client.List()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d personas\n", len(personas))
}
```

### Python SDK

```python
from fr0g_ai_aip import Client

client = Client("http://localhost:8080")

# Create persona
persona = client.personas.create({
    "name": "Python Expert",
    "topic": "Python Programming",
    "prompt": "You are an expert Python programmer..."
})

# List personas
personas = client.personas.list()
print(f"Found {len(personas)} personas")
```

### JavaScript SDK

```javascript
const { Client } = require('@fr0g-vibe/fr0g-ai-aip');

const client = new Client('http://localhost:8080');

// Create persona
const persona = await client.personas.create({
  name: 'JavaScript Expert',
  topic: 'JavaScript Programming',
  prompt: 'You are an expert JavaScript programmer...'
});

// List personas
const personas = await client.personas.list();
console.log(`Found ${personas.length} personas`);
```

## MCP Integration

The system is designed for Model Context Protocol (MCP) integration:

### MCP Server Configuration

```json
{
  "mcpServers": {
    "fr0g-ai-aip": {
      "command": "fr0g-ai-aip",
      "args": ["-mcp"],
      "env": {
        "FR0G_SERVER_URL": "http://localhost:8080"
      }
    }
  }
}
```

### MCP Tools

The system exposes these MCP tools:

- `get_persona`: Retrieve persona by ID or name
- `list_personas`: List all available personas
- `create_identity`: Create identity from persona
- `get_community_stats`: Get community analytics
- `generate_community`: Create new community

### MCP Resources

Available MCP resources:

- `persona://{id}`: Individual persona data
- `identity://{id}`: Individual identity data  
- `community://{id}`: Community data and stats
- `community://{id}/members`: Community member list

## Performance

**Typical Response Times:**
- GET requests: < 10ms
- POST/PUT requests: < 50ms
- Community generation: 100ms - 2s (depending on size)

**Throughput:**
- Personas: 1000+ ops/sec
- Identities: 500+ ops/sec  
- Communities: 10+ generations/sec

**Storage:**
- Memory: Unlimited (RAM-bound)
- File: Unlimited (disk-bound)
- Recommended: < 10,000 personas, < 100,000 identities per instance
