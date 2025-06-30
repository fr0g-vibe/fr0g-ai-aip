# Model Context Protocol (MCP) Integration

This document describes how fr0g-ai-aip integrates with the Model Context Protocol (MCP) to provide AI personas as contextual tools and resources.

## Overview

The fr0g-ai-aip system is designed to work seamlessly with MCP-compatible AI systems, providing:

- **Tools**: Functions that AI models can call to interact with personas and communities
- **Resources**: Structured data that AI models can access for context
- **Prompts**: Pre-configured prompts and templates for common use cases

## MCP Server Configuration

### Basic Configuration

Add fr0g-ai-aip to your MCP configuration:

```json
{
  "mcpServers": {
    "fr0g-ai-aip": {
      "command": "fr0g-ai-aip",
      "args": ["-mcp"],
      "env": {
        "FR0G_SERVER_URL": "http://localhost:8080",
        "FR0G_CLIENT_TYPE": "rest"
      }
    }
  }
}
```

### Advanced Configuration

For production deployments with authentication:

```json
{
  "mcpServers": {
    "fr0g-ai-aip": {
      "command": "fr0g-ai-aip",
      "args": ["-mcp", "-config", "/path/to/config.yaml"],
      "env": {
        "FR0G_SERVER_URL": "https://api.fr0g-ai-aip.com",
        "FR0G_API_KEY": "${FR0G_API_KEY}",
        "FR0G_CLIENT_TYPE": "rest"
      }
    }
  }
}
```

## Available Tools

### Persona Management Tools

#### `get_persona`

Retrieve a specific persona by ID or name.

**Parameters:**
- `id` (string, optional): Persona ID
- `name` (string, optional): Persona name (if ID not provided)

**Returns:**
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
  "rag": ["security-frameworks.md"]
}
```

**Example Usage:**
```
Get the Security Expert persona to help with threat analysis.
```

#### `list_personas`

List all available personas with optional filtering.

**Parameters:**
- `topic` (string, optional): Filter by topic
- `limit` (integer, optional): Maximum number of results (default: 20)

**Returns:**
```json
[
  {
    "id": "abc123",
    "name": "Security Expert",
    "topic": "Cybersecurity"
  },
  {
    "id": "def456", 
    "name": "Go Developer",
    "topic": "Golang Programming"
  }
]
```

#### `create_persona`

Create a new persona (requires appropriate permissions).

**Parameters:**
- `name` (string, required): Persona name
- `topic` (string, required): Subject area
- `prompt` (string, required): System prompt
- `context` (object, optional): Additional context
- `rag` (array, optional): RAG document references

**Returns:**
```json
{
  "id": "new123",
  "name": "New Expert",
  "topic": "New Domain",
  "prompt": "You are an expert in..."
}
```

### Identity Management Tools

#### `create_identity`

Create an identity instance from a persona.

**Parameters:**
- `persona_id` (string, required): Base persona ID
- `name` (string, required): Identity name
- `description` (string, optional): Description
- `attributes` (object, optional): Rich attributes

**Returns:**
```json
{
  "id": "identity123",
  "persona_id": "abc123",
  "name": "Alice Johnson",
  "description": "Senior security analyst",
  "rich_attributes": {
    "age": 32,
    "experience": "10 years",
    "location": "Seattle"
  }
}
```

#### `get_identity`

Retrieve an identity with its persona context.

**Parameters:**
- `id` (string, required): Identity ID
- `include_persona` (boolean, optional): Include full persona data

**Returns:**
```json
{
  "identity": {
    "id": "identity123",
    "name": "Alice Johnson",
    "description": "Senior security analyst"
  },
  "persona": {
    "id": "abc123",
    "name": "Security Expert",
    "prompt": "You are a cybersecurity expert..."
  }
}
```

### Community Tools

#### `generate_community`

Generate a community with specified characteristics.

**Parameters:**
- `name` (string, required): Community name
- `type` (string, required): Community type
- `size` (integer, required): Target size
- `config` (object, required): Generation configuration

**Returns:**
```json
{
  "id": "community123",
  "name": "Tech Startup Team",
  "type": "professional",
  "size": 25,
  "diversity": 0.78,
  "member_ids": ["identity1", "identity2", "..."]
}
```

#### `get_community_stats`

Get detailed analytics for a community.

**Parameters:**
- `id` (string, required): Community ID

**Returns:**
```json
{
  "community_id": "community123",
  "member_count": 25,
  "average_age": 32.4,
  "diversity_index": 0.78,
  "political_spread": {
    "liberal": 0.4,
    "moderate": 0.4,
    "conservative": 0.2
  },
  "engagement_score": 0.82
}
```

#### `simulate_community_discussion`

Simulate a discussion within a community on a given topic.

**Parameters:**
- `community_id` (string, required): Community ID
- `topic` (string, required): Discussion topic
- `participants` (integer, optional): Number of participants (default: 5)
- `rounds` (integer, optional): Discussion rounds (default: 3)

**Returns:**
```json
{
  "topic": "Remote work policies",
  "participants": [
    {
      "identity_id": "identity1",
      "name": "Alice Johnson",
      "perspective": "Supports flexible remote work",
      "contributions": ["Initial statement", "Response to Bob", "Final thoughts"]
    }
  ],
  "summary": "The discussion revealed diverse perspectives on remote work...",
  "sentiment": 0.6,
  "consensus_level": 0.4
}
```

## Available Resources

### Persona Resources

#### `persona://{id}`

