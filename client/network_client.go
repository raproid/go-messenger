package client

import (
    "crypto/tls"
    "crypto/x509"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net"
    "secure-messenger/shared"
    "time"
)

type NetworkClient struct {
    conn     net.Conn
    protocol *shared.Protocol
    Session  *Session
    config   *Config
}

type Session struct {
    Token    string
    User     *shared.User
    LastSeen time.Time
}

func NewNetworkClient() *NetworkClient {
    config, _ := LoadConfig()
    return &NetworkClient{config: config}
}

func (nc *NetworkClient) Connect() error {
    var conn net.Conn
    var err error
    
    if nc.config.UseTLS {
        // Load server certificate
        certPool := x509.NewCertPool()
        certData, err := ioutil.ReadFile(nc.config.CertPath)
        if err != nil {
            return fmt.Errorf("failed to read certificate: %v", err)
        }
        
        if !certPool.AppendCertsFromPEM(certData) {
            return fmt.Errorf("failed to parse certificate")
        }
        
        // Create TLS connection
        tlsConfig := &tls.Config{
            RootCAs:            certPool,
            ServerName:         nc.config.ServerName,
            InsecureSkipVerify: true, // Skip certificate verification for self-signed certs (development only)
        }
        
        conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", nc.config.ServerAddress, nc.config.ServerPort), tlsConfig)
        if err != nil {
            return fmt.Errorf("failed to connect to server: %v", err)
        }
    } else {
        // Plain TCP connection
        conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", nc.config.ServerAddress, nc.config.ServerPort))
        if err != nil {
            return fmt.Errorf("failed to connect to server: %v", err)
        }
    }
    
    nc.conn = conn
    nc.protocol = shared.NewProtocol(conn)
    
    return nil
}

func (nc *NetworkClient) Disconnect() error {
    if nc.conn != nil {
        return nc.conn.Close()
    }
    return nil
}

func (nc *NetworkClient) Register(username, email, password string) (*shared.AuthResponse, error) {
    req := &shared.RegisterRequest{
        Username: username,
        Email:    email,
        Password: password,
    }
    
    data, _ := json.Marshal(req)
    request := map[string]interface{}{
        "action": "register",
        "data":   string(data),
    }
    
    if err := nc.protocol.SendMessage(request); err != nil {
        return nil, err
    }
    
    msg, err := nc.protocol.ReadMessage()
    if err != nil {
        return nil, err
    }
    
    var response shared.AuthResponse
    responseData, _ := json.Marshal(msg)
    if err := json.Unmarshal(responseData, &response); err != nil {
        return nil, err
    }
    
    if response.Success {
        nc.Session = &Session{
            Token:    response.Token,
            User:     response.User,
            LastSeen: time.Now(),
        }
    }
    
    return &response, nil
}

func (nc *NetworkClient) Login(username, password string) (*shared.AuthResponse, error) {
    req := &shared.LoginRequest{
        Username: username,
        Password: password,
    }
    
    data, _ := json.Marshal(req)
    request := map[string]interface{}{
        "action": "login",
        "data":   string(data),
    }
    
    if err := nc.protocol.SendMessage(request); err != nil {
        return nil, err
    }
    
    msg, err := nc.protocol.ReadMessage()
    if err != nil {
        return nil, err
    }
    
    var response shared.AuthResponse
    responseData, _ := json.Marshal(msg)
    if err := json.Unmarshal(responseData, &response); err != nil {
        return nil, err
    }
    
    if response.Success {
        nc.Session = &Session{
            Token:    response.Token,
            User:     response.User,
            LastSeen: time.Now(),
        }
    }
    
    return &response, nil
}

func (nc *NetworkClient) SendMessage(to, content string) error {
    if nc.Session == nil {
        return fmt.Errorf("not authenticated")
    }
    
    req := &shared.MessageRequest{
        To:      to,
        Content: content,
    }
    
    data, _ := json.Marshal(req)
    request := map[string]interface{}{
        "action": "send_message",
        "token":  nc.Session.Token,
        "data":   string(data),
    }
    
    return nc.protocol.SendMessage(request)
}

func (nc *NetworkClient) SendChannelMessage(channelID, content string) error {
    if nc.Session == nil {
        return fmt.Errorf("not authenticated")
    }
    
    req := &shared.ChannelMessageRequest{
        ChannelID: channelID,
        Content:   content,
    }
    
    data, _ := json.Marshal(req)
    request := map[string]interface{}{
        "action": "send_channel_message",
        "token":  nc.Session.Token,
        "data":   string(data),
    }
    
    return nc.protocol.SendMessage(request)
}

