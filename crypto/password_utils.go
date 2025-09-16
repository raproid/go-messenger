package crypto

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "golang.org/x/crypto/pbkdf2"
)

const (
    SaltLength = 32
    Iterations = 100000
)

func HashPassword(password string) (string, string, error) {
    // Generate random salt
    salt := make([]byte, SaltLength)
    if _, err := rand.Read(salt); err != nil {
        return "", "", err
    }
    
    // Hash password with salt
    hash := pbkdf2.Key([]byte(password), salt, Iterations, 32, sha256.New)
    
    // Encode salt and hash
    saltStr := base64.StdEncoding.EncodeToString(salt)
    hashStr := base64.StdEncoding.EncodeToString(hash)
    
    return hashStr, saltStr, nil
}

func VerifyPassword(password, hash, salt string) bool {
    // Decode salt
    saltBytes, err := base64.StdEncoding.DecodeString(salt)
    if err != nil {
        return false
    }
    
    // Decode hash
    hashBytes, err := base64.StdEncoding.DecodeString(hash)
    if err != nil {
        return false
    }
    
    // Hash password with same salt
    testHash := pbkdf2.Key([]byte(password), saltBytes, Iterations, 32, sha256.New)
    
    // Compare hashes
    return fmt.Sprintf("%x", testHash) == fmt.Sprintf("%x", hashBytes)
}
