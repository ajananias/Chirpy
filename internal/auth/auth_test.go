package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashing(t *testing.T) {
	password := "argon2test"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Failed to hash password: %v", err)
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
		t.Errorf("Verification process crashed: %v", err)
	}
	if !match {
		t.Error("Verification failed: valid password was rejected")
	}

	wrongMatch, err := CheckPasswordHash("WrongPassword", hash)
	if err != nil {
		t.Errorf("Verification process crashed: %v", err)
	}
	if wrongMatch {
		t.Error("Verification error: Invalid password was accepted")
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "testjwt"
	expiresIn := time.Hour

	// create a JWT
	jwtToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Error creating the JWT: %v", err)
	}
	// validate
	validUserID, err := ValidateJWT(jwtToken, tokenSecret)
	if err != nil {
		t.Errorf("Couldn't validate user id with generated token and correct secret: %v", err)
	}
	if validUserID != userID {
		t.Errorf("Incorrect user validation: expected %s & obtained %s", userID, validUserID)
	}

	// Expired token
	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("Token creation failed: %v", err)
	}
	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Errorf("Test failed: Unable to catch Expired Token")
	}

	// Wrong token secret
	wrongSecret := "wrongsecret"
	_, err = ValidateJWT(jwtToken, wrongSecret)
	if err == nil {
		t.Error("Test failed: Validated the Wrong Secret")
	}

	// Check for authorization header
	var emptyHeaders http.Header
	bearerToken, err := GetBearerToken(emptyHeaders)
	if err == nil {
		t.Error("Obtained token from empty header")
	}

	header := http.Header{
		"Authorization": []string{"Bearer ${jwtTokenTest}"},
	}
	bearerToken, err = GetBearerToken(header)
	if err != nil {
		t.Errorf("Couldn't get token: %v", err)
	}
	if strings.HasPrefix(bearerToken, "Bearer ") {
		t.Errorf("Unable to strip prefix: Obtained %s", bearerToken)
	}

}