func (nc *NetworkClient) GetMessages(otherUserID string, limit int) ([]*shared.Message, error) {
    if nc.Session == nil {
        return nil, fmt.Errorf("not authenticated")
    }
    
    request := map[string]interface{}{
        "action":        "get_messages",
        "token":         nc.Session.Token,
        "other_user_id": otherUserID,
        "limit":         limit,
    }
    
    if err := nc.protocol.SendMessage(request); err != nil {
        return nil, err
    }
    
    data, err := nc.protocol.ReadMessage()
    if err != nil {
        return nil, err
    }
    
    var response struct {
        Success  bool              `json:"success"`
        Messages []*shared.Message `json:"messages"`
        Error    string            `json:"error"`
    }
    
    responseData, _ := json.Marshal(data)
    if err := json.Unmarshal(responseData, &response); err != nil {
        return nil, err
    }
    
    if !response.Success {
        return nil, fmt.Errorf(response.Error)
    }
    
    return response.Messages, nil
}

func (nc *NetworkClient) GetChannelMessages(channelID string, limit int) ([]*shared.Message, error) {
    if nc.Session == nil {
        return nil, fmt.Errorf("not authenticated")
    }
    
    request := map[string]interface{}{
        "action":     "get_channel_messages",
        "token":      nc.Session.Token,
        "channel_id": channelID,
        "limit":      limit,
    }
    
    if err := nc.protocol.SendMessage(request); err != nil {
        return nil, err
    }
    
    data, err := nc.protocol.ReadMessage()
    if err != nil {
        return nil, err
    }
    
    var response struct {
        Success  bool              `json:"success"`
        Messages []*shared.Message `json:"messages"`
        Error    string            `json:"error"`
    }
    
    responseData, _ := json.Marshal(data)
    if err := json.Unmarshal(responseData, &response); err != nil {
        return nil, err
    }
    
    if !response.Success {
        return nil, fmt.Errorf(response.Error)
    }
    
    return response.Messages, nil
}

func (nc *NetworkClient) CreateChannel(name, description string, members []string) (*shared.Channel, error) {
    if nc.Session == nil {
        return nil, fmt.Errorf("not authenticated")
    }
    
    req := &shared.ChannelRequest{
        Name:        name,
        Description: description,
        Members:     members,
    }
    
    data, _ := json.Marshal(req)
    request := map[string]interface{}{
        "action": "create_channel",
        "token":  nc.Session.Token,
        "data":   string(data),
    }
    
    if err := nc.protocol.SendMessage(request); err != nil {
        return nil, err
    }
    
    msg, err := nc.protocol.ReadMessage()
    if err != nil {
        return nil, err
    }
    
    var response struct {
        Success bool             `json:"success"`
        Channel *shared.Channel `json:"channel"`
        Error   string           `json:"error"`
    }
    
    responseData, _ := json.Marshal(msg)
    if err := json.Unmarshal(responseData, &response); err != nil {
        return nil, err
    }
    
    if !response.Success {
        return nil, fmt.Errorf(response.Error)
    }
    
    return response.Channel, nil
}

func (nc *NetworkClient) GetUserChannels() ([]*shared.Channel, error) {
    if nc.Session == nil {
        return nil, fmt.Errorf("not authenticated")
    }
    
    request := map[string]interface{}{
        "action": "get_user_channels",
        "token":  nc.Session.Token,
    }
    
    if err := nc.protocol.SendMessage(request); err != nil {
        return nil, err
    }
    
    data, err := nc.protocol.ReadMessage()
    if err != nil {
        return nil, err
    }
    
    var response struct {
        Success  bool              `json:"success"`
        Channels []*shared.Channel `json:"channels"`
        Error    string            `json:"error"`
    }
    
    responseData, _ := json.Marshal(data)
    if err := json.Unmarshal(responseData, &response); err != nil {
        return nil, err
    }
    
    if !response.Success {
        return nil, fmt.Errorf(response.Error)
    }
    
    return response.Channels, nil
}

func (nc *NetworkClient) IsAuthenticated() bool {
    return nc.Session != nil
}

func (nc *NetworkClient) GetUser() *shared.User {
    if nc.Session != nil {
        return nc.Session.User
    }
    return nil
}
