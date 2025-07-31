package services

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"myllm/internal/models"
)

// EnhancedLocalProvider implements AIProvider with configurable intent recognition
type EnhancedLocalProvider struct {
	config     *models.IntentConfig
	compiled   *CompiledConfig
	configPath string
}

// CompiledConfig holds pre-compiled patterns for performance
type CompiledConfig struct {
	IntentRegexes map[string][]*regexp.Regexp
	EntityRegexes map[string][]*regexp.Regexp
	KeywordMap    map[string][]string
	PhraseMap     map[string][]string
	SynonymMap    map[string]string
}

// NewEnhancedLocalProvider creates a new enhanced local AI provider
func NewEnhancedLocalProvider(configPath string) (AIProvider, error) {
	var config *models.IntentConfig
	var err error

	// Try to load from file, fallback to default
	if configPath != "" {
		config, err = models.LoadIntentConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	} else {
		config = models.GetDefaultConfig()
	}

	// Compile patterns for performance
	compiled, err := compileConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to compile config: %w", err)
	}

	return &EnhancedLocalProvider{
		config:     config,
		compiled:   compiled,
		configPath: configPath,
	}, nil
}

// compileConfig pre-compiles all regex patterns for performance
func compileConfig(config *models.IntentConfig) (*CompiledConfig, error) {
	compiled := &CompiledConfig{
		IntentRegexes: make(map[string][]*regexp.Regexp),
		EntityRegexes: make(map[string][]*regexp.Regexp),
		KeywordMap:    make(map[string][]string),
		PhraseMap:     make(map[string][]string),
		SynonymMap:    make(map[string]string),
	}

	// Compile intent regexes
	for intentName, intent := range config.Intents {
		var regexes []*regexp.Regexp
		for _, pattern := range intent.Regex {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid regex for intent %s: %w", intentName, err)
			}
			regexes = append(regexes, re)
		}
		compiled.IntentRegexes[intentName] = regexes
		compiled.KeywordMap[intentName] = intent.Keywords
		compiled.PhraseMap[intentName] = intent.Phrases
	}

	// Compile entity regexes
	for entityName, entity := range config.Entities {
		var regexes []*regexp.Regexp
		for _, pattern := range entity.Regex {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid regex for entity %s: %w", entityName, err)
			}
			regexes = append(regexes, re)
		}
		compiled.EntityRegexes[entityName] = regexes
	}

	// Build synonym map
	for word, synonyms := range config.Synonyms {
		for _, synonym := range synonyms {
			compiled.SynonymMap[synonym] = word
		}
	}

	return compiled, nil
}

// ExtractIntent extracts intent using enhanced local processing
func (p *EnhancedLocalProvider) ExtractIntent(ctx context.Context, text string) (*models.Intent, error) {
	normalizedText := p.normalizeText(text)

	// Get intent with confidence score
	intentResult := p.classifyIntent(normalizedText)

	// Extract entities
	entities := p.extractEntities(text)

	// Build the intent structure
	result := &models.Intent{
		Task: intentResult.Intent,
		Vars: make(map[string]interface{}),
	}

	// Map extracted entities to variables
	for entityType, value := range entities {
		result.Vars[entityType] = value
	}

	// Add confidence score
	result.Vars["confidence"] = intentResult.Confidence

	// Check for missing required fields and generate follow-up questions
	if intentResult.Intent != "UNKNOWN" {
		p.addMissingFieldsAndFollowUp(result, intentResult.Intent)
	}

	return result, nil
}

// addMissingFieldsAndFollowUp checks for missing required fields and adds follow-up questions
func (p *EnhancedLocalProvider) addMissingFieldsAndFollowUp(intent *models.Intent, intentName string) {
	intentPattern, exists := p.config.Intents[intentName]
	if !exists {
		return
	}

	var missing []string
	var followUp []string

	// Check which required fields are missing
	for _, requiredField := range intentPattern.Required {
		if value, exists := intent.Vars[requiredField]; !exists || value == "" {
			missing = append(missing, requiredField)
		}
	}

	// Generate follow-up questions for missing fields
	for _, field := range missing {
		question := p.generateFollowUpQuestion(intentName, field, intentPattern.FollowUp)
		if question != "" {
			followUp = append(followUp, question)
		}
	}

	// Update intent with missing fields and follow-up questions
	intent.Missing = missing
	intent.FollowUp = followUp
	intent.IsComplete = len(missing) == 0
}

