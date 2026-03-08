package phonemes

import "testing"

func TestEstonian_BasicConversion(t *testing.T) {
	c := New("estonian")
	tests := []struct {
		word     string
		phonemes string
		lang     string
	}{
		{"la", "l a", "spanish"},
		{"tere", "t e r e", "spanish"},
		{"küla", "k y l a", "mandarin"},    // ü → mandarin
		{"öö", "9", "cantonese"},            // ö → cantonese, doubled vowel = single phoneme
		{"õhtu", "7 x t u", "mandarin"},     // õ → mandarin
		{"äike", "E i k e", "cantonese"},    // ä → cantonese
		{"hmm", "h m=", "cantonese"},        // word override
		{"kass", "k a s s", "spanish"},      // geminate consonant
		{"saal", "s a l", "spanish"},        // doubled vowel = single phoneme
		{"tsikk", "t s i k k", "spanish"},   // ts digraph + geminate
	}
	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			r := c.Convert(tt.word)
			if r.Phonemes != tt.phonemes {
				t.Errorf("Convert(%q).Phonemes = %q, want %q", tt.word, r.Phonemes, tt.phonemes)
			}
			if r.Language != tt.lang {
				t.Errorf("Convert(%q).Language = %q, want %q", tt.word, r.Language, tt.lang)
			}
		})
	}
}

func TestEstonian_SpecialMarkers(t *testing.T) {
	c := New("estonian")
	for _, marker := range []string{"-", "+", "sil", "br", "sp", "ap", ""} {
		r := c.Convert(marker)
		if r.Phonemes != "" {
			t.Errorf("Convert(%q) should return empty, got %q", marker, r.Phonemes)
		}
	}
}

func TestKarelian_BasicConversion(t *testing.T) {
	c := New("karelian")
	tests := []struct {
		word     string
		phonemes string
		lang     string
	}{
		{"kala", "k a l a", "spanish"},
		{"kyly", "k y l y", "mandarin"},     // y → mandarin
		{"öä", "9 E", "cantonese"},          // ö → cantonese
		{"hmm", "h m=", "cantonese"},        // word override
		{"čakku", "ch a k k u", "spanish"},  // č mapped + geminate
		{"šakki", "sh a k k i", "spanish"},  // š mapped + geminate
	}
	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			r := c.Convert(tt.word)
			if r.Phonemes != tt.phonemes {
				t.Errorf("Convert(%q).Phonemes = %q, want %q", tt.word, r.Phonemes, tt.phonemes)
			}
			if r.Language != tt.lang {
				t.Errorf("Convert(%q).Language = %q, want %q", tt.word, r.Language, tt.lang)
			}
		})
	}
}

func TestKarelian_PalatalizationSkip(t *testing.T) {
	c := New("karelian")
	// Palatalization marks should be skipped.
	r := c.Convert("d'a")
	if r.Phonemes != "d a" {
		t.Errorf("Convert(\"d'a\").Phonemes = %q, want %q", r.Phonemes, "d a")
	}
	// All palatalization variants.
	for _, mark := range []string{"'", "\u2018", "\u2019", "\u02BC"} {
		word := "d" + mark + "a"
		r := c.Convert(word)
		if r.Phonemes != "d a" {
			t.Errorf("Convert(%q).Phonemes = %q, want %q", word, r.Phonemes, "d a")
		}
	}
}

func TestGerman_LanguageSelection(t *testing.T) {
	c := New("german")
	tests := []struct {
		word string
		lang string
	}{
		{"über", "mandarin"},   // ü → mandarin
		{"schön", "cantonese"}, // ö → cantonese
		{"ach", "spanish"},     // ch after a → spanish (ach-Laut)
		{"ich", "japanese"},    // ch after i → japanese (ich-Laut)
		{"mächtig", "japanese"}, // ch after non-back-vowel → japanese (priority 3 > ä priority 4)
		{"rot", "korean"},     // r → korean
		{"groß", "korean"},    // r → korean (priority 5 > ß priority 6)
		{"tag", "japanese"},    // simple → japanese
	}
	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			r := c.Convert(tt.word)
			if r.Language != tt.lang {
				t.Errorf("Convert(%q).Language = %q, want %q", tt.word, r.Language, tt.lang)
			}
		})
	}
}

func TestGerman_BasicConversion(t *testing.T) {
	c := New("german")
	tests := []struct {
		word     string
		phonemes string
		lang     string
	}{
		{"tag", "t a g", "japanese"},
		{"rot", "4 o t", "korean"},          // r → korean '4'
		{"über", "y p e r\\`", "mandarin"},  // ü → mandarin
	}
	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			r := c.Convert(tt.word)
			if r.Phonemes != tt.phonemes {
				t.Errorf("Convert(%q).Phonemes = %q, want %q", tt.word, r.Phonemes, tt.phonemes)
			}
			if r.Language != tt.lang {
				t.Errorf("Convert(%q).Language = %q, want %q", tt.word, r.Language, tt.lang)
			}
		})
	}
}

