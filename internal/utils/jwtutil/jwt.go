package jwtutil

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"github.com/microservices-simulator-api/internal/config"
	"github.com/microservices-simulator-api/internal/users"
)

var (
	signingMethod = jwt.SigningMethodRS512
)

type (
	Authentication struct {
		cfg        config.SecurityConfig
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
	}

	Claims = jwt.Claims

	Token = jwt.Token
)

func NewAuthentication(cfg config.SecurityConfig) (*Authentication, error) {
	publicKey, err := rsaPublicKey(cfg.PublicKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := rsaPrivateKey(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &Authentication{cfg, publicKey, privateKey}, nil
}

func (a *Authentication) Generate(user users.User) (string, error) {
	claims := NewUserClaims(user.Id, user.Name, a.cfg.TokenDomain, a.JIT(), a.cfg.TokenMaxAge)
	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(a.privateKey)
}

func (a *Authentication) Parse(accessToken string) (interface{}, error) {
	claims := &UserClaims{}

	_, err := jwt.ParseWithClaims(accessToken, claims, a.GetPublicKey)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (a *Authentication) GetPublicKey(token *jwt.Token) (interface{}, error) {
	return a.publicKey, nil
}

func (a *Authentication) JIT() string {
	return a.cfg.Id
}

func (a *Authentication) NewClaims() jwt.Claims {
	return new(UserClaims)
}

func rsaPrivateKey(key string) (*rsa.PrivateKey, error) {
	pvKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(pvKey)
}

func rsaPublicKey(key string) (*rsa.PublicKey, error) {
	pbKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(pbKey)
}
