package utility

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

var saltLen = 8

func init() {
	rand.Seed(time.Now().UnixNano())
}

func HashPassword(password string) (hashed string, salt string, err error) {
	salt = randomSalt()
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return string(hash), salt, nil
}

func PasswordIsValid(password, salt, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	return err == nil
}

func randomSalt() string {
	var salt string
	for i := 0; i < saltLen; i++ {
		char := rand.Int31n(122) + 41
		salt += string(char)
	}

	return salt
}
