// Package phonemes converts lyrics to SynthV phonemes using adaptive
// multi-language selection (Mandarin, Cantonese, Spanish).
package phonemes

import (
	"slices"
	"strings"
)

// Result holds the conversion result for a single note.
type Result struct {
	Language string // target SynthV language (e.g., "mandarin", "english")
	Phoneset string // phoneset for the language (e.g., "xsampa", "romaji", "arpabet")
	Phonemes string // space-separated phoneme string
}

// Converter converts lyrics to phonemes for a specific source language.
type Converter struct {
	// selectLangs returns the target languages that can render a word, best
	// first. Languages after the first fit the word just as well, so the
	// voice's own language can be preferred among them.
	selectLangs func(string) []string
	tables      map[string]*phoneTable
	normalize   func(string) string
	skip        func(rune) bool
	vowels      string            // source-language vowels; doubled vowels emit a single phoneme
	words       map[string]Result // whole-word overrides checked before table conversion

	prefer       string // voice language, preferred among equally good candidates
	preferAlways bool   // use the voice language even when it fits worse
}

// Prefer sets the voice's own language. It is used when several languages
// render a word equally well, or -- with always -- whenever a table for it
// exists, even where another language would sound closer.
func (c *Converter) Prefer(lang string, always bool) {
	c.prefer = lang
	c.preferAlways = always
}

// New creates a Converter for the given source language.
// Supported: "estonian", "karelian", "latvian", "german", "latin".
func New(lang string) *Converter {
	switch strings.ToLower(lang) {
	case "estonian":
		return newEstonian()
	case "karelian":
		return newKarelian()
	case "latvian":
		return newLatvian()
	case "german":
		return newGerman()
	case "latin":
		return newLatin()
	default:
		return nil
	}
}

// SetWord adds a whole-word override that bypasses table conversion.
func (c *Converter) SetWord(word string, r Result) {
	if c.words == nil {
		c.words = map[string]Result{}
	}
	c.words[strings.ToLower(word)] = r
}

// Convert converts a lyrics word to phonemes.
// Returns empty Result for special markers ("-", "+", "sil", "br", etc.).
func (c *Converter) Convert(word string) Result {
	lower := strings.ToLower(word)

	switch lower {
	case "-", "+", "sil", "br", "sp", "ap", "":
		return Result{}
	}

	if r, ok := c.words[lower]; ok {
		return r
	}

	lang := c.pickLang(lower)
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
		Phonemes: table.convert(w, c.skip, c.vowels),
	}
}

// pickLang chooses the target language for a word: the voice's own language
// when it is among the candidates (or always, if asked), otherwise the best
// candidate for the word.
func (c *Converter) pickLang(word string) string {
	candidates := c.selectLangs(word)
	if c.prefer != "" && c.tables[c.prefer] != nil {
		if c.preferAlways || slices.Contains(candidates, c.prefer) {
			return c.prefer
		}
	}
	return candidates[0]
}

type phoneTable struct {
	// Digraphs are checked first (two-char sequences like "sh", "ng", "ts").
	// Also used for special geminates that don't simply double (Korean "ss"→"s_t").
	digraphs map[string][]string
	// Expanded maps individual runes to multiple phonemes (checked before singles).
	// Unlike singles, expanded entries skip geminate handling.
	expanded map[rune][]string
	// Singles maps individual runes to their phoneme.
	// Geminates (doubled chars) automatically double the single phoneme.
	singles map[rune]string
}

func (t *phoneTable) convert(word string, skip func(rune) bool, vowels string) string {
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

		// Check expanded (single chars → multiple phonemes, no geminate handling).
		if t.expanded != nil {
			if ph, ok := t.expanded[ch]; ok {
				phonemes = append(phonemes, ph...)
				continue
			}
		}

		ph, ok := t.singles[ch]
		if !ok {
			continue
		}

		// Check geminate (doubled character).
		if i+1 < len(runes) && runes[i+1] == ch {
			// Doubled vowels emit a single phoneme to avoid rearticulation.
			if strings.ContainsRune(vowels, ch) {
				phonemes = append(phonemes, ph)
			} else {
				phonemes = append(phonemes, ph, ph)
			}
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
