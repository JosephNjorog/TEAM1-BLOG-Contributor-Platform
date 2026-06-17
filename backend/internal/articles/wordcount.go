package articles

import (
	"regexp"
	"strings"
)

var tagRE = regexp.MustCompile(`<[^>]*>`)

// CountWords strips HTML markup produced by the rich text editor and counts
// whitespace-separated tokens, used to recompute word_count server-side on
// every save so the contributor's editor and the moderator's queue agree.
func CountWords(htmlContent string) int {
	text := tagRE.ReplaceAllString(htmlContent, " ")
	fields := strings.Fields(text)
	return len(fields)
}
