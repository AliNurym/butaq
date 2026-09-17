package locales

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed *.json
var embeddedLocales embed.FS

// Locale represents a declarative syntax projection for a human language.
type Locale struct {
	Code     string            `json:"code"`
	Name     string            `json:"name"`
	Keywords map[string]string `json:"keywords"` // localized keyword -> internal token name ("болсын" -> "VAR")
	Builtins map[string]string `json:"builtins"` // localized builtin -> canonical name ("мәтін_ұзындығы" -> "str_len")
	Messages map[string]string `json:"messages,omitempty"` // localized compiler diagnostic messages

	// Computed reverse lookups
	TokenToKeyword map[string]string `json:"-"` // token name -> preferred localized keyword ("VAR" -> "let")
	BuiltinToLocal map[string]string `json:"-"` // canonical builtin -> localized name ("str_len" -> "str_length")
}

// Message returns localized message for key, or fallback if not found.
func (l *Locale) Message(key string, fallback string) string {
	if l != nil && l.Messages != nil {
		if msg, ok := l.Messages[key]; ok && msg != "" {
			return msg
		}
	}
	return fallback
}

var (
	registryLock sync.RWMutex
	registry     = make(map[string]*Locale)
	initialized  = false
)

func init() {
	Init()
}

// Init loads all embedded and directory locales.
func Init() {
	registryLock.Lock()
	defer registryLock.Unlock()

	if initialized {
		return
	}

	// 1. Load embedded JSON files
	entries, err := embeddedLocales.ReadDir(".")
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				data, err := embeddedLocales.ReadFile(entry.Name())
				if err == nil {
					var loc Locale
					if err := json.Unmarshal(data, &loc); err == nil {
						initReverseMaps(&loc)
						registry[strings.ToLower(loc.Code)] = &loc
					}
				}
			}
		}
	}

	// 2. Load from local filesystem "locales" directory if present (allows user to drop in new languages!)
	if dirEntries, err := os.ReadDir("locales"); err == nil {
		for _, entry := range dirEntries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				path := filepath.Join("locales", entry.Name())
				data, err := os.ReadFile(path)
				if err == nil {
					var loc Locale
					if err := json.Unmarshal(data, &loc); err == nil {
						initReverseMaps(&loc)
						registry[strings.ToLower(loc.Code)] = &loc
					}
				}
			}
		}
	}

	initialized = true
}

func initReverseMaps(loc *Locale) {
	loc.TokenToKeyword = make(map[string]string)
	loc.BuiltinToLocal = make(map[string]string)

	for word, tokName := range loc.Keywords {
		// Prefer the shortest or first defined keyword for transpilation output
		if existing, ok := loc.TokenToKeyword[tokName]; !ok || len(word) < len(existing) {
			loc.TokenToKeyword[tokName] = word
		}
	}

	for localFunc, canonFunc := range loc.Builtins {
		if existing, ok := loc.BuiltinToLocal[canonFunc]; !ok || len(localFunc) < len(existing) {
			loc.BuiltinToLocal[canonFunc] = localFunc
		}
	}
}

// Get returns the locale by code (e.g. "kk", "en", "it", "ru").
func Get(code string) (*Locale, bool) {
	registryLock.RLock()
	defer registryLock.RUnlock()

	loc, ok := registry[strings.ToLower(strings.TrimSpace(code))]
	return loc, ok
}

// Default returns the default reference locale (Kazakh "kk").
func Default() *Locale {
	if loc, ok := Get("kk"); ok {
		return loc
	}
	// Fallback to any available
	registryLock.RLock()
	defer registryLock.RUnlock()
	for _, loc := range registry {
		return loc
	}
	return &Locale{Code: "kk", Name: "Қазақша", Keywords: map[string]string{}, Builtins: map[string]string{}}
}

// Available returns a list of all registered locales.
func Available() []*Locale {
	registryLock.RLock()
	defer registryLock.RUnlock()

	res := make([]*Locale, 0, len(registry))
	for _, loc := range registry {
		res = append(res, loc)
	}
	return res
}

// AutoDetect examines source code tokens and returns the most likely matching locale.
func AutoDetect(source string) *Locale {
	registryLock.RLock()
	defer registryLock.RUnlock()

	words := strings.FieldsFunc(source, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '(' || r == ')' || r == '{' || r == '}' || r == '[' || r == ']' || r == ',' || r == '"'
	})

	scores := make(map[string]int)
	for _, w := range words {
		for code, loc := range registry {
			if _, ok := loc.Keywords[w]; ok {
				scores[code]++
			}
			if _, ok := loc.Builtins[w]; ok {
				scores[code]++
			}
		}
	}

	bestCode := "kk"
	maxScore := 0
	for code, score := range scores {
		if score > maxScore {
			maxScore = score
			bestCode = code
		}
	}

	if loc, ok := registry[bestCode]; ok {
		return loc
	}
	return Default()
}

// KeywordForToken returns the localized keyword string for an internal token name.
func (l *Locale) KeywordForToken(tokenName string) string {
	if l == nil {
		return tokenName
	}
	if kw, ok := l.TokenToKeyword[tokenName]; ok {
		return kw
	}
	return tokenName
}

// LocalNameForBuiltin returns the localized name for a canonical builtin function.
func (l *Locale) LocalNameForBuiltin(canonical string) string {
	if l == nil {
		return canonical
	}
	if name, ok := l.BuiltinToLocal[canonical]; ok {
		return name
	}
	return canonical
}

// CanonicalForLocalBuiltin returns the canonical builtin name for a localized function call.
func (l *Locale) CanonicalForLocalBuiltin(local string) string {
	if l == nil {
		return local
	}
	if canon, ok := l.Builtins[local]; ok {
		return canon
	}
	return local
}

// RegisterFromFile allows loading a custom locale JSON file dynamically at runtime.
func RegisterFromFile(filePath string) (*Locale, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locale file: %w", err)
	}

	var loc Locale
	if err := json.Unmarshal(data, &loc); err != nil {
		return nil, fmt.Errorf("invalid locale JSON: %w", err)
	}

	initReverseMaps(&loc)

	registryLock.Lock()
	registry[strings.ToLower(loc.Code)] = &loc
	registryLock.Unlock()

	return &loc, nil
}

// CanonicalizeBuiltin resolves any localized builtin name to its canonical Kazakh compiler representation.
func CanonicalizeBuiltin(name string) string {
	registryLock.RLock()
	defer registryLock.RUnlock()

	// 1. If it's already a Kazakh builtin name, return as-is
	if kk, ok := registry["kk"]; ok {
		if _, exists := kk.Builtins[name]; exists {
			return name
		}
	}

	// 2. Search if name exists in any locale's builtins
	canonical := ""
	for _, loc := range registry {
		if canon, exists := loc.Builtins[name]; exists {
			canonical = canon
			break
		}
	}

	if canonical != "" {
		if kk, ok := registry["kk"]; ok {
			if kkName, exists := kk.BuiltinToLocal[canonical]; exists {
				return kkName
			}
		}
	}

	return name
}
