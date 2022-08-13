package security

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

// saltSize
// Define salt size
const saltSize = 16

type Encrypt struct {
	Password string
	Salt     string
}

func EncryptPassword(password string) *Encrypt {
	saltByte := generateRandomSalt(saltSize)
	return &Encrypt{Password: hashPassword(password, saltByte), Salt: base64.URLEncoding.EncodeToString(saltByte)}
}

// generateRandomSalt
// Generate 16 bytes randomly and securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func generateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}

// hashPassword
// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a base64 encoded string
func hashPassword(password string, salt []byte) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a base64 encoded string
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

// ComparePassword
// Check if two passwords compare
func ComparePassword(hashedPassword, currPassword string, saltStr string) bool {
	salt, err := base64.URLEncoding.DecodeString(saltStr)
	if err != nil {
		return false
	}
	var currPasswordHash = hashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}
