package articles

import "fmt"

// transitions whitelists every allowed Status -> Status move in the
// platform's article lifecycle (PRD section 6). Each entry maps the
// resulting status to the set of statuses it may be entered from.
var transitions = map[Status][]Status{
	StatusSubmitted:         {StatusDraft},
	StatusChangesRequested:  {StatusSubmitted, StatusResubmitted},
	StatusResubmitted:       {StatusChangesRequested},
	StatusEditorialApproved: {StatusSubmitted, StatusResubmitted},
	StatusBannerUploaded:    {StatusEditorialApproved},
	StatusPublished:         {StatusBannerUploaded},
	StatusPaymentInitiated:  {StatusPublished},
	StatusPaymentConfirmed:  {StatusPaymentInitiated},
}

// CanTransition reports whether moving from `from` to `to` is a legal step
// in the article state machine.
func CanTransition(from, to Status) bool {
	for _, allowed := range transitions[to] {
		if allowed == from {
			return true
		}
	}
	return false
}

func ValidateTransition(from, to Status) error {
	if !CanTransition(from, to) {
		return fmt.Errorf("cannot transition article from %q to %q", from, to)
	}
	return nil
}
