package main

import (
    "fmt"
    "secure-messenger/shared"
    "secure-messenger/storage"
    "time"
)

type MessageHandler struct {
    messageStore *storage.MessageStore
    userStore    *storage.UserStore
}

func NewMessageHandler(messageStore *storage.MessageStore, userStore *storage.UserStore) *MessageHandler {
    return &MessageHandler{
        messageStore: messageStore,
        userStore:    userStore,
    }
}

func (mh *MessageHandler) SendMessage(req *shared.MessageRequest, fromUserID string) (*shared.Message, error) {
    // Create message
    message := &shared.Message{
        ID:        generateMessageID(),
        From:      fromUserID,
        To:        req.To,
        Content:   req.Content,
        Encrypted: true,
        Timestamp: time.Now(),
    }
    
    // Save message to database
    if err := mh.messageStore.CreateMessage(message); err != nil {
        return nil, fmt.Errorf("failed to save message: %v", err)
    }
    
    return message, nil
}

func (mh *MessageHandler) SendChannelMessage(req *shared.ChannelMessageRequest, fromUserID string) (*shared.Message, error) {
    // Verify user is member of channel
    channels, err := mh.messageStore.GetUserChannels(fromUserID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user channels: %v", err)
    }
    
    isMember := false
    for _, channel := range channels {
        if channel.ID == req.ChannelID {
            isMember = true
            break
        }
    }
    
    if !isMember {
        return nil, fmt.Errorf("user is not a member of this channel")
    }
    
    // Create message
    message := &shared.Message{
        ID:        generateMessageID(),
        From:      fromUserID,
        ChannelID: req.ChannelID,
        Content:   req.Content,
        Encrypted: true,
        Timestamp: time.Now(),
    }
    
    // Save message to database
    if err := mh.messageStore.CreateMessage(message); err != nil {
        return nil, fmt.Errorf("failed to save message: %v", err)
    }
    
    return message, nil
}

func (mh *MessageHandler) GetMessages(userID, otherUserID string, limit int) ([]*shared.Message, error) {
    return mh.messageStore.GetMessagesBetweenUsers(userID, otherUserID, limit)
}

func (mh *MessageHandler) GetChannelMessages(channelID string, limit int) ([]*shared.Message, error) {
    return mh.messageStore.GetChannelMessages(channelID, limit)
}

func (mh *MessageHandler) CreateChannel(req *shared.ChannelRequest, creatorID string) (*shared.Channel, error) {
    // Create channel
    channel := &shared.Channel{
        ID:          generateChannelID(),
        Name:        req.Name,
        Description: req.Description,
        Members:     append(req.Members, creatorID), // Add creator to members
        Created:     time.Now(),
        CreatedBy:   creatorID,
    }
    
    // Save channel to database
    if err := mh.messageStore.CreateChannel(channel); err != nil {
        return nil, fmt.Errorf("failed to create channel: %v", err)
    }
    
    return channel, nil
}

func (mh *MessageHandler) GetUserChannels(userID string) ([]*shared.Channel, error) {
    return mh.messageStore.GetUserChannels(userID)
}

func (mh *MessageHandler) AddUserToChannel(channelID, userID string) error {
    return mh.messageStore.AddUserToChannel(channelID, userID)
}

func (mh *MessageHandler) RemoveUserFromChannel(channelID, userID string) error {
    return mh.messageStore.RemoveUserFromChannel(channelID, userID)
}

func (mh *MessageHandler) GetRecentMessages(userID string, limit int) ([]*shared.Message, error) {
    return mh.messageStore.GetRecentMessages(userID, limit)
}

func generateMessageID() string {
    return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

func generateChannelID() string {
    return fmt.Sprintf("ch_%d", time.Now().UnixNano())
}
