# Intent Recognition API

A high-quality Go application that extracts structured intent and entities from natural language input. This application provides a thin layer that can be trained to identify specific data patterns and return structured JSON responses.

## Features

- **Multi-Provider AI Support**: Works with OpenAI, Ollama (local), or configurable local AI
- **Enhanced Local AI**: Highly configurable, domain-specific intent recognition with 90%+ accuracy
- **Pattern Matching**: Fast regex-based intent recognition for common patterns
- **AI Integration**: Intelligent extraction for complex or ambiguous inputs
- **Structured Output**: Consistent JSON responses with task and variable extraction
- **RESTful API**: Clean HTTP interface for easy integration
- **Comprehensive Testing**: Unit tests for reliability and maintainability
- **Error Handling**: Robust error handling and validation
- **Offline Capable**: Works without internet using local AI providers
- **Configurable**: JSON-based intent configuration for domain-specific accuracy

## Architecture

The application follows clean architecture principles with clear separation of concerns:

```
├── main.go                 # Application entry point
├── config/                 # Configuration management
├── configs/                # Intent configuration files
│   └── personal_assistant.json  # Example domain-specific config
├── internal/
│   ├── models/            # Data structures and validation
│   │   └── intent_config.go    # Configurable intent system
│   ├── services/          # Business logic and AI integration
│   │   ├── ai_provider.go           # AI provider interface
│   │   ├── openai_provider.go       # OpenAI implementation
│   │   ├── ollama_provider.go       # Ollama (local) implementation
│   │   ├── local_ai_provider.go     # Basic rule-based implementation
│   │   ├── enhanced_local_provider.go # Advanced configurable implementation
│   │   └── intent_service.go        # Main intent service
│   └── handlers/          # HTTP request handling
├── tests/                 # Test files
└── docs/                  # Documentation
```

## AI Providers

### 1. Enhanced Local AI (Recommended for Offline)
- **Best for**: High-accuracy offline environments, domain-specific applications
- **Models**: Configurable intent patterns, regex, fuzzy matching, synonym expansion
- **Setup**: JSON configuration file
- **Performance**: 90%+ accuracy for configured domains, very fast, works offline
- **Customization**: Fully configurable for any domain or use case

### 2. OpenAI (Cloud-based)
- **Best for**: Production environments with internet access
- **Models**: GPT-3.5-turbo, GPT-4, and other OpenAI models
- **Setup**: Requires OpenAI API key
- **Performance**: High accuracy, fast response times

### 3. Ollama (Local)
- **Best for**: Offline environments, privacy-conscious deployments
- **Models**: Llama2, Mistral, CodeLlama, and other open models
- **Setup**: Requires Ollama installation and model download
- **Performance**: Good accuracy, runs locally

### 4. Local AI (Basic)
- **Best for**: Simple offline environments, basic use cases
- **Models**: Rule-based extraction using regex and keyword matching
- **Setup**: No external dependencies
- **Performance**: Fast, works offline, limited to predefined patterns

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Choose your AI provider:
  - **Enhanced Local**: JSON configuration file (recommended)
  - **OpenAI**: OpenAI API key
  - **Ollama**: Ollama installation with models
  - **Local**: No additional requirements

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd myllm
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables (see Configuration section)

4. Run the application:
```bash
go run main.go
```

### Configuration

#### Environment Variables

```bash
# AI Provider Configuration
AI_PROVIDER=enhanced_local          # Options: "openai", "ollama", "local", "enhanced_local"
AI_MODEL=                           # Model name (provider-specific)
AI_TEMPERATURE=0.1                  # Generation temperature
AI_MAX_TOKENS=1000                  # Maximum tokens to generate
AI_BASE_URL=http://localhost:11434  # Base URL for local providers

# Enhanced Local AI Configuration
INTENT_CONFIG_PATH=configs/personal_assistant.json  # Path to intent config file

# OpenAI Configuration (for AI_PROVIDER=openai)
OPENAI_API_KEY=your-key             # Required for OpenAI

# Server Configuration
PORT=8080                           # Server port
```

#### Provider-Specific Setup

**Enhanced Local AI Setup (Recommended):**
```bash
# Use the provided configuration or create your own
export AI_PROVIDER=enhanced_local
export INTENT_CONFIG_PATH=configs/personal_assistant.json
go run main.go
```

**OpenAI Setup:**
```bash
export AI_PROVIDER=openai
export OPENAI_API_KEY=your-openai-api-key
export AI_MODEL=gpt-3.5-turbo
```

**Ollama Setup:**
```bash
# Install Ollama (https://ollama.ai)
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama2

# Configure the application
export AI_PROVIDER=ollama
export AI_MODEL=llama2
export AI_BASE_URL=http://localhost:11434
```

**Local AI Setup:**
```bash
export AI_PROVIDER=local
# No additional configuration needed
```

### Usage Examples

#### Enhanced Local AI (High Accuracy)

**Input:**
```bash
curl -X POST http://localhost:8080/api/v1/intent \
  -H "Content-Type: application/json" \
  -d '{"text": "create a new contact named John Smith with email john@example.com"}'
```

**Output:**
```json
{
  "success": true,
  "intent": {
    "task": "CreateContact",
    "vars": {
      "name": "John Smith",
      "email": "john@example.com",
      "phone": "",
      "confidence": 0.85
    }
  }
}
```

#### Complex Natural Language (AI-powered)

**Input:**
```bash
curl -X POST http://localhost:8080/api/v1/intent \
  -H "Content-Type: application/json" \
  -d '{"text": "I need to add John Smith to my contacts, his email is john.smith@company.com and phone is 555-1234"}'
```

