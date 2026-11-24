package auth

import (
	"testing"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

func TestJwtFunction(t *testing.T) {
	id := uuid.New()
	secretKey := "secret-key"
	expiresAt := time.Minute * 2

	token, err := MakeJwt(id, secretKey, expiresAt)
	if err != nil {
		t.Fatal("function failed to process the input")
	}

	claims := &jwt.RegisteredClaims{}
	NewToken, err := jwt.ParseWithClaims(
		token,
		claims,
		func(*jwt.Token) (any, error) { return []byte(secretKey), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil || !NewToken.Valid {
		t.Fatal("invalid token")
	}
	if claims.Issuer != "todo-access" || claims.Subject != id.String() {
		t.Fatalf("unexpected claims: issuer=%q subject=%q", claims.Issuer, claims.Subject)
	}

	if claims.IssuedAt == nil ||
		claims.ExpiresAt == nil ||
		claims.ExpiresAt.Sub(claims.IssuedAt.Time) != 2*time.Minute {
		t.Fatalf("unexpected token lifetime: issued=%v expires=%v", claims.IssuedAt, claims.ExpiresAt)
	}
}

func TestJwtValidation(t *testing.T) {
	id := uuid.New()
	secret := "test-secret"
	valid, err := MakeJwt(id, secret, time.Hour)
	if err != nil {
		t.Fatal("Couldn't create a valid jwt")
	}

	result, err := ValidateJwt(valid, secret)
	if err != nil || result != id {
		t.Fatalf("The error: %v and result: %s", err, result)
	}

	expired, err := MakeJwt(id, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Error: %v", expired)
	}

	badRequest, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Subject:   "Not uuid",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString([]byte(secret))

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	for _, tc := range []struct {
		name   string
		token  string
		secret string
	}{
		{name: "Wrong key", token: valid, secret: "test-secret"},
		{name: "Expired", token: expired, secret: secret},
		{name: "Bad request", token: badRequest, secret: secret},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ValidateJwt(tc.token, tc.secret); err == nil {
				t.Fatal("expected token validation error")
			}
		})
	}

}
