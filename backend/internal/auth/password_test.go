package auth

import "testing"

func TestHashAndCheckPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "correct-password" {
		t.Fatal("HashPassword returned the plaintext unchanged")
	}
	if !CheckPassword(hash, "correct-password") {
		t.Error("CheckPassword rejected the correct password")
	}
	if CheckPassword(hash, "wrong-password") {
		t.Error("CheckPassword accepted an incorrect password")
	}
}

func TestHashPasswordIsSalted(t *testing.T) {
	hashA, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	hashB, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hashA == hashB {
		t.Error("two hashes of the same password were identical - bcrypt should salt each call")
	}
}
