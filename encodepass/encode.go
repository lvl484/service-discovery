package encodepass

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	configTime    = 3
	configMemory  = 64 * 1024
	configThreads = 1
	configKeyLen  = 64
)

// EncodePassword returns encoded password
func EncodePassword(c *PasswordConfig, pass string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	passByte := []byte(pass)

	hash := argon2.IDKey(passByte, salt, c.time, c.memory, c.threads, c.keyLen)
	log.Println(hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)

	return full, nil
}

// ComparePassword returns true if password matches
func ComparePassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	c := &PasswordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)

	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])

	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])

	if err != nil {
		return false, err
	}

	c.keyLen = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}

// CompareEncodePassword returns true if password matches
func CompareEncodedPassword(pass1, pass2 string) bool {
	return (pass1 == pass2)
}

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// NewPasswordConfig returns config for encode
func NewPasswordConfig() *PasswordConfig {
	return &PasswordConfig{
		time:    configTime,
		memory:  configMemory,
		threads: configThreads,
		keyLen:  configKeyLen,
	}
}
