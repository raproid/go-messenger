package main

import (
    "crypto/rand"
    "encoding/hex"
    "secure-messenger/crypto"
    "secure-messenger/shared"
    "secure-messenger/storage"
    "time"
)

type AuthManager struct {
    userStore *storage.UserStore
}

func NewAuthManager(userStore *storage.UserStore) *AuthManager {
    return &AuthManager{userStore: userStore}
}

func (am *AuthManager) Register(req *shared.RegisterRequest) (*shared.AuthResponse, error) {
    // Check if user already exists
    _, _, _, err := am.userStore.GetUserByUsername(req.Username)
    if err == nil {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Username already exists",
        }, nil
    }
    
    // Check if email already exists
    _, err = am.userStore.GetUserByEmail(req.Email)
    if err == nil {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Email already exists",
        }, nil
    }
    
    // Hash password
    passwordHash, passwordSalt, err := crypto.HashPassword(req.Password)
    if err != nil {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Failed to hash password",
        }, nil
    }
    
    // Create user
    user := &shared.User{
        ID:       generateID(),
        Username: req.Username,
        Email:    req.Email,
        Created:  time.Now(),
    }
    
    // Save user to database
    if err := am.userStore.CreateUser(user, passwordHash, passwordSalt); err != nil {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Failed to create user",
        }, nil
    }
    
    // Generate session token
    token := generateSessionToken()
    
    // Create session
    if err := am.userStore.CreateSession(token, user.ID); err != nil {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Failed to create session",
        }, nil
    }
    
    return &shared.AuthResponse{
        Success: true,
        Token:   token,
        User:    user,
    }, nil
}

func (am *AuthManager) Login(req *shared.LoginRequest) (*shared.AuthResponse, error) {
    // Get user from database
    user, passwordHash, passwordSalt, err := am.userStore.GetUserByUsername(req.Username)
    if err != nil {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Invalid username or password",
        }, nil
    }
    
    // Verify password
    if !crypto.VerifyPassword(req.Password, passwordHash, passwordSalt) {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Invalid username or password",
        }, nil
    }
    
    // Generate session token
    token := generateSessionToken()
    
    // Create session
    if err := am.userStore.CreateSession(token, user.ID); err != nil {
        return &shared.AuthResponse{
            Success: false,
            Error:   "Failed to create session",
        }, nil
    }
    
    return &shared.AuthResponse{
        Success: true,
        Token:   token,
        User:    user,
    }, nil
}

func (am *AuthManager) ValidateSession(token string) (*shared.User, error) {
    user, err := am.userStore.GetSession(token)
    if err != nil {
        return nil, err
    }
    
    // Update last seen
    am.userStore.UpdateSessionLastSeen(token)
    
    return user, nil
}

func (am *AuthManager) Logout(token string) error {
    return am.userStore.DeleteSession(token)
}

func generateID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func generateSessionToken() string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}