func TestGerman_SchNormalization(t *testing.T) {
	c := New("german")
	// "sch" should be normalized to ʃ and mapped appropriately.
	r := c.Convert("schlag")
	// "schlag" → Japanese (simple, no special chars after normalization)
	if r.Language != "japanese" {
		t.Errorf("Convert(\"schlag\").Language = %q, want %q", r.Language, "japanese")
	}
	// sch→sh, l→r, a→a, g→g in Japanese
	if r.Phonemes != "sh r a g" {
		t.Errorf("Convert(\"schlag\").Phonemes = %q, want %q", r.Phonemes, "sh r a g")
	}
}

func TestGerman_WordInitialStSp(t *testing.T) {
	c := New("german")
	// Word-initial "st" should be palatalized in Mandarin.
	// "stein" → mandarin (has no special chars, but let's check via forced word)
	// Actually "stein" has "ei" digraph and is simple → japanese.
	// Let's use a word that forces Mandarin: "stück" (has ü).
	r := c.Convert("stück")
	if r.Language != "mandarin" {
		t.Errorf("Convert(\"stück\").Language = %q, want %q", r.Language, "mandarin")
	}
	// After normalization: "ʃtück" → "ʃtyck" (ü→y, ß already handled)
	// ʃt digraph → {"s`", "t"}, y → "y", k → "k", k → "k"
	// Wait: "stück" → normalize: "sch" not present, starts with "st" → "ʃtück"
	// Then ß→s (none), x→ks (none). Result: "ʃtück"
	// Mandarin table: ʃt→{"s`","t"}, ü→"y", ck→{"k"}
	if r.Phonemes != "s` t y k" {
		t.Errorf("Convert(\"stück\").Phonemes = %q, want %q", r.Phonemes, "s` t y k")
	}

	// Non-initial "st" should NOT be palatalized.
	// "fast" → japanese (simple). In japanese table, no "st" digraph.
	r = c.Convert("fast")
	if r.Phonemes != "f a s t" {
		t.Errorf("Convert(\"fast\").Phonemes = %q, want %q", r.Phonemes, "f a s t")
	}
}

func TestGerman_SetWord(t *testing.T) {
	c := New("german")
	c.SetWord("test", Result{
		Language: "english",
		Phoneset: "arpabet",
		Phonemes: "t eh s t",
	})
	r := c.Convert("test")
	if r.Phonemes != "t eh s t" {
		t.Errorf("Convert(\"test\").Phonemes = %q, want %q", r.Phonemes, "t eh s t")
	}
	if r.Language != "english" {
		t.Errorf("Convert(\"test\").Language = %q, want %q", r.Language, "english")
	}
}

func TestNew_UnknownLanguage(t *testing.T) {
	c := New("french")
	if c != nil {
		t.Error("New(\"french\") should return nil for unsupported language")
	}
}

func TestGerman_GeminateHandling(t *testing.T) {
	c := New("german")
	// "bett" → japanese (simple). "tt" = geminate consonant → "t t"
	r := c.Convert("bett")
	if r.Phonemes != "b e t t" {
		t.Errorf("Convert(\"bett\").Phonemes = %q, want %q", r.Phonemes, "b e t t")
	}
	// "see" → japanese (simple). "ee" = doubled vowel → single "e"
	r = c.Convert("see")
	if r.Phonemes != "s e" {
		t.Errorf("Convert(\"see\").Phonemes = %q, want %q", r.Phonemes, "s e")
	}
}

func TestGerman_Digraphs(t *testing.T) {
	c := New("german")
	// "ng" digraph in Japanese
	r := c.Convert("gang")
	// "gang" → japanese: g→g, a→a, ng→{"N","g"}
	if r.Phonemes != "g a N g" {
		t.Errorf("Convert(\"gang\").Phonemes = %q, want %q", r.Phonemes, "g a N g")
	}
	// "ei" digraph in Japanese
	r = c.Convert("ein")
	// "ein" → japanese: ei→{"a","i"}, n→n
	if r.Phonemes != "a i n" {
		t.Errorf("Convert(\"ein\").Phonemes = %q, want %q", r.Phonemes, "a i n")
	}
}

func TestPhoneset(t *testing.T) {
	tests := []struct {
		lang     string
		phoneset string
	}{
		{"japanese", "romaji"},
		{"english", "arpabet"},
		{"mandarin", "xsampa"},
		{"cantonese", "xsampa"},
		{"korean", "xsampa"},
		{"spanish", "xsampa"},
	}
	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := phonesetForLanguage(tt.lang)
			if got != tt.phoneset {
				t.Errorf("phonesetForLanguage(%q) = %q, want %q", tt.lang, got, tt.phoneset)
			}
		})
	}
}