// generateFollowUpQuestion generates a follow-up question for a missing field
func (p *EnhancedLocalProvider) generateFollowUpQuestion(intentName, field string, customFollowUp []string) string {
	// Try to use custom follow-up questions first
	for _, question := range customFollowUp {
		if strings.Contains(strings.ToLower(question), strings.ToLower(field)) {
			return question
		}
	}

	// Generate default questions based on field type
	switch field {
	case "title":
		return "What should I call this " + strings.ToLower(intentName) + "?"
	case "name":
		return "What's the name?"
	case "email":
		return "What's the email address?"
	case "phone":
		return "What's the phone number?"
	case "date":
		return "When should this " + strings.ToLower(intentName) + " be scheduled?"
	case "time":
		return "What time should this " + strings.ToLower(intentName) + " be?"
	case "duration":
		return "How long should this " + strings.ToLower(intentName) + " last?"
	case "location":
		return "Where should this " + strings.ToLower(intentName) + " take place?"
	case "description":
		return "Can you provide more details about this " + strings.ToLower(intentName) + "?"
	case "priority":
		return "What priority should this " + strings.ToLower(intentName) + " have?"
	default:
		return "What " + field + " should I use for this " + strings.ToLower(intentName) + "?"
	}
}

// IntentResult holds intent classification with confidence
type IntentResult struct {
	Intent     string
	Confidence float64
}

// classifyIntent determines the intent with confidence scoring
func (p *EnhancedLocalProvider) classifyIntent(text string) IntentResult {
	var bestIntent string = "UNKNOWN"
	var bestScore float64 = 0.0

	// Score each intent
	intentScores := make(map[string]float64)

	for intentName, intent := range p.config.Intents {
		score := p.calculateIntentScore(text, intentName, intent)
		intentScores[intentName] = score

		// Apply priority boost
		priorityBoost := float64(intent.Priority) * 0.1
		score += priorityBoost

		if score > bestScore {
			bestScore = score
			bestIntent = intentName
		}
	}

	// Check confidence threshold
	threshold := p.config.Confidence[bestIntent]
	if threshold == 0 {
		threshold = 0.5 // Default threshold
	}

	if bestScore < threshold {
		bestIntent = "UNKNOWN"
		bestScore = 0.0
	}

	return IntentResult{
		Intent:     bestIntent,
		Confidence: math.Min(bestScore, 1.0),
	}
}

// calculateIntentScore calculates a confidence score for an intent
func (p *EnhancedLocalProvider) calculateIntentScore(text, intentName string, intent models.IntentPattern) float64 {
	score := 0.0

	// 1. Regex matching (highest weight)
	for _, re := range p.compiled.IntentRegexes[intentName] {
		if re.MatchString(text) {
			score += 0.8
			break
		}
	}

	// 2. Exact phrase matching (high weight)
	textLower := strings.ToLower(text)
	for _, phrase := range p.compiled.PhraseMap[intentName] {
		if strings.Contains(textLower, strings.ToLower(phrase)) {
			score += 0.6
			break
		}
	}

	// 3. Keyword matching with fuzzy scoring
	keywords := p.compiled.KeywordMap[intentName]
	keywordScore := 0.0
	matchedKeywords := 0

	for _, keyword := range keywords {
		// Exact match
		if strings.Contains(textLower, strings.ToLower(keyword)) {
			keywordScore += 0.4
			matchedKeywords++
		} else {
			// Fuzzy match using synonym expansion
			synonyms := p.getSynonyms(keyword)
			for _, synonym := range synonyms {
				if strings.Contains(textLower, strings.ToLower(synonym)) {
					keywordScore += 0.3
					matchedKeywords++
					break
				}
			}
		}
	}

	// Normalize keyword score
	if len(keywords) > 0 {
		keywordScore = keywordScore / float64(len(keywords))
	}

	score += keywordScore

	// 4. Word overlap scoring
	textWords := p.tokenize(text)
	intentWords := p.getIntentWords(intent)
	overlap := p.calculateWordOverlap(textWords, intentWords)
	score += overlap * 0.2

	// 5. Length bonus (longer, more specific queries get higher scores)
	if len(text) > 20 {
		score += 0.1
	}

	return score
}

