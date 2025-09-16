package client

import (
    "encoding/json"
    "os"
    "secure-messenger/shared"
)

type MessageHandler struct {
    messagePath string
}

func NewMessageHandler() *MessageHandler {
    // Use project-relative path
    messagePath := "messages.json"
    return &MessageHandler{messagePath: messagePath}
}

func (mh *MessageHandler) SaveMessage(message *shared.Message) error {
    messages, err := mh.LoadMessages()
    if err != nil {
        messages = []*shared.Message{}
    }
    
    messages = append(messages, message)
    
    data, err := json.Marshal(messages)
    if err != nil {
        return err
    }
    
    // No need to create directory for project-relative path
    
    return os.WriteFile(mh.messagePath, data, 0644)
}

func (mh *MessageHandler) LoadMessages() ([]*shared.Message, error) {
    data, err := os.ReadFile(mh.messagePath)
    if err != nil {
        if os.IsNotExist(err) {
            return []*shared.Message{}, nil
        }
        return nil, err
    }
    
    var messages []*shared.Message
    if err := json.Unmarshal(data, &messages); err != nil {
        return nil, err
    }
    
    return messages, nil
}

func (mh *MessageHandler) GetMessagesWithUser(userID string) ([]*shared.Message, error) {
    messages, err := mh.LoadMessages()
    if err != nil {
        return nil, err
    }
    
    var filtered []*shared.Message
    for _, msg := range messages {
        if (msg.From == userID && msg.To != "") || (msg.To == userID && msg.From != "") {
            filtered = append(filtered, msg)
        }
    }
    
    return filtered, nil
}

func (mh *MessageHandler) GetChannelMessages(channelID string) ([]*shared.Message, error) {
    messages, err := mh.LoadMessages()
    if err != nil {
        return nil, err
    }
    
    var filtered []*shared.Message
    for _, msg := range messages {
        if msg.ChannelID == channelID {
            filtered = append(filtered, msg)
        }
    }
    
    return filtered, nil
}

func (mh *MessageHandler) ClearMessages() error {
    return os.Remove(mh.messagePath)
}
