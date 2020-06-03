package pwhash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// PasswordHasher provides methods to calculate password hash and to check if
// password is valid
type PasswordHasher struct {
	saltSize   int
	iterations int
}

// New returns a new password hasher
func New(saltSize, iterations int) (*PasswordHasher, error) {
	if saltSize == 0 || iterations < 1000 {
		return nil, fmt.Errorf("saltSize or the number of iterations are too low")
	}
	return &PasswordHasher{
		saltSize:   saltSize,
		iterations: iterations,
	}, nil
}

// HashPassword generates random salt and calculates password hash with that
// random salt
func (ph *PasswordHasher) HashPassword(password string) (string, error) {
	// Generate a random salt value, 128 bits.
	salt := make([]byte, ph.saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("couldn't get random data: %v", err)
	}

	return hashWithSalt(password, salt, ph.iterations), nil
}

func hashWithSalt(password string, salt []byte, iterations int) string {
	// Generate the hash.  This should be a little painful, adjust ITERATIONS
	// if it needs performance tweeking.  Greatly depends on the hardware.
	// NOTE: We store these details with the returned hash, so changes will not
	// affect our ability to do password compares.
	hash := pbkdf2.Key([]byte(password), salt, iterations, sha512.Size, sha512.New)

	// Build up the parameters and hash into a single string so we can compare
	// other string to the same hash.  Note that the hash algorithm is hard-
	// coded here, as it is above.  Introducing alternate encodings must support
	// old encodings as well, and build this string appropriately.
	var buffer bytes.Buffer

	buffer.WriteString("PBKDF2$")
	buffer.WriteString("sha512$")
	buffer.WriteString(strconv.Itoa(iterations))
	buffer.WriteString("$")
	buffer.WriteString(base64.StdEncoding.EncodeToString(salt))
	buffer.WriteString("$")
	buffer.WriteString(base64.StdEncoding.EncodeToString(hash))

	return buffer.String()
}

// Validate ensures that the password matches the hash, if not it returns an error
func (ph *PasswordHasher) Validate(password, hash string) error {
	hashSplit := strings.Split(hash, "$")
	if len(hashSplit) != 5 {
		return fmt.Errorf("invalid password hash")
	}

	iterations, err := strconv.Atoi(hashSplit[2])
	if err != nil {
		return fmt.Errorf("invalid password hash")
	}
	salt, err := base64.StdEncoding.DecodeString(hashSplit[3])
	if err != nil {
		return fmt.Errorf("invalid password hash")
	}
	newHash := hashWithSalt(password, salt, iterations)
	if subtle.ConstantTimeCompare([]byte(newHash), []byte(hash)) != 1 {
		return fmt.Errorf("invalid password")
	}
	return nil
}
