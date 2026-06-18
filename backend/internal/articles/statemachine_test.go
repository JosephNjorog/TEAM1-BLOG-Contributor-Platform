package articles

import "testing"

func TestCanTransition_AllowedMoves(t *testing.T) {
	allowed := []struct {
		from, to Status
	}{
		{StatusDraft, StatusSubmitted},
		{StatusSubmitted, StatusChangesRequested},
		{StatusResubmitted, StatusChangesRequested},
		{StatusChangesRequested, StatusResubmitted},
		{StatusSubmitted, StatusEditorialApproved},
		{StatusResubmitted, StatusEditorialApproved},
		{StatusEditorialApproved, StatusBannerUploaded},
		{StatusBannerUploaded, StatusPublished},
		{StatusPublished, StatusPaymentInitiated},
		{StatusPaymentInitiated, StatusPaymentConfirmed},
	}
	for _, c := range allowed {
		if !CanTransition(c.from, c.to) {
			t.Errorf("CanTransition(%s, %s) = false, want true", c.from, c.to)
		}
	}
}

func TestCanTransition_DisallowedMoves(t *testing.T) {
	disallowed := []struct {
		from, to Status
	}{
		{StatusDraft, StatusPublished},                    // skipping the entire pipeline
		{StatusDraft, StatusEditorialApproved},            // can't approve a draft directly
		{StatusPublished, StatusDraft},                    // no going backwards
		{StatusEditorialApproved, StatusPublished},        // banner step can't be skipped
		{StatusChangesRequested, StatusEditorialApproved}, // must resubmit first
		{StatusPaymentConfirmed, StatusPaymentInitiated},  // no going backwards
		{StatusDraft, StatusDraft},                        // not a real transition
	}
	for _, c := range disallowed {
		if CanTransition(c.from, c.to) {
			t.Errorf("CanTransition(%s, %s) = true, want false", c.from, c.to)
		}
	}
}

func TestValidateTransition(t *testing.T) {
	if err := ValidateTransition(StatusDraft, StatusSubmitted); err != nil {
		t.Errorf("ValidateTransition(draft, submitted) returned error: %v", err)
	}
	if err := ValidateTransition(StatusDraft, StatusPublished); err == nil {
		t.Error("ValidateTransition(draft, published) returned nil, want an error")
	}
}

func TestStatus_Editable(t *testing.T) {
	editable := []Status{StatusDraft, StatusChangesRequested}
	for _, s := range editable {
		if !s.Editable() {
			t.Errorf("%s.Editable() = false, want true", s)
		}
	}

	notEditable := []Status{
		StatusSubmitted, StatusResubmitted, StatusEditorialApproved,
		StatusBannerUploaded, StatusPublished, StatusPaymentInitiated, StatusPaymentConfirmed,
	}
	for _, s := range notEditable {
		if s.Editable() {
			t.Errorf("%s.Editable() = true, want false", s)
		}
	}
}
