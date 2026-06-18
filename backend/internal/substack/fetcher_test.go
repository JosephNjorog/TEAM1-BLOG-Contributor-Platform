package substack

import "testing"

func TestParsePubDate(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr bool // true means we expect the zero value back
	}{
		{"named zone (what Substack's feed actually emits)", "Mon, 08 Jun 2026 13:30:58 GMT", false},
		{"numeric offset", "Mon, 08 Jun 2026 13:30:58 +0000", false},
		{"garbage", "not a date", true},
		{"empty", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := parsePubDate(c.raw)
			if c.wantErr {
				if !got.IsZero() {
					t.Errorf("parsePubDate(%q) = %v, want the zero value", c.raw, got)
				}
				return
			}
			if got.IsZero() {
				t.Errorf("parsePubDate(%q) returned the zero value, want a parsed date", c.raw)
			}
			if got.Year() != 2026 {
				t.Errorf("parsePubDate(%q).Year() = %d, want 2026", c.raw, got.Year())
			}
		})
	}
}
