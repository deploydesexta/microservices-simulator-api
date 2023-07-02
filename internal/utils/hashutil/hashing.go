package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

type (
	PasswordHashing interface {
		Hash(password string) (string, error)
		Compare(password string, hash string) error
	}

	BcryptPasswordHashing struct {
		salts int
	}
)

func Sha256(data string) string {
	hashing := sha256.New()
	hashing.Write([]byte(data))
	return hex.EncodeToString(hashing.Sum(nil))
}

func NewPasswordHashing(salts int) *BcryptPasswordHashing {
	return &BcryptPasswordHashing{salts}
}

func (b *BcryptPasswordHashing) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.salts)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b *BcryptPasswordHashing) Compare(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