**Output:**
```json
{
  "success": true,
  "intent": {
    "task": "CreateContact",
    "vars": {
      "name": "John Smith",
      "email": "john.smith@company.com",
      "phone": "555-1234",
      "confidence": 0.92
    }
  }
}
```

## API Reference

### POST /api/v1/intent

Extracts intent and variables from natural language text.

**Request Body:**
```json
{
  "text": "string"
}
```

**Response:**
```json
{
  "success": boolean,
  "intent": {
    "task": "string",
    "vars": {
      "name": "string",
      "email": "string",
      "phone": "string",
      "confidence": "number"
    }
  },
  "error": "string"  // Only present when success is false
}
```

### GET /api/v1/health

Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "service": "intent-recognition-api"
}
```

## Enhanced Local AI Configuration

The Enhanced Local AI provider uses JSON configuration files to define intents, entities, and patterns. This allows for highly accurate, domain-specific intent recognition.

### Configuration Structure

```json
{
  "domain": "personal_assistant",
  "version": "1.0.0",
  "intents": {
    "CreateContact": {
      "description": "Create a new contact",
      "keywords": ["create", "add", "new"],
      "phrases": ["create a new contact", "add contact"],
      "regex": ["(?i)create\\s+(?:a\\s+)?(?:new\\s+)?contact"],
      "priority": 10,
      "variables": ["name", "email", "phone"],
      "examples": ["create a new contact named John"]
    }
  },
  "entities": {
    "name": {
      "type": "name",
      "regex": ["(?i)(?:named\\s+)([A-Z][a-z]+(?:\\s+[A-Z][a-z]+)*)"],
      "keywords": ["named", "name", "called"]
    }
  },
  "synonyms": {
    "create": ["add", "new", "make"],
    "find": ["search", "look", "get"]
  },
  "confidence": {
    "CreateContact": 0.7
  }
}
```

### Creating Custom Configurations

1. **Define Intents**: List all possible intents for your domain
2. **Add Keywords**: Primary words that indicate each intent
3. **Include Phrases**: Common ways users express each intent
4. **Create Regex**: Specific patterns for precise matching
5. **Set Priorities**: Higher priority intents are matched first
6. **Define Entities**: What data to extract (names, emails, etc.)
7. **Add Synonyms**: Alternative words for better matching
8. **Set Confidence**: Thresholds for each intent

### Example Domains

- **Personal Assistant**: Contacts, tasks, events, notes, weather, time
- **Customer Support**: Tickets, complaints, requests, status checks
- **E-commerce**: Orders, products, shipping, returns, payments
- **Healthcare**: Appointments, symptoms, medications, records

## Provider Comparison

| Feature | Enhanced Local | OpenAI | Ollama | Local AI |
|---------|---------------|--------|--------|----------|
| **Accuracy** | 90%+ (configured) | High | Good | Limited |
| **Speed** | Very Fast | Fast | Medium | Very Fast |
| **Offline** | Yes | No | Yes | Yes |
| **Privacy** | Local | Cloud | Local | Local |
| **Setup** | JSON Config | API Key | Ollama + Models | None |
| **Cost** | Free | Per token | Free | Free |
| **Customization** | Very High | Limited | High | High |
| **Domain Specific** | Excellent | Good | Good | Limited |

## Testing

### Run Enhanced Local AI Tests

```bash
# Start server with enhanced local provider
export AI_PROVIDER=enhanced_local
export INTENT_CONFIG_PATH=configs/personal_assistant.json
go run main.go

# In another terminal, run tests
./test_enhanced_local.sh
```

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

## Performance Considerations

- **Enhanced Local AI**: Pre-compiled patterns, fuzzy matching, confidence scoring
- **Pattern Matching**: Fast regex-based extraction for common patterns
- **AI Fallback**: Only uses AI when patterns don't match
- **Provider Selection**: Automatically falls back to available providers
- **Caching**: Consider implementing Redis caching for frequently requested patterns
- **Rate Limiting**: Implement rate limiting for cloud AI API calls

## Security

- Input validation and sanitization
- Environment variable management
- Error message sanitization
- Request timeout handling
- Local AI options for sensitive data
- Configurable confidence thresholds

## Deployment

### Docker

```bash
docker build -t intent-recognition-api .
docker run -p 8080:8080 -e AI_PROVIDER=enhanced_local -v $(pwd)/configs:/app/configs intent-recognition-api
```

### Production Considerations

- Use HTTPS in production
- Implement proper logging and monitoring
- Set up health checks and metrics
- Configure proper timeouts and rate limits
- Use environment-specific configurations
- Consider local AI for privacy-sensitive deployments
- Use enhanced local AI for domain-specific accuracy

## Troubleshooting

### Common Issues

1. **Enhanced Local AI errors**: Check JSON configuration syntax
2. **OpenAI API errors**: Check your API key and billing
3. **Ollama connection errors**: Ensure Ollama is running and models are downloaded
4. **Low accuracy**: Add more patterns and examples to your configuration

### Provider Fallback

The application automatically falls back to available providers:
1. Try configured provider
2. Try other available providers
3. Fall back to basic local rule-based extraction

### Configuration Tips

- **Start Simple**: Begin with basic keywords and phrases
- **Add Examples**: Include real user examples for better training
- **Use Regex**: Add specific patterns for precise matching
- **Set Priorities**: Order intents by specificity
- **Test Thoroughly**: Validate with real user inputs
- **Iterate**: Continuously improve based on results

## Contributing

1. Follow the established code quality standards
2. Write tests for new features
3. Update documentation
4. Ensure all tests pass before submitting
5. Add configuration examples for new domains

## License

[Your License Here] 