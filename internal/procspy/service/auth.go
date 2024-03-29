package service

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

type Auth struct {
	key    []byte
	pubKey []byte
	ttl    time.Duration
}

func NewAuth() *Auth {
	ret := &Auth{
		ttl: time.Hour * 24,
	}

	err := ret.GenerateKeys()
	if err != nil {
		log.Printf("[Auth] Error generating keys: %s", err)
		return nil
	}

	return ret
}

func (a *Auth) GenerateKeys() error {
	log.Printf("[Auth] Creating new key")

	bitSize := 4096

	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Printf("[Auth] Error generating key: %s", err)
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

	log.Printf("[Auth] Keys created - pub: %s", string(a.pubKey))
	return nil
}

func (a *Auth) GetPubKey() (string, error) {
	return string(a.pubKey), nil
}

func (a *Auth) CreateToken(user string, content map[string]string) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(a.key)
	if err != nil {
		return "", fmt.Errorf("[Auth] Create: parse key: %w", err)
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
		return "", fmt.Errorf("[Auth] Create: sign token: %w", err)
	}

	return token, nil
}

func (a *Auth) Validate(token string) (map[string]string, bool, error) {
	expired := false
	key, err := jwt.ParseRSAPublicKeyFromPEM(a.pubKey)
	if err != nil {
		return nil, expired, fmt.Errorf("[Auth] validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("[Auth] unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return nil, expired, fmt.Errorf("[Auth] validate error: %w", err)
	}

	if !tok.Valid {
		return nil, expired, fmt.Errorf("[Auth] validate: invalid token format")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, expired, fmt.Errorf("[Auth] validate: invalid claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		expired = true
		return nil, expired, fmt.Errorf("[Auth] validate: expired token")
	}

	if !claims.VerifyIssuedAt(time.Now().Unix(), true) {
		return nil, expired, fmt.Errorf("[Auth] validate: invalid iat")
	}

	if !claims.VerifyNotBefore(time.Now().Unix(), true) {
		expired = true
		return nil, expired, fmt.Errorf("[Auth] validate: invalid nbf")
	}

	content, ok := claims["dat"].(map[string]string)

	if !ok {

		return nil, expired, fmt.Errorf("[Auth] validate: invalid sub")
	}

	if sub, ok := claims["sub"].(string); !ok || sub != "procspy" {
		return nil, expired, fmt.Errorf("[Auth] validate: invalid sub")
	}

	if user, ok := content["user"]; !ok || user == "" {
		return nil, expired, fmt.Errorf("[Auth] validate: invalid user")
	} else {
		fmt.Printf("[Auth] validate: %s\n", user)
	}

	return content, expired, nil
}

func (a *Auth) Cypher(data string) (string, error) {
	return Cypher(data, a.pubKey)
}

func (a *Auth) Decypher(data string) (string, error) {
	return Decypher(data, a.key)
}

func Cypher(data string, key []byte) (string, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return "", fmt.Errorf("[Auth] cypher: parse key: %w", err)
	}

	enc, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))
	if err != nil {
		return "", fmt.Errorf("[Auth] cypher: encrypt: %w", err)
	}

	return string(enc), nil
}

func Decypher(data string, key []byte) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return "", fmt.Errorf("[Auth] decypher: parse key: %w", err)
	}

	dec, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(data))
	if err != nil {
		return "", fmt.Errorf("[Auth] decypher: decrypt: %w", err)
	}

	return string(dec), nil
}
