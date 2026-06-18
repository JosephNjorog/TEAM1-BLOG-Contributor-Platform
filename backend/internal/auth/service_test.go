package auth

import "testing"

func TestIsValidAvalancheAddress(t *testing.T) {
	cases := []struct {
		addr string
		want bool
	}{
		{"0x1234567890123456789012345678901234567890", true},
		{"0xABCDEFabcdef00000000000000000000000000AB", true},
		{"1234567890123456789012345678901234567890", false},     // missing 0x
		{"0x12345678901234567890123456789012345678", false},     // too short
		{"0x123456789012345678901234567890123456789012", false}, // too long
		{"0xZZZZ567890123456789012345678901234567890", false},   // non-hex chars
		{"", false},
	}
	for _, c := range cases {
		if got := IsValidAvalancheAddress(c.addr); got != c.want {
			t.Errorf("IsValidAvalancheAddress(%q) = %v, want %v", c.addr, got, c.want)
		}
	}
}
