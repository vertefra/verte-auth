package utils

import (
	"github/vertefra/verte_auth_server/config"
	"math/rand"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Encrypt will hash a string and return a encrypted string value taking care of all
// the conversion between strings and byte slices
func Encrypt(p string) (string, error) {
	ph := []byte(p)
	hashed, err := bcrypt.GenerateFromPassword(ph, 10)
	if err != nil {
		config.Err("Error Ashing password")
		config.Err("Error: ", err)
		return "", err
	}

	return string(hashed[:]), nil
}

// GenerateToken generates a jwt token used for authentication
// User Id is used to sign the token, d represent expiration time
// in hours
func GenerateToken(ID string, d int, key string) (string, error) {

	claims := jwt.MapClaims{}

	claims["ID"] = ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(d)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(key))
	if err != nil {
		config.Err("\nGenerateToken: Error generating the token")
		config.Err("Error: ", err)
		return "", err
	}

	return t, nil

}

// GenerateRanKey a random key 8 char long
func GenerateRanKey(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(){}:")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// IsEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
