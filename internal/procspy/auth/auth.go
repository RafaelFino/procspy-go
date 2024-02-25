package procspy_auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Authorization struct {
	key    []byte
	pubKey []byte
	ttl    time.Duration
}

func NewAuthorization() *Authorization {
	ret := &Authorization{}

	err := ret.GenerateKeys()
	if err != nil {
		log.Printf("[Authorization] Error generating keys: %s", err)
		return nil
	}

	return ret
}

func (a *Authorization) GenerateKeys() error {
	log.Printf("[Authorization] Creating new key")

	bitSize := 4096

	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Printf("[Authorization] Error generating key: %s", err)
		return err
	}

	a.key = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	pub := key.Public()

	a.pubKey = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	log.Printf("[Authorization] Keys created - pub: %s", string(a.pubKey))
	return nil
}

func (a *Authorization) GetPubKey() (string, error) {
	return string(a.pubKey), nil
}

func (a *Authorization) CreateToken(user string, content interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(a.key)
	if err != nil {
		return "", fmt.Errorf("[Authorization] Create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content               // Our custom data.
	claims["exp"] = now.Add(a.ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()            // The time at which the token was issued.
	claims["nbf"] = now.Unix()            // The time before which the token must be disregarded.
	claims["sub"] = user                  // The subject of the token.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("[Authorization] Create: sign token: %w", err)
	}

	return token, nil
}

func (a *Authorization) Validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(a.pubKey)
	if err != nil {
		return nil, fmt.Errorf("[Authorization] validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("[Authorization] unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("[Authorization] validate error: %w", err)
	}

	if !tok.Valid {
		return nil, fmt.Errorf("[Authorization] validate: invalid token format")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("[Authorization] validate: invalid claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("[Authorization] validate: expired token")
	}

	if !claims.VerifyIssuedAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("[Authorization] validate: invalid iat")
	}

	if !claims.VerifyNotBefore(time.Now().Unix(), true) {
		return nil, fmt.Errorf("[Authorization] validate: invalid nbf")
	}

	user, ok := claims["sub"]

	if !ok {
		return nil, fmt.Errorf("[Authorization] validate: invalid sub")
	}

	fmt.Printf("[Authorization] validate: %s\n", user)

	return claims["dat"], nil
}
