package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJwt(id uuid.UUID, secretKey string, expiresAt time.Duration) (string, error) {
	createdAt := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.RegisteredClaims{
			Issuer:    "todo-access",
			IssuedAt:  jwt.NewNumericDate(createdAt),
			ExpiresAt: jwt.NewNumericDate(createdAt.Add(expiresAt)),
			Subject:   id.String(),
		})
	return token.SignedString([]byte(secretKey))
}

func ValidateJwt(tokenInput, secretToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenInput,
		jwt.RegisteredClaims{},
		func(*jwt.Token) (any, error) { return []byte(secretToken), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return uuid.Nil(), fmt.Errorf("Couldn't validate the token")
	}

	if !token.Valid {
		return uuid.Nil(), fmt.Errorf("Invalid token")
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil(), fmt.Errorf("Couldn't get the subject")
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil(), fmt.Errorf("Can't parse the uuid")
	}
	return id, nil
}

func GetBearer(r http.Request) (string, error) {
	part := strings.Fields(r.Header.Get("Authorization"))
	if len(part) != 2 || strings.EqualFold(part[0], "Bearer") {
		return "", fmt.Errorf("The Header is invalid")
	}
	return part[1], nil

}
