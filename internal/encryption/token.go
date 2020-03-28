package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
	"time"
)

func GenerateVerificationToken(email string, expiry time.Duration, c *config.Configuration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      email,
		"expiryDate": time.Now().Add(expiry).Unix(),
		"timestamp":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.Keys.PasswordHashingKey))
	if err != nil {
		// TODO: log this
	}
	return tokenString
}

func GenerateSecurityToken(password string, c *config.Configuration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password":  password,
		"timestamp": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.Keys.PasswordHashingKey))
	if err != nil {
		// TODO: log this
	}
	return tokenString
}

func GenerateFromPassword(password string, cfg *config.Configuration) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := GenerateRandomBytes(cfg.ArgonParams.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, cfg.ArgonParams.Iterations, cfg.ArgonParams.Memory, cfg.ArgonParams.Parallelism, cfg.ArgonParams.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, cfg.ArgonParams.Memory, cfg.ArgonParams.Iterations, cfg.ArgonParams.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// For guidance and an outline process for choosing appropriate parameters see https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04#section-4.
func GenerateEncodedHash(pw string, cfg *config.Configuration) (string, error) {
	encodedHash, err := GenerateFromPassword(pw, cfg)
	if err != nil {
		return "", err
	}

	return encodedHash, nil
}
