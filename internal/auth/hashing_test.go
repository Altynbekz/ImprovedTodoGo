package auth

import (
	"testing"
)

func TestHashing(t *testing.T) {
	tests := []struct {
		name          string
		plainPassword string
	}{
		{name: "Normal password", plainPassword: "Altynbek"},
		{name: "Different capital letter", plainPassword: "altynbek"},
		{name: "Empty", plainPassword: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := HashPassword(tc.plainPassword)
			if err != nil {
				t.Fatalf("The hash did not hash %v", err)
			}
			if hash == "" {
				t.Fatalf("Hash returned an empty string")
			}
			match, err := CheckPassword(tc.plainPassword, hash)
			if err != nil {
				t.Fatal("The check password failed to work")
			}
			if !match {
				t.Fatal("The hash does not match when given the same password")
			}
		})
	}
}

func TestWrongHashFunction(t *testing.T) {
	hash, err := HashPassword("Altynbek")
	if err != nil {
		t.Fatalf("Couldnt hash the password u gave me")
	}

	match, err := CheckPassword("Wrong", hash)
	if err != nil {
		t.Fatal("Couldnt use the check function to verify ur password")
	}
	if match {
		t.Fatal("The function did not work properly it gave true when it should have been false ")
	}
}
