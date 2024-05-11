package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltSize     = 12
	iter         = 20000
)

func getRandomSalt(size int) []byte {
	salt := make([]byte, size)
	l := len(allowedChars)
	for i := range salt {
		salt[i] = allowedChars[rand.Intn(l)]
	}
	return salt
}

func CheckPbkdf2SHA256(password, encoded string) (bool, error) {
	parts := strings.SplitN(encoded, "$", 4)
	if len(parts) != 4 {
		return false, errors.New("Hash must consist of 4 segments")
	}
	iter, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, fmt.Errorf("Wrong number of iterations: %v", err)
	}
	salt := []byte(parts[2])
	k, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, fmt.Errorf("Wrong hash encoding: %v", err)
	}
	dk := pbkdf2.Key([]byte(password), salt, iter, sha256.Size, sha256.New)
	return bytes.Equal(k, dk), nil
}

func MakePbkdf2SHA256(password string) string {
	salt := getRandomSalt(saltSize)
	dk := pbkdf2.Key([]byte(password), salt, iter, sha256.Size, sha256.New)
	b64Hash := base64.StdEncoding.EncodeToString(dk)
	return fmt.Sprintf("%s$%d$%s$%s", "pbkdf2_sha256", iter, salt, b64Hash)
}
