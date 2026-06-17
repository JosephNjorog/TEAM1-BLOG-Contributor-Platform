package substack

import "time"

// Post is a single entry fetched from the publication's feed, before it's
// matched to a contributor and persisted.
type Post struct {
	SubstackPostID string
	Title          string
	URL            string
	Author         string
	PublishedAt    time.Time
}

type Article struct {
	ID             string
	ContributorID  *string
	SubstackPostID string
	Title          string
	URL            string
	PublishedAt    time.Time
	SyncedAt       time.Time
}
