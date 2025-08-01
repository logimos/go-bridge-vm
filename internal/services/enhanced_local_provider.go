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
		fmt.Printf("Loading intent configuration from: %s\n", configPath)
		config, err = models.LoadIntentConfig(configPath)
		if err != nil {
			fmt.Printf("Failed to load config from %s: %v\n", configPath, err)
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
		fmt.Printf("Successfully loaded configuration with domain: %s\n", config.Domain)
	} else {
		fmt.Printf("No config path provided, using default configuration\n")
		config = models.GetDefaultConfig()
		fmt.Printf("Using default configuration with domain: %s\n", config.Domain)
	}

	// Log available intents
	fmt.Printf("Available intents (%d):\n", len(config.Intents))
	for intentName, intent := range config.Intents {
		fmt.Printf("  - %s: %s (priority: %d, required: %v)\n",
			intentName, intent.Description, intent.Priority, intent.Required)
	}

	// Log available entities
	fmt.Printf("Available entities (%d):\n", len(config.Entities))
	for entityName, entity := range config.Entities {
		fmt.Printf("  - %s: %s\n", entityName, entity.Description)
	}

	// Compile patterns for performance
	compiled, err := compileConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to compile config: %w", err)
	}

	fmt.Printf("Configuration compilation completed successfully\n")

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
		// Use a more natural intent name for the question
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "What should I call this " + intentDisplayName + "?"
	case "name":
		return "What's the name?"
	case "email":
		return "What's the email address?"
	case "phone":
		return "What's the phone number?"
	case "date":
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "When should this " + intentDisplayName + " be scheduled?"
	case "time":
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "What time should this " + intentDisplayName + " be?"
	case "duration":
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "How long should this " + intentDisplayName + " last?"
	case "location":
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "Where should this " + intentDisplayName + " take place?"
	case "description":
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "Can you provide more details about this " + intentDisplayName + "?"
	case "priority":
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "What priority should this " + intentDisplayName + " have?"
	default:
		intentDisplayName := p.getIntentDisplayName(intentName)
		return "What " + field + " should I use for this " + intentDisplayName + "?"
	}
}