Access individual persona data.

**URI Pattern:** `persona://abc123`

**Content:**
```json
{
  "id": "abc123",
  "name": "Security Expert",
  "topic": "Cybersecurity",
  "prompt": "You are a cybersecurity expert with extensive knowledge...",
  "context": {
    "domain": "enterprise security",
    "experience": "15 years",
    "certifications": "CISSP, CISM"
  },
  "rag": [
    "security-frameworks.md",
    "incident-response-playbook.pdf",
    "threat-intelligence-feeds.json"
  ]
}
```

#### `persona://list`

Access the complete persona catalog.

**URI Pattern:** `persona://list`

**Content:**
```json
{
  "personas": [
    {
      "id": "abc123",
      "name": "Security Expert",
      "topic": "Cybersecurity",
      "description": "Expert in enterprise security and threat analysis"
    },
    {
      "id": "def456",
      "name": "Go Developer", 
      "topic": "Golang Programming",
      "description": "Expert in Go programming and best practices"
    }
  ],
  "total_count": 2,
  "last_updated": "2024-01-01T00:00:00Z"
}
```

### Identity Resources

#### `identity://{id}`

Access individual identity data with persona context.

**URI Pattern:** `identity://identity123`

**Content:**
```json
{
  "identity": {
    "id": "identity123",
    "persona_id": "abc123",
    "name": "Alice Johnson",
    "description": "Senior cybersecurity analyst at TechCorp",
    "background": "10 years experience in enterprise security...",
    "rich_attributes": {
      "age": 32,
      "gender": "female",
      "education": "master",
      "location": {
        "city": "Seattle",
        "timezone": "America/Los_Angeles"
      },
      "interests": ["cybersecurity", "hiking", "photography"]
    }
  },
  "persona": {
    "id": "abc123",
    "name": "Security Expert",
    "prompt": "You are a cybersecurity expert..."
  }
}
```

### Community Resources

#### `community://{id}`

Access community data and member information.

**URI Pattern:** `community://community123`

**Content:**
```json
{
  "community": {
    "id": "community123",
    "name": "Tech Startup Team",
    "type": "professional",
    "size": 25,
    "diversity": 0.78,
    "cohesion": 0.65
  },
  "members": [
    {
      "identity_id": "identity1",
      "name": "Alice Johnson",
      "role": "security_lead",
      "influence": 0.8
    }
  ],
  "statistics": {
    "average_age": 32.4,
    "gender_ratio": {"male": 0.52, "female": 0.48},
    "political_spread": {"liberal": 0.4, "moderate": 0.4, "conservative": 0.2}
  }
}
```

#### `community://{id}/members`

Access detailed member list for a community.

**URI Pattern:** `community://community123/members`

**Content:**
```json
{
  "community_id": "community123",
  "member_count": 25,
  "members": [
    {
      "identity": {
        "id": "identity1",
        "name": "Alice Johnson",
        "description": "Senior security analyst"
      },
      "persona": {
        "id": "abc123",
        "name": "Security Expert",
        "topic": "Cybersecurity"
      },
      "community_role": "security_lead",
      "influence_score": 0.8,
      "activity_level": 0.9
    }
  ]
}
```

## Available Prompts

### Persona Consultation Prompts

#### `consult_expert`

Template for consulting with a specific persona.

**Parameters:**
- `persona_id`: ID of the persona to consult
- `question`: Question or topic to discuss
- `context`: Additional context for the consultation

**Template:**
```
You are now acting as {persona.name}, {persona.topic}.

{persona.prompt}

Additional context:
{context}

Please respond to the following question from your perspective as {persona.name}:

{question}

Consider your expertise in {persona.topic} and any relevant context provided above.
```

#### `multi_expert_panel`

Template for consulting multiple personas on a topic.

**Parameters:**
- `persona_ids`: Array of persona IDs
- `topic`: Discussion topic
- `format`: Response format (discussion, debate, consensus)

**Template:**
```
You are facilitating a panel discussion between multiple experts on the topic: {topic}

Panel members:
{#each personas}
- {name}: {topic} expert
{/each}

Please provide perspectives from each expert, considering their unique expertise and viewpoints. Format the response as a {format}.

Each expert should:
1. Present their initial perspective
2. Respond to other experts' points
3. Provide actionable recommendations

Topic for discussion: {topic}
```

### Community Simulation Prompts

#### `community_discussion`

Template for simulating community discussions.

**Parameters:**
- `community_id`: Community ID
- `topic`: Discussion topic
- `participant_count`: Number of participants

