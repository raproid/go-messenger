package main

import (
    "encoding/json"
    "log"
    "net"
    "secure-messenger/shared"
    "secure-messenger/storage"
)

type Server struct {
    db           *storage.Database
    userStore    *storage.UserStore
    messageStore *storage.MessageStore
    authManager  *AuthManager
    messageHandler *MessageHandler
}

func NewServer(db *storage.Database) *Server {
    userStore := storage.NewUserStore(db.GetDB())
    messageStore := storage.NewMessageStore(db.GetDB())
    
    return &Server{
        db:            db,
        userStore:     userStore,
        messageStore:  messageStore,
        authManager:   NewAuthManager(userStore),
        messageHandler: NewMessageHandler(messageStore, userStore),
    }
}

func (s *Server) HandleConnection(conn net.Conn) {
    defer conn.Close()
    
    protocol := shared.NewProtocol(conn)
    
    for {
        // Read message from client
        msg, err := protocol.ReadMessage()
        if err != nil {
            log.Printf("Failed to read message: %v", err)
            return
        }
        
        // Process message
        response, err := s.processMessage(msg, protocol)
        if err != nil {
            log.Printf("Failed to process message: %v", err)
            continue
        }
        
        // Send response
        if response != nil {
            if err := protocol.SendMessage(response); err != nil {
                log.Printf("Failed to send response: %v", err)
                return
            }
        }
    }
}

func (s *Server) processMessage(msg map[string]interface{}, protocol *shared.Protocol) (map[string]interface{}, error) {
    action, ok := msg["action"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid action",
        }, nil
    }
    
    switch action {
    case "register":
        return s.handleRegister(msg)
    case "login":
        return s.handleLogin(msg)
    case "send_message":
        return s.handleSendMessage(msg)
    case "send_channel_message":
        return s.handleSendChannelMessage(msg)
    case "get_messages":
        return s.handleGetMessages(msg)
    case "get_channel_messages":
        return s.handleGetChannelMessages(msg)
    case "create_channel":
        return s.handleCreateChannel(msg)
    case "get_user_channels":
        return s.handleGetUserChannels(msg)
    case "add_user_to_channel":
        return s.handleAddUserToChannel(msg)
    case "remove_user_from_channel":
        return s.handleRemoveUserFromChannel(msg)
    case "get_recent_messages":
        return s.handleGetRecentMessages(msg)
    default:
        return map[string]interface{}{
            "success": false,
            "error":   "Unknown action",
        }, nil
    }
}

func (s *Server) handleRegister(msg map[string]interface{}) (map[string]interface{}, error) {
    data, ok := msg["data"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request data",
        }, nil
    }
    
    var req shared.RegisterRequest
    if err := json.Unmarshal([]byte(data), &req); err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request format",
        }, nil
    }
    
    response, err := s.authManager.Register(&req)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success": response.Success,
        "token":   response.Token,
        "user":    response.User,
        "error":   response.Error,
    }, nil
}

func (s *Server) handleLogin(msg map[string]interface{}) (map[string]interface{}, error) {
    data, ok := msg["data"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request data",
        }, nil
    }
    
    var req shared.LoginRequest
    if err := json.Unmarshal([]byte(data), &req); err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request format",
        }, nil
    }
    
    response, err := s.authManager.Login(&req)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success": response.Success,
        "token":   response.Token,
        "user":    response.User,
        "error":   response.Error,
    }, nil
}

func (s *Server) handleSendMessage(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    user, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    data, ok := msg["data"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request data",
        }, nil
    }
    
    var req shared.MessageRequest
    if err := json.Unmarshal([]byte(data), &req); err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request format",
        }, nil
    }
    
    message, err := s.messageHandler.SendMessage(&req, user.ID)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success": true,
        "message": message,
    }, nil
}

func (s *Server) handleSendChannelMessage(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    user, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    data, ok := msg["data"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request data",
        }, nil
    }
    
    var req shared.ChannelMessageRequest
    if err := json.Unmarshal([]byte(data), &req); err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request format",
        }, nil
    }
    
    message, err := s.messageHandler.SendChannelMessage(&req, user.ID)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success": true,
        "message": message,
    }, nil
}

func (s *Server) handleGetMessages(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    user, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    otherUserID, ok := msg["other_user_id"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Other user ID required",
        }, nil
    }
    
    limit := 50
    if l, ok := msg["limit"].(float64); ok {
        limit = int(l)
    }
    
    messages, err := s.messageHandler.GetMessages(user.ID, otherUserID, limit)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success":  true,
        "messages": messages,
    }, nil
}

func (s *Server) handleGetChannelMessages(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    _, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    channelID, ok := msg["channel_id"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Channel ID required",
        }, nil
    }
    
    limit := 50
    if l, ok := msg["limit"].(float64); ok {
        limit = int(l)
    }
    
    messages, err := s.messageHandler.GetChannelMessages(channelID, limit)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success":  true,
        "messages": messages,
    }, nil
}

func (s *Server) handleCreateChannel(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    user, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    data, ok := msg["data"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request data",
        }, nil
    }
    
    var req shared.ChannelRequest
    if err := json.Unmarshal([]byte(data), &req); err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid request format",
        }, nil
    }
    
    channel, err := s.messageHandler.CreateChannel(&req, user.ID)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success": true,
        "channel": channel,
    }, nil
}

func (s *Server) handleGetUserChannels(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    user, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    channels, err := s.messageHandler.GetUserChannels(user.ID)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success":  true,
        "channels": channels,
    }, nil
}

func (s *Server) handleAddUserToChannel(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    _, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    channelID, ok := msg["channel_id"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Channel ID required",
        }, nil
    }
    
    userID, ok := msg["user_id"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "User ID required",
        }, nil
    }
    
    if err := s.messageHandler.AddUserToChannel(channelID, userID); err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success": true,
    }, nil
}

func (s *Server) handleRemoveUserFromChannel(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    _, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    channelID, ok := msg["channel_id"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Channel ID required",
        }, nil
    }
    
    userID, ok := msg["user_id"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "User ID required",
        }, nil
    }
    
    if err := s.messageHandler.RemoveUserFromChannel(channelID, userID); err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success": true,
    }, nil
}

func (s *Server) handleGetRecentMessages(msg map[string]interface{}) (map[string]interface{}, error) {
    token, ok := msg["token"].(string)
    if !ok {
        return map[string]interface{}{
            "success": false,
            "error":   "Authentication required",
        }, nil
    }
    
    // Validate session
    user, err := s.authManager.ValidateSession(token)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   "Invalid session",
        }, nil
    }
    
    limit := 50
    if l, ok := msg["limit"].(float64); ok {
        limit = int(l)
    }
    
    messages, err := s.messageHandler.GetRecentMessages(user.ID, limit)
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "success":  true,
        "messages": messages,
    }, nil
}
