package articles

import "testing"

func TestCountWords(t *testing.T) {
	cases := []struct {
		name string
		html string
		want int
	}{
		{"empty", "", 0},
		{"plain text", "hello world", 2},
		{"single paragraph", "<p>Hello world, this is a test.</p>", 6},
		{"multiple paragraphs", "<p>First paragraph here.</p><p>Second one too.</p>", 6},
		{"headings and links", "<h2>A Title</h2><p>Some <a href=\"https://x.com\">linked</a> text.</p>", 5},
		{"only markup, no text", "<p></p><br/><img src=\"x.jpg\"/>", 0},
		{"whitespace-only content", "<p>   </p>\n\t", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := CountWords(c.html); got != c.want {
				t.Errorf("CountWords(%q) = %d, want %d", c.html, got, c.want)
			}
		})
	}
}
