package crypto

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "os"
    "path/filepath"
)

const (
    KeySize = 2048
)

type KeyManager struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
}

func NewKeyManager() *KeyManager {
    return &KeyManager{}
}

func (km *KeyManager) GenerateKeys() error {
    privateKey, err := rsa.GenerateKey(rand.Reader, KeySize)
    if err != nil {
        return err
    }
    
    km.privateKey = privateKey
    km.publicKey = &privateKey.PublicKey
    
    return nil
}

func (km *KeyManager) LoadKeys(keyPath string) error {
    // Load private key
    privateKeyPath := filepath.Join(keyPath, "private.pem")
    privateKeyData, err := os.ReadFile(privateKeyPath)
    if err != nil {
        return err
    }
    
    block, _ := pem.Decode(privateKeyData)
    if block == nil {
        return fmt.Errorf("failed to decode private key PEM")
    }
    
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return err
    }
    
    km.privateKey = privateKey
    km.publicKey = &privateKey.PublicKey
    
    return nil
}

func (km *KeyManager) SaveKeys(keyPath string) error {
    if err := os.MkdirAll(keyPath, 0700); err != nil {
        return err
    }
    
    // Save private key
    privateKeyPath := filepath.Join(keyPath, "private.pem")
    privateKeyData := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(km.privateKey),
    })
    
    if err := os.WriteFile(privateKeyPath, privateKeyData, 0600); err != nil {
        return err
    }
    
    // Save public key
    publicKeyPath := filepath.Join(keyPath, "public.pem")
    publicKeyBytes, err := x509.MarshalPKIXPublicKey(km.publicKey)
    if err != nil {
        return err
    }
    publicKeyData := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicKeyBytes,
    })
    
    if err := os.WriteFile(publicKeyPath, publicKeyData, 0644); err != nil {
        return err
    }
    
    return nil
}

func (km *KeyManager) GetPublicKey() *rsa.PublicKey {
    return km.publicKey
}

func (km *KeyManager) GetPrivateKey() *rsa.PrivateKey {
    return km.privateKey
}