// extractEntities extracts entities using configurable patterns
func (p *EnhancedLocalProvider) extractEntities(text string) map[string]string {
	entities := make(map[string]string)

	for entityName, entity := range p.config.Entities {
		// Try regex patterns first
		for _, re := range p.compiled.EntityRegexes[entityName] {
			matches := re.FindStringSubmatch(text)
			if len(matches) > 1 {
				entities[entityName] = matches[1]
				break
			}
		}

		// If no regex match, try keyword-based extraction
		if entities[entityName] == "" {
			value := p.extractEntityByKeywords(text, entityName, entity)
			if value != "" {
				entities[entityName] = value
			}
		}
	}

	return entities
}

// extractEntityByKeywords extracts entities using keyword context
func (p *EnhancedLocalProvider) extractEntityByKeywords(text, entityName string, entity models.EntityPattern) string {
	words := strings.Fields(text)

	for i, word := range words {
		wordLower := strings.ToLower(word)

		// Check if this word is a keyword for the entity
		for _, keyword := range entity.Keywords {
			if strings.Contains(wordLower, strings.ToLower(keyword)) {
				// Look for the entity value after the keyword
				if i+1 < len(words) {
					nextWord := words[i+1]
					// Clean the word
					nextWord = strings.Trim(nextWord, ".,!?;:")
					if len(nextWord) > 0 {
						return nextWord
					}
				}
			}
		}
	}

	return ""
}

// normalizeText performs advanced text normalization
func (p *EnhancedLocalProvider) normalizeText(text string) string {
	// Convert to lowercase
	normalized := strings.ToLower(text)

	// Remove extra whitespace
	normalized = strings.Join(strings.Fields(normalized), " ")

	// Remove punctuation (but keep some for context)
	normalized = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) && r != '@' && r != '.' && r != '_' && r != '-' {
			return ' '
		}
		return r
	}, normalized)

	// Clean up multiple spaces
	normalized = strings.Join(strings.Fields(normalized), " ")

	return normalized
}

// tokenize splits text into meaningful tokens
func (p *EnhancedLocalProvider) tokenize(text string) []string {
	// Simple tokenization - can be enhanced with NLP libraries
	words := strings.Fields(strings.ToLower(text))
	var tokens []string

	for _, word := range words {
		// Remove common stop words
		if !p.isStopWord(word) {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

// isStopWord checks if a word is a common stop word
func (p *EnhancedLocalProvider) isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "up": true, "about": true,
		"into": true, "through": true, "during": true, "before": true, "after": true,
		"above": true, "below": true, "between": true, "among": true,
		"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
		"have": true, "has": true, "had": true, "do": true, "does": true, "did": true,
		"will": true, "would": true, "could": true, "should": true, "may": true, "might": true,
	}

	return stopWords[word]
}

// getSynonyms returns synonyms for a word
func (p *EnhancedLocalProvider) getSynonyms(word string) []string {
	if synonyms, exists := p.config.Synonyms[word]; exists {
		return synonyms
	}
	return []string{}
}

// getIntentWords gets all words associated with an intent
func (p *EnhancedLocalProvider) getIntentWords(intent models.IntentPattern) []string {
	var words []string
	words = append(words, intent.Keywords...)

	// Add words from phrases
	for _, phrase := range intent.Phrases {
		phraseWords := strings.Fields(strings.ToLower(phrase))
		words = append(words, phraseWords...)
	}

	return words
}

// calculateWordOverlap calculates word overlap between text and intent
func (p *EnhancedLocalProvider) calculateWordOverlap(textWords, intentWords []string) float64 {
	if len(intentWords) == 0 {
		return 0.0
	}

	textSet := make(map[string]bool)
	for _, word := range textWords {
		textSet[word] = true
	}

	overlap := 0
	for _, word := range intentWords {
		if textSet[word] {
			overlap++
		}
	}

	return float64(overlap) / float64(len(intentWords))
}

// Name returns the provider name
func (p *EnhancedLocalProvider) Name() string {
	if p.configPath != "" {
		return fmt.Sprintf("Enhanced Local AI (%s)", p.config.Domain)
	}
	return "Enhanced Local AI (Default)"
}

// IsAvailable checks if enhanced local AI is available
func (p *EnhancedLocalProvider) IsAvailable() bool {
	return true // Always available
}

// GetConfig returns the current configuration
func (p *EnhancedLocalProvider) GetConfig() *models.IntentConfig {
	return p.config
}
