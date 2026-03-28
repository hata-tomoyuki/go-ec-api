package auth

import "golang.org/x/crypto/bcrypt"

func hashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPasswordHash(plain, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
