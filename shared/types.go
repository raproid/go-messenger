package shared

import (
    "time"
)

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Created  time.Time `json:"created"`
}

type Message struct {
    ID        string    `json:"id"`
    From      string    `json:"from"`
    To        string    `json:"to"`
    ChannelID string    `json:"channel_id"`
    Content   string    `json:"content"`
    Encrypted bool      `json:"encrypted"`
    Timestamp time.Time `json:"timestamp"`
}

type Channel struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Members     []string `json:"members"`
    Created     time.Time `json:"created"`
    CreatedBy   string   `json:"created_by"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Success bool   `json:"success"`
    Token   string `json:"token"`
    User    *User  `json:"user"`
    Error   string `json:"error"`
}

type MessageRequest struct {
    To      string `json:"to"`
    Content string `json:"content"`
}

type ChannelMessageRequest struct {
    ChannelID string `json:"channel_id"`
    Content   string `json:"content"`
}

type ChannelRequest struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Members     []string `json:"members"`
}
