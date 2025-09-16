package main

import (
    "crypto/tls"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "secure-messenger/storage"
)

func main() {
    // Initialize database
    db, err := storage.NewDatabase("./data")
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.Close()
    
    // Initialize server
    srv := NewServer(db)
    
    // Load TLS certificate
    cert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key")
    if err != nil {
        log.Fatal("Failed to load TLS certificate:", err)
    }
    
    config := &tls.Config{Certificates: []tls.Certificate{cert}}
    
    // Create TLS listener
    listener, err := tls.Listen("tcp", ":8080", config)
    if err != nil {
        log.Fatal("Failed to create TLS listener:", err)
    }
    defer listener.Close()
    
    fmt.Println("�� Secure Messenger Server started on :8080")
    fmt.Println("📡 Listening for encrypted connections...")
    
    // Handle graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        fmt.Println("\n🛑 Shutting down server...")
        listener.Close()
        os.Exit(0)
    }()
    
    // Accept connections
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        
        go srv.HandleConnection(conn)
    }
}