// getIntentDisplayName converts intent names to user-friendly display names
func (p *EnhancedLocalProvider) getIntentDisplayName(intentName string) string {
	switch intentName {
	case "CreateEvent":
		return "event"
	case "CreateTask":
		return "task"
	case "CreateContact":
		return "contact"
	case "CreateNote":
		return "note"
	case "Weather":
		return "weather"
	case "Time":
		return "time"
	case "Calculator":
		return "calculation"
	default:
		// Convert camelCase to lowercase with spaces
		result := ""
		for i, char := range intentName {
			if i > 0 && unicode.IsUpper(char) {
				result += " "
			}
			result += string(unicode.ToLower(char))
		}
		return result
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

	// Extract name first (can be quoted)
	for entityName, entity := range p.config.Entities {
		if entityName == "name" {
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
	}

	// Extract title (can be quoted, but don't override name)
	for entityName, entity := range p.config.Entities {
		if entityName == "title" {
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
	}

	// Extract other entities
	for entityName, entity := range p.config.Entities {
		if entityName == "name" || entityName == "title" {
			continue // Already processed
		}

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

	// Only use keyword-based extraction for specific entity types that have clear patterns
	switch entityName {
	case "name":
		// First try to extract names in quotes (most reliable)
		quotePattern := regexp.MustCompile(`"([^"]+)"`)
		matches := quotePattern.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}

		// Look for name patterns like "named John", "contact Alice", "for Bob"
		for i, word := range words {
			wordLower := strings.ToLower(word)
			if wordLower == "named" || wordLower == "name" || wordLower == "contact" || wordLower == "person" {
				if i+1 < len(words) {
					nextWord := strings.Trim(words[i+1], ".,!?;:")
					// Check if it looks like a name (starts with capital letter, not a common word)
					if len(nextWord) > 0 && unicode.IsUpper(rune(nextWord[0])) && !p.isStopWord(strings.ToLower(nextWord)) {
						// Only take the first word if it's followed by "with" or other indicators
						if i+2 < len(words) && strings.ToLower(words[i+2]) == "with" {
							return nextWord
						}
						// Check if next word is a stop word or indicator
						if i+2 < len(words) {
							nextNextWord := strings.ToLower(words[i+2])
							if nextNextWord == "with" || nextNextWord == "email" || nextNextWord == "phone" {
								return nextWord
							}
						}
						return nextWord
					}
				}
			}
		}

	case "email":
		// Look for email patterns like "email alice@example.com"
		for i, word := range words {
			wordLower := strings.ToLower(word)
			if wordLower == "email" || wordLower == "e-mail" || wordLower == "mail" {
				if i+1 < len(words) {
					nextWord := strings.Trim(words[i+1], ".,!?;:")
					// Check if it looks like an email
					if strings.Contains(nextWord, "@") && strings.Contains(nextWord, ".") {
						return nextWord
					}
				}
			}
		}

	case "phone":
		// Look for phone patterns like "phone 555-123-4567"
		for i, word := range words {
			wordLower := strings.ToLower(word)
			if wordLower == "phone" || wordLower == "mobile" || wordLower == "cell" {
				if i+1 < len(words) {
					nextWord := strings.Trim(words[i+1], ".,!?;:")
					// Check if it looks like a phone number (contains digits and possibly dashes/parentheses)
					if strings.ContainsAny(nextWord, "0123456789") && (strings.Contains(nextWord, "-") || strings.Contains(nextWord, "(") || len(nextWord) >= 10) {
						return nextWord
					}
				}
			}
		}

	case "date":
		// Look for date patterns like "tomorrow", "today", "next week"
		for _, word := range words {
			wordLower := strings.ToLower(strings.Trim(word, ".,!?;:"))
			if wordLower == "today" || wordLower == "tomorrow" || wordLower == "yesterday" {
				return wordLower
			}
		}

	case "time":
		// Look for time patterns like "at 2pm", "at 3:30"
		for i, word := range words {
			wordLower := strings.ToLower(word)
			if wordLower == "at" {
				if i+1 < len(words) {
					nextWord := strings.Trim(words[i+1], ".,!?;:")
					// Check if it looks like a time (contains digits and possibly : or am/pm)
					if strings.ContainsAny(nextWord, "0123456789") && (strings.Contains(nextWord, ":") || strings.Contains(strings.ToLower(nextWord), "am") || strings.Contains(strings.ToLower(nextWord), "pm")) {
						return nextWord
					}
				}
			}
		}

	case "location":
		// Look for location patterns like "in New York", "at the office"
		for i, word := range words {
			wordLower := strings.ToLower(word)
			if wordLower == "in" || wordLower == "at" {
				if i+1 < len(words) {
					nextWord := strings.Trim(words[i+1], ".,!?;:")
					// Check if it looks like a location (starts with capital letter, not a common word)
					if len(nextWord) > 0 && unicode.IsUpper(rune(nextWord[0])) && !p.isStopWord(strings.ToLower(nextWord)) {
						return nextWord
					}
				}
			}
		}

	case "title":
		// First try to extract titles in quotes (most reliable)
		quotePattern := regexp.MustCompile(`"([^"]+)"`)
		matches := quotePattern.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}

		// Look for title patterns like "called buy groceries", "for team meeting"
		for i, word := range words {
			wordLower := strings.ToLower(word)
			if wordLower == "called" || wordLower == "titled" || wordLower == "named" || wordLower == "for" {
				// Collect all words after the keyword until we hit a stop word or punctuation
				var titleWords []string
				for j := i + 1; j < len(words); j++ {
					nextWord := strings.Trim(words[j], ".,!?;:")
					if len(nextWord) == 0 {
						continue
					}
					nextWordLower := strings.ToLower(nextWord)

					// Stop if we hit a stop word or time/date indicator
					if p.isStopWord(nextWordLower) ||
						nextWordLower == "tomorrow" ||
						nextWordLower == "today" ||
						nextWordLower == "yesterday" ||
						nextWordLower == "at" ||
						nextWordLower == "with" ||
						nextWordLower == "email" ||
						nextWordLower == "phone" ||
						nextWordLower == "on" ||
						nextWordLower == "in" ||
						nextWordLower == "to" {
						break
					}
					titleWords = append(titleWords, nextWord)
				}
				if len(titleWords) > 0 {
					return strings.Join(titleWords, " ")
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
