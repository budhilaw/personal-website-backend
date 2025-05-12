package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 parameters
var (
	argon2Time    = uint32(1)
	argon2Memory  = uint32(64 * 1024) // 64MB
	argon2Threads = uint8(4)
	argon2KeyLen  = uint32(32)
	saltLen       = 16
)

// HashPassword hashes a password using Argon2id
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password with Argon2id
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Base64 encode the salt and hash for storage
	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	// Return the formatted hash string
	// Format: $argon2id$v=19$m=memory,t=time,p=threads$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2Memory, argon2Time, argon2Threads, saltBase64, hashBase64)

	return encodedHash, nil
}

// VerifyPassword compares a password with a hashed password
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the hash string
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the provided password with the same parameters and salt
	otherHash := argon2.IDKey([]byte(password), salt, params.time, params.memory, params.threads, params.keyLen)

	// Compare the computed hash with the stored hash
	// Use constant-time comparison to prevent timing attacks
	match := subtle.ConstantTimeCompare(hash, otherHash) == 1
	return match, nil
}

// argon2Params holds the parameters for Argon2
type argon2Params struct {
	memory  uint32
	time    uint32
	threads uint8
	keyLen  uint32
}

// decodeHash decodes an Argon2id hash string
func decodeHash(encodedHash string) (params *argon2Params, salt, hash []byte, err error) {
	// Check if hash format is correct
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported hash algorithm")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, errors.New("invalid hash version")
	}
	if version != 19 {
		return nil, nil, nil, errors.New("unsupported hash version")
	}

	// Parse the parameters
	params = &argon2Params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.time, &params.threads)
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash parameters")
	}

	// Decode the salt
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, errors.New("invalid salt encoding")
	}

	// Set the key length based on the decoded hash length
	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash encoding")
	}
	params.keyLen = uint32(len(hash))

	if params.keyLen == 0 {
		return nil, nil, nil, errors.New("invalid key length")
	}

	return params, salt, hash, nil
}
