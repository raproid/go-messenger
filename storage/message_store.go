package storage

import (
    "database/sql"
    "fmt"
    "secure-messenger/shared"
    "time"
)

type MessageStore struct {
    db *sql.DB
}

func NewMessageStore(db *sql.DB) *MessageStore {
    return &MessageStore{db: db}
}

func (ms *MessageStore) CreateMessage(message *shared.Message) error {
    query := `
    INSERT INTO messages (id, from_user, to_user, channel_id, content, encrypted, timestamp)
    VALUES (?, ?, ?, ?, ?, ?, ?)`
    
    _, err := ms.db.Exec(query, message.ID, message.From, message.To, message.ChannelID, message.Content, message.Encrypted, message.Timestamp)
    return err
}

func (ms *MessageStore) GetMessagesBetweenUsers(user1ID, user2ID string, limit int) ([]*shared.Message, error) {
    query := `
    SELECT id, from_user, to_user, content, encrypted, timestamp
    FROM messages
    WHERE (from_user = ? AND to_user = ?) OR (from_user = ? AND to_user = ?)
    ORDER BY timestamp DESC
    LIMIT ?`
    
    rows, err := ms.db.Query(query, user1ID, user2ID, user2ID, user1ID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var messages []*shared.Message
    for rows.Next() {
        var msg shared.Message
        err := rows.Scan(&msg.ID, &msg.From, &msg.To, &msg.Content, &msg.Encrypted, &msg.Timestamp)
        if err != nil {
            return nil, err
        }
        messages = append(messages, &msg)
    }
    
    return messages, nil
}

func (ms *MessageStore) GetChannelMessages(channelID string, limit int) ([]*shared.Message, error) {
    query := `
    SELECT id, from_user, to_user, content, encrypted, timestamp
    FROM messages
    WHERE channel_id = ?
    ORDER BY timestamp DESC
    LIMIT ?`
    
    rows, err := ms.db.Query(query, channelID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var messages []*shared.Message
    for rows.Next() {
        var msg shared.Message
        err := rows.Scan(&msg.ID, &msg.From, &msg.To, &msg.Content, &msg.Encrypted, &msg.Timestamp)
        if err != nil {
            return nil, err
        }
        messages = append(messages, &msg)
    }
    
    return messages, nil
}

func (ms *MessageStore) CreateChannel(channel *shared.Channel) error {
    // Start transaction
    tx, err := ms.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Create channel
    query := `
    INSERT INTO channels (id, name, description, created_by, created_at)
    VALUES (?, ?, ?, ?, ?)`
    
    _, err = tx.Exec(query, channel.ID, channel.Name, channel.Description, channel.CreatedBy, channel.Created)
    if err != nil {
        return err
    }
    
    // Add members
    for _, memberID := range channel.Members {
        memberQuery := `
        INSERT INTO channel_members (channel_id, user_id, joined_at)
        VALUES (?, ?, ?)`
        
        _, err = tx.Exec(memberQuery, channel.ID, memberID, time.Now())
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (ms *MessageStore) GetChannel(channelID string) (*shared.Channel, error) {
    var channel shared.Channel
    
    query := `
    SELECT id, name, description, created_by, created_at
    FROM channels WHERE id = ?`
    
    row := ms.db.QueryRow(query, channelID)
    err := row.Scan(&channel.ID, &channel.Name, &channel.Description, &channel.CreatedBy, &channel.Created)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("channel not found")
        }
        return nil, err
    }
    
    // Get members
    membersQuery := `
    SELECT user_id FROM channel_members WHERE channel_id = ?`
    
    rows, err := ms.db.Query(membersQuery, channelID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var members []string
    for rows.Next() {
        var memberID string
        err := rows.Scan(&memberID)
        if err != nil {
            return nil, err
        }
        members = append(members, memberID)
    }
    
    channel.Members = members
    return &channel, nil
}

func (ms *MessageStore) GetUserChannels(userID string) ([]*shared.Channel, error) {
    query := `
    SELECT c.id, c.name, c.description, c.created_by, c.created_at
    FROM channels c
    JOIN channel_members cm ON c.id = cm.channel_id
    WHERE cm.user_id = ?
    ORDER BY c.created_at DESC`
    
    rows, err := ms.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var channels []*shared.Channel
    for rows.Next() {
        var channel shared.Channel
        err := rows.Scan(&channel.ID, &channel.Name, &channel.Description, &channel.CreatedBy, &channel.Created)
        if err != nil {
            return nil, err
        }
        
        // Get members for this channel
        membersQuery := `
        SELECT user_id FROM channel_members WHERE channel_id = ?`
        
        memberRows, err := ms.db.Query(membersQuery, channel.ID)
        if err != nil {
            return nil, err
        }
        
        var members []string
        for memberRows.Next() {
            var memberID string
            err := memberRows.Scan(&memberID)
            if err != nil {
                memberRows.Close()
                return nil, err
            }
            members = append(members, memberID)
        }
        memberRows.Close()
        
        channel.Members = members
        channels = append(channels, &channel)
    }
    
    return channels, nil
}

func (ms *MessageStore) AddUserToChannel(channelID, userID string) error {
    query := `
    INSERT INTO channel_members (channel_id, user_id, joined_at)
    VALUES (?, ?, ?)`
    
    _, err := ms.db.Exec(query, channelID, userID, time.Now())
    return err
}

func (ms *MessageStore) RemoveUserFromChannel(channelID, userID string) error {
    query := `DELETE FROM channel_members WHERE channel_id = ? AND user_id = ?`
    _, err := ms.db.Exec(query, channelID, userID)
    return err
}

func (ms *MessageStore) GetRecentMessages(userID string, limit int) ([]*shared.Message, error) {
    query := `
    SELECT m.id, m.from_user, m.to_user, m.channel_id, m.content, m.encrypted, m.timestamp
    FROM messages m
    WHERE m.from_user = ? OR m.to_user = ? OR m.channel_id IN (
        SELECT channel_id FROM channel_members WHERE user_id = ?
    )
    ORDER BY m.timestamp DESC
    LIMIT ?`
    
    rows, err := ms.db.Query(query, userID, userID, userID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var messages []*shared.Message
    for rows.Next() {
        var msg shared.Message
        err := rows.Scan(&msg.ID, &msg.From, &msg.To, &msg.ChannelID, &msg.Content, &msg.Encrypted, &msg.Timestamp)
        if err != nil {
            return nil, err
        }
        messages = append(messages, &msg)
    }
    
    return messages, nil
}