**Template:**
```
Simulate a discussion within the {community.name} community on the topic: {topic}

Community characteristics:
- Type: {community.type}
- Size: {community.size} members
- Diversity: {community.diversity}
- Political spread: {community.political_distribution}

Participants ({participant_count} members):
{#each participants}
- {name}: {description} (Political leaning: {political_leaning}, Activity: {activity_level})
{/each}

Simulate a realistic discussion where each participant contributes based on their background, political views, and personality. Include:
1. Initial positions from each participant
2. Interactions and responses between members
3. Evolution of the discussion
4. Potential consensus or disagreement areas

Topic: {topic}
```

#### `demographic_analysis`

Template for analyzing demographic impacts on opinions.

**Parameters:**
- `community_id`: Community ID
- `issue`: Issue to analyze
- `demographic_focus`: Specific demographic to focus on

**Template:**
```
Analyze how different demographic groups within {community.name} might respond to: {issue}

Community demographics:
- Average age: {community.average_age}
- Gender ratio: {community.gender_ratio}
- Education levels: {community.education_distribution}
- Geographic spread: {community.location_spread}

Focus on {demographic_focus} and provide:
1. Likely perspectives from different demographic segments
2. Potential areas of agreement and disagreement
3. Factors that might influence opinions
4. Recommendations for addressing diverse viewpoints

Issue to analyze: {issue}
```

## Integration Examples

### Claude Desktop Integration

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

### Custom MCP Client

```python
import asyncio
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def main():
    server_params = StdioServerParameters(
        command="fr0g-ai-aip",
        args=["-mcp"],
        env={"FR0G_SERVER_URL": "http://localhost:8080"}
    )
    
    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # Initialize the session
            await session.initialize()
            
            # List available tools
            tools = await session.list_tools()
            print(f"Available tools: {[tool.name for tool in tools.tools]}")
            
            # Get a persona
            result = await session.call_tool("get_persona", {"name": "Security Expert"})
            print(f"Persona: {result.content}")
            
            # Create an identity
            identity_result = await session.call_tool("create_identity", {
                "persona_id": result.content["id"],
                "name": "Alice Johnson",
                "description": "Senior security analyst"
            })
            print(f"Created identity: {identity_result.content}")

if __name__ == "__main__":
    asyncio.run(main())
```

### JavaScript/Node.js Integration

```javascript
const { Client } = require('@modelcontextprotocol/sdk/client/index.js');
const { StdioClientTransport } = require('@modelcontextprotocol/sdk/client/stdio.js');

async function main() {
  const transport = new StdioClientTransport({
    command: 'fr0g-ai-aip',
    args: ['-mcp'],
    env: { FR0G_SERVER_URL: 'http://localhost:8080' }
  });

  const client = new Client({
    name: "fr0g-ai-aip-client",
    version: "1.0.0"
  }, {
    capabilities: {}
  });

  await client.connect(transport);

  // List personas
  const personas = await client.callTool({
    name: "list_personas",
    arguments: { limit: 10 }
  });

  console.log('Available personas:', personas.content);

  // Generate a community
  const community = await client.callTool({
    name: "generate_community",
    arguments: {
      name: "Test Community",
      type: "professional",
      size: 10,
      config: {
        age_distribution: {
          mean: 35,
          std_dev: 10,
          min_age: 25,
          max_age: 55
        },
        political_spread: 0.6,
        interest_spread: 0.8
      }
    }
  });

  console.log('Generated community:', community.content);

  await client.close();
}

main().catch(console.error);
```

## Best Practices

### Tool Usage

1. **Cache persona data**: Personas don't change frequently, so cache them locally
2. **Batch operations**: Use list operations when working with multiple items
3. **Error handling**: Always handle potential errors from tool calls
4. **Rate limiting**: Respect API rate limits when making frequent calls

### Resource Access

1. **Use specific URIs**: Access specific resources rather than listing when possible
2. **Monitor updates**: Check resource timestamps to detect changes
3. **Efficient filtering**: Use query parameters to filter data at the source

### Prompt Templates

1. **Parameterize effectively**: Use clear parameter names and provide defaults
2. **Include context**: Always provide relevant context for better responses
3. **Format consistently**: Use consistent formatting for better AI understanding

## Troubleshooting

### Common Issues

**Tool not found**
```
Error: Tool 'get_persona' not found
```
- Ensure fr0g-ai-aip MCP server is running
- Check MCP configuration is correct
- Verify tool name spelling

**Connection refused**
```
Error: Connection refused to http://localhost:8080
```
- Start the fr0g-ai-aip HTTP server: `fr0g-ai-aip -server`
- Check server URL in environment variables
- Verify firewall settings

**Authentication failed**
```
Error: 401 Unauthorized
```
- Set FR0G_API_KEY environment variable
- Check API key is valid
- Verify authentication configuration

### Debug Mode

Enable debug logging:

```bash
FR0G_LOG_LEVEL=debug fr0g-ai-aip -mcp
```

### Health Checks

Test MCP server health:

```bash
# Test basic connectivity
curl http://localhost:8080/health

# Test MCP server
echo '{"jsonrpc": "2.0", "method": "initialize", "id": 1}' | fr0g-ai-aip -mcp
```
