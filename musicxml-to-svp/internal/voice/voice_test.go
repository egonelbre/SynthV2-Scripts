package voice

import "testing"

func TestParseVoicePart(t *testing.T) {
	tests := []struct {
		name string
		want VoicePart
	}{
		{"Soprano", Soprano},
		{"Soprano 1", Soprano},
		{"Alto", Alto},
		{"Tenor", Tenor},
		{"Bass", Bass},
		{"Bass 2", Bass},
		{"Baritone", Baritone},
		{"Mezzo-Soprano", MezzoSoprano},
		{"Solo Soprano", Soprano},
		{"Solo Alto", Alto},
		{"Solo Alto 2", Alto},
		{"Solo Bass", Bass},
		{"Solo Tenor", Tenor},
		{"solo soprano", Soprano},
		{"SOLO BASS", Bass},
		{"Piano", Unknown},
		{"Solo", Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseVoicePart(tt.name)
			if got != tt.want {
				t.Errorf("ParseVoicePart(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestParsePartNum(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Bass", 1},
		{"Bass 2", 2},
		{"Soprano 1", 1},
		{"Solo Bass", 1},
		{"Solo Bass 2", 2},
		{"Solo Soprano 1", 1},
		{"Alto 3a", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePartNum(tt.name)
			if got != tt.want {
				t.Errorf("ParsePartNum(%q) = %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

func TestParseIsSolo(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"Solo Alto", true},
		{"Solo Soprano", true},
		{"solo bass 2", true},
		{"SOLO Tenor", true},
		{"Soprano 1", false},
		{"Bass", false},
		{"Alto", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseIsSolo(tt.name)
			if got != tt.want {
				t.Errorf("ParseIsSolo(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
