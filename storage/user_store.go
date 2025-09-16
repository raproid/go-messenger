package storage

import (
    "database/sql"
    "fmt"
    "secure-messenger/shared"
    "time"
)

type UserStore struct {
    db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
    return &UserStore{db: db}
}

func (us *UserStore) CreateUser(user *shared.User, passwordHash, passwordSalt string) error {
    query := `
    INSERT INTO users (id, username, email, password_hash, password_salt, created_at)
    VALUES (?, ?, ?, ?, ?, ?)`
    
    _, err := us.db.Exec(query, user.ID, user.Username, user.Email, passwordHash, passwordSalt, user.Created)
    return err
}

func (us *UserStore) GetUserByUsername(username string) (*shared.User, string, string, error) {
    var user shared.User
    var passwordHash, passwordSalt string
    
    query := `
    SELECT id, username, email, password_hash, password_salt, created_at
    FROM users WHERE username = ?`
    
    row := us.db.QueryRow(query, username)
    err := row.Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &passwordSalt, &user.Created)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, "", "", fmt.Errorf("user not found")
        }
        return nil, "", "", err
    }
    
    return &user, passwordHash, passwordSalt, nil
}

func (us *UserStore) GetUserByID(id string) (*shared.User, error) {
    var user shared.User
    
    query := `
    SELECT id, username, email, created_at
    FROM users WHERE id = ?`
    
    row := us.db.QueryRow(query, id)
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Created)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}

func (us *UserStore) GetUserByEmail(email string) (*shared.User, error) {
    var user shared.User
    
    query := `
    SELECT id, username, email, created_at
    FROM users WHERE email = ?`
    
    row := us.db.QueryRow(query, email)
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Created)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}

func (us *UserStore) UpdateUserPublicKey(userID, publicKey string) error {
    query := `UPDATE users SET public_key = ? WHERE id = ?`
    _, err := us.db.Exec(query, publicKey, userID)
    return err
}

func (us *UserStore) GetUserPublicKey(userID string) (string, error) {
    var publicKey string
    
    query := `SELECT public_key FROM users WHERE id = ?`
    row := us.db.QueryRow(query, userID)
    err := row.Scan(&publicKey)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("user not found")
        }
        return "", err
    }
    
    return publicKey, nil
}

func (us *UserStore) CreateSession(token, userID string) error {
    query := `
    INSERT INTO sessions (token, user_id, created_at, last_seen)
    VALUES (?, ?, ?, ?)`
    
    now := time.Now()
    _, err := us.db.Exec(query, token, userID, now, now)
    return err
}

func (us *UserStore) GetSession(token string) (*shared.User, error) {
    var user shared.User
    
    query := `
    SELECT u.id, u.username, u.email, u.created_at
    FROM users u
    JOIN sessions s ON u.id = s.user_id
    WHERE s.token = ?`
    
    row := us.db.QueryRow(query, token)
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Created)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("session not found")
        }
        return nil, err
    }
    
    return &user, nil
}

func (us *UserStore) UpdateSessionLastSeen(token string) error {
    query := `UPDATE sessions SET last_seen = ? WHERE token = ?`
    _, err := us.db.Exec(query, time.Now(), token)
    return err
}

func (us *UserStore) DeleteSession(token string) error {
    query := `DELETE FROM sessions WHERE token = ?`
    _, err := us.db.Exec(query, token)
    return err
}

func (us *UserStore) GetAllUsers() ([]*shared.User, error) {
    query := `
    SELECT id, username, email, created_at
    FROM users
    ORDER BY created_at DESC`
    
    rows, err := us.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*shared.User
    for rows.Next() {
        var user shared.User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Created)
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    
    return users, nil
}
