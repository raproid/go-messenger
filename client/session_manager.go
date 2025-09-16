package client

import (
    "encoding/json"
    "os"
    "time"
)

type SessionManager struct {
    sessionPath string
}

func NewSessionManager() *SessionManager {
    // Use project-relative path
    sessionPath := "session.json"
    return &SessionManager{sessionPath: sessionPath}
}

func (sm *SessionManager) SaveSession(session *Session) error {
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }
    
    // No need to create directory for project-relative path
    
    return os.WriteFile(sm.sessionPath, data, 0600)
}

func (sm *SessionManager) LoadSession() (*Session, error) {
    data, err := os.ReadFile(sm.sessionPath)
    if err != nil {
        return nil, err
    }
    
    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, err
    }
    
    return &session, nil
}

func (sm *SessionManager) ClearSession() error {
    return os.Remove(sm.sessionPath)
}

func (sm *SessionManager) HasValidSession() bool {
    session, err := sm.LoadSession()
    if err != nil {
        return false
    }
    
    // Check if session is not too old (e.g., 7 days)
    if time.Since(session.LastSeen) > 7*24*time.Hour {
        sm.ClearSession()
        return false
    }
    
    return true
}
