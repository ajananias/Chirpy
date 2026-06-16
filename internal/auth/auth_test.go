package auth

import (
	"strings"
	"testing"
)

func TestHashing(t *testing.T) {
	password := "argon2test"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Validate structural format properties
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("Expected hash to start with $argon2id$, got: %s", hash)
	}

	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		t.Errorf("Expected 5 hash components, gt %d", len(parts)-1)
	}

	// Test verification correctness
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Verification process crashed: %v", err)
	}
	if !match {
		t.Error("Verification failed: valid password was rejected")
	}

	wrongMatch, err := CheckPasswordHash("WrongPassword", hash)
	if err != nil {
		t.Fatalf("Verification process crashed: %v", err)
	}
	if wrongMatch {
		t.Error("Verification error: Invalid password was accepted")
	}
}
