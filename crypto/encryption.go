package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/json"
    "encoding/pem"
    "fmt"
    "io"
)

type EncryptionManager struct {
    keyManager *KeyManager
}

func NewEncryptionManager(keyManager *KeyManager) *EncryptionManager {
    return &EncryptionManager{keyManager: keyManager}
}

func (em *EncryptionManager) EncryptMessage(message string, recipientPublicKey *rsa.PublicKey) (string, error) {
    // Generate random AES key
    aesKey := make([]byte, 32)
    if _, err := rand.Read(aesKey); err != nil {
        return "", err
    }
    
    // Encrypt message with AES
    encryptedMessage, err := em.encryptAES(message, aesKey)
    if err != nil {
        return "", err
    }
    
    // Encrypt AES key with RSA
    encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, recipientPublicKey, aesKey, nil)
    if err != nil {
        return "", err
    }
    
    // Combine encrypted key and message
    result := map[string]string{
        "key":     base64.StdEncoding.EncodeToString(encryptedKey),
        "message": base64.StdEncoding.EncodeToString(encryptedMessage),
    }
    
    data, _ := json.Marshal(result)
    return base64.StdEncoding.EncodeToString(data), nil
}

func (em *EncryptionManager) DecryptMessage(encryptedData string) (string, error) {
    // Decode base64
    data, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return "", err
    }
    
    // Parse JSON
    var result map[string]string
    if err := json.Unmarshal(data, &result); err != nil {
        return "", err
    }
    
    // Decode encrypted key and message
    encryptedKey, err := base64.StdEncoding.DecodeString(result["key"])
    if err != nil {
        return "", err
    }
    
    encryptedMessage, err := base64.StdEncoding.DecodeString(result["message"])
    if err != nil {
        return "", err
    }
    
    // Decrypt AES key with RSA
    aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, em.keyManager.GetPrivateKey(), encryptedKey, nil)
    if err != nil {
        return "", err
    }
    
    // Decrypt message with AES
    return em.decryptAES(encryptedMessage, aesKey)
}

func (em *EncryptionManager) encryptAES(plaintext string, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return ciphertext, nil
}

func (em *EncryptionManager) decryptAES(ciphertext, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    
    return string(plaintext), nil
}

func (em *EncryptionManager) LoadPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(publicKeyPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to decode public key PEM")
    }
    
    publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    
    rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("not an RSA public key")
    }
    
    return rsaPublicKey, nil
}
