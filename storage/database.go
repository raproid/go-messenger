package storage

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    _ "github.com/mattn/go-sqlite3"
)

type Database struct {
    db *sql.DB
}

func NewDatabase(dataDir string) (*Database, error) {
    // Create data directory if it doesn't exist
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        return nil, err
    }
    
    // Open SQLite database
    dbPath := filepath.Join(dataDir, "messenger.db")
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    database := &Database{db: db}
    
    // Create tables
    if err := database.createTables(); err != nil {
        return nil, err
    }
    
    return database, nil
}

func (d *Database) createTables() error {
    // Users table
    usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        password_salt TEXT NOT NULL,
        public_key TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    
    // Messages table
    messagesTable := `
    CREATE TABLE IF NOT EXISTS messages (
        id TEXT PRIMARY KEY,
        from_user TEXT NOT NULL,
        to_user TEXT,
        channel_id TEXT,
        content TEXT NOT NULL,
        encrypted BOOLEAN DEFAULT 0,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (from_user) REFERENCES users(id),
        FOREIGN KEY (to_user) REFERENCES users(id),
        FOREIGN KEY (channel_id) REFERENCES channels(id)
    );`
    
    // Channels table
    channelsTable := `
    CREATE TABLE IF NOT EXISTS channels (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        created_by TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (created_by) REFERENCES users(id)
    );`
    
    // Channel members table
    channelMembersTable := `
    CREATE TABLE IF NOT EXISTS channel_members (
        channel_id TEXT NOT NULL,
        user_id TEXT NOT NULL,
        joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (channel_id, user_id),
        FOREIGN KEY (channel_id) REFERENCES channels(id),
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
    
    // Sessions table
    sessionsTable := `
    CREATE TABLE IF NOT EXISTS sessions (
        token TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
    
    tables := []string{usersTable, messagesTable, channelsTable, channelMembersTable, sessionsTable}
    
    for _, table := range tables {
        if _, err := d.db.Exec(table); err != nil {
            return fmt.Errorf("failed to create table: %v", err)
        }
    }
    
    return nil
}

func (d *Database) Close() error {
    return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
    return d.db
}
