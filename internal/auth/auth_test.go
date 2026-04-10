package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"password123", true},
		{"wrongpassword", false},
	}

	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	for _, c := range cases {
		match, err := CheckPasswordHash(c.input, hashedPassword)
		if err != nil {
			t.Errorf("Error checking password hash: %v", err)
			continue
		}
		if match != c.expected {
			t.Errorf("For input '%s', expected %v but got %v", c.input, c.expected, match)
		}
	}
}

func TestJWT(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	tokenSecret := "my_secret_key"
	expiresIn := 1 * time.Hour
	token, err := MakeJWT(uuid.MustParse(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Error creating JWT: %v", err)
	}

	returnedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}

	if returnedUserID.String() != userID {
		t.Errorf("Expected user ID '%s' but got '%s'", userID, returnedUserID.String())
	}

	tokenSecretWrong := "wrong_secret_key"
	_, err = ValidateJWT(token, tokenSecretWrong)
	if err == nil {
		t.Errorf("Expected error when validating JWT with wrong secret, but got none")
	}

	expiredToken, err := MakeJWT(uuid.MustParse(userID), tokenSecret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Error creating expired JWT: %v", err)
	}

	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Errorf("Expected error when validating expired JWT, but got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer my_token")

	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("Error getting bearer token: %v", err)
	}

	expectedToken := "my_token"
	if token != expectedToken {
		t.Errorf("Expected token '%s' but got '%s'", expectedToken, token)
	}

	headersMissing := http.Header{}
	_, err = GetBearerToken(headersMissing)
	if err == nil {
		t.Errorf("Expected error when authorization header is missing, but got none")
	}

	headersInvalid := http.Header{}
	headersInvalid.Set("Authorization", "InvalidHeader")
	_, err = GetBearerToken(headersInvalid)
	if err == nil {
		t.Errorf("Expected error when bearer token is missing, but got none")
	}
}
