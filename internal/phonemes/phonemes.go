// Package phonemes converts lyrics to SynthV phonemes using adaptive
// multi-language selection (Mandarin, Cantonese, Japanese, English, Korean).
package phonemes

import "strings"

// Result holds the conversion result for a single note.
type Result struct {
	Language string // target SynthV language (e.g., "mandarin", "english")
	Phoneset string // phoneset for the language (e.g., "xsampa", "romaji", "arpabet")
	Phonemes string // space-separated phoneme string
}

// Converter converts lyrics to phonemes for a specific source language.
type Converter struct {
	selectLang func(string) string
	tables     map[string]*phoneTable
	normalize  func(string) string
	skip       func(rune) bool
}

// New creates a Converter for the given source language.
// Supported: "estonian", "karelian".
func New(lang string) *Converter {
	switch strings.ToLower(lang) {
	case "estonian":
		return newEstonian()
	case "karelian":
		return newKarelian()
	default:
		return nil
	}
}

// Convert converts a lyrics word to phonemes.
// Returns empty Result for special markers ("-", "+", "sil", "br", etc.).
func (c *Converter) Convert(word string) Result {
	lower := strings.ToLower(word)

	switch lower {
	case "-", "+", "sil", "br", "sp", "ap", "":
		return Result{}
	}

	lang := c.selectLang(lower)
	table := c.tables[lang]
	if table == nil {
		return Result{Language: lang}
	}

	w := lower
	if c.normalize != nil {
		w = c.normalize(w)
	}

	return Result{
		Language: lang,
		Phoneset: phonesetForLanguage(lang),
		Phonemes: table.convert(w, c.skip),
	}
}

type phoneTable struct {
	// Digraphs are checked first (two-char sequences like "sh", "ng", "ts").
	// Also used for special geminates that don't simply double (Korean "ss"→"s_t").
	digraphs map[string][]string
	// Singles maps individual runes to their phoneme.
	// Geminates (doubled chars) automatically double the single phoneme.
	singles map[rune]string
}

func (t *phoneTable) convert(word string, skip func(rune) bool) string {
	runes := []rune(word)
	var phonemes []string

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		if skip != nil && skip(ch) {
			continue
		}

		// Check digraphs (two-char sequences).
		if i+1 < len(runes) {
			pair := string(runes[i : i+2])
			if ph, ok := t.digraphs[pair]; ok {
				phonemes = append(phonemes, ph...)
				i++
				continue
			}
		}

		ph, ok := t.singles[ch]
		if !ok {
			continue
		}

		// Check geminate (doubled character) → double the single phoneme.
		if i+1 < len(runes) && runes[i+1] == ch {
			phonemes = append(phonemes, ph, ph)
			i++
			continue
		}

		phonemes = append(phonemes, ph)
	}

	return strings.Join(phonemes, " ")
}

func phonesetForLanguage(lang string) string {
	switch lang {
	case "japanese":
		return "romaji"
	case "english":
		return "arpabet"
	default:
		return "xsampa"
	}
}

func hasOnlyBasicVowels(word string, special string) bool {
	for _, r := range word {
		if strings.ContainsRune(special, r) {
			return false
		}
	}
	return true
}
