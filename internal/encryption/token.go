package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/clog"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
)

const (
	errIncompatibleVersion = "version de argon2 incompatible"
)

func GenerateVerificationToken(email string, expiry time.Duration, c *config.EnigmaConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      email,
		"expiryDate": time.Now().Add(expiry).Unix(),
		"timestamp":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.Keys.PasswordHashingKey))
	if err != nil {
		clog.Error("Error generating verification token", "generate-verification-token", err, nil)
		return "", err
	}
	return tokenString, nil
}

func GenerateSecurityToken(password string, c *config.EnigmaConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password":  password,
		"timestamp": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.Keys.PasswordHashingKey))
	if err != nil {
		clog.Error("Error generating security token", "generate-security-token", err, nil)
		return "", err
	}

	return tokenString, nil
}

func generateFromPassword(password string, cfg *config.EnigmaConfig) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(cfg.ArgonParams.SaltLength)
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

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		clog.Error("Error generating random bytes", "generate-random-bytes", err, nil)
		return nil, err
	}

	return b, nil
}

// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// For guidance and an outline process for choosing appropriate parameters see https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04#section-4.
func GenerateEncodedHash(pw string, cfg *config.EnigmaConfig) (string, error) {
	encodedHash, err := generateFromPassword(pw, cfg)
	if err != nil {
		clog.Error("Error generating encoded hash", "generate-encoded-hash", err, nil)
		return "", err
	}

	return encodedHash, nil
}

func DecodeHash(encodedHash string) (p *config.ArgonParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		msg := "encoded hash string is not 6"
		clog.Error(msg, "decode-hash", errors.New(msg), map[string]string{"len": strconv.Itoa(len(vals))})
		return nil, nil, nil, errors.New(msg)
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		clog.Error("Error finding arguments in string", "decode-hash", err, nil)
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		err := errors.New(errIncompatibleVersion)
		clog.Error(errIncompatibleVersion, "decode-hash", err, map[string]string{"version": strconv.Itoa(version)})
		return nil, nil, nil, err
	}

	p = &config.ArgonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		clog.Error("Error scanning params", "decode-hash", err, nil)
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		clog.Error("Error decoding salt base64 string", "decode-hash", err, nil)
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		clog.Error("Error decoding hash base64 string", "decode-hash", err, nil)
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
