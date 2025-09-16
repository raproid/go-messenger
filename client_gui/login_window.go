package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "secure-messenger/client"
    "time"
)

type LoginWindow struct {
    app    fyne.App
    window fyne.Window
    client *client.NetworkClient
}

func NewLoginWindow(app fyne.App) *LoginWindow {
    w := app.NewWindow("Secure Messenger - Login")
    w.Resize(fyne.NewSize(400, 300))
    w.CenterOnScreen()
    
    lw := &LoginWindow{
        app:    app,
        window: w,
        client: client.NewNetworkClient(),
    }
    
    lw.setupUI()
    return lw
}

func (lw *LoginWindow) setupUI() {
    // Username field
    usernameEntry := widget.NewEntry()
    usernameEntry.SetPlaceHolder("Username")
    
    // Password field
    passwordEntry := widget.NewPasswordEntry()
    passwordEntry.SetPlaceHolder("Password")
    
    // Email field (for registration)
    emailEntry := widget.NewEntry()
    emailEntry.SetPlaceHolder("Email")
    
    // Login button
    loginBtn := widget.NewButton("Login", func() {
        lw.handleLogin(usernameEntry.Text, passwordEntry.Text)
    })
    
    // Register button
    registerBtn := widget.NewButton("Register", func() {
        lw.handleRegister(usernameEntry.Text, emailEntry.Text, passwordEntry.Text)
    })
    
    // Server config button
    configBtn := widget.NewButton("Server Config", func() {
        lw.showServerConfig()
    })
    
    // Certificate button
    certBtn := widget.NewButton("Certificate", func() {
        lw.showCertificateWindow()
    })
    
    // Layout
    form := container.NewVBox(
        widget.NewLabel("Secure Messenger"),
        widget.NewSeparator(),
        usernameEntry,
        passwordEntry,
        emailEntry,
        container.NewHBox(loginBtn, registerBtn),
        widget.NewSeparator(),
        container.NewHBox(configBtn, certBtn),
    )
    
    lw.window.SetContent(form)
}

func (lw *LoginWindow) handleLogin(username, password string) {
    if username == "" || password == "" {
        dialog.ShowError(fmt.Errorf("Please enter username and password"), lw.window)
        return
    }
    
    // Connect to server
    if err := lw.client.Connect(); err != nil {
        dialog.ShowError(fmt.Errorf("Failed to connect to server: %v", err), lw.window)
        return
    }
    defer lw.client.Disconnect()
    
    // Attempt login
    response, err := lw.client.Login(username, password)
    if err != nil {
        dialog.ShowError(fmt.Errorf("Login failed: %v", err), lw.window)
        return
    }
    
    if !response.Success {
        dialog.ShowError(fmt.Errorf("Login failed: %s", response.Error), lw.window)
        return
    }
    
    // Save session
    sessionManager := client.NewSessionManager()
    session := &client.Session{
        Token:    response.Token,
        User:     response.User,
        LastSeen: time.Now(),
    }
    sessionManager.SaveSession(session)
    
    // Show chat window
    chatWindow := NewChatWindow(lw.app)
    chatWindow.window.Show()
    lw.window.Hide()
}

func (lw *LoginWindow) handleRegister(username, email, password string) {
    if username == "" || email == "" || password == "" {
        dialog.ShowError(fmt.Errorf("Please enter all fields"), lw.window)
        return
    }
    
    // Connect to server
    if err := lw.client.Connect(); err != nil {
        dialog.ShowError(fmt.Errorf("Failed to connect to server: %v", err), lw.window)
        return
    }
    defer lw.client.Disconnect()
    
    // Attempt registration
    response, err := lw.client.Register(username, email, password)
    if err != nil {
        dialog.ShowError(fmt.Errorf("Registration failed: %v", err), lw.window)
        return
    }
    
    if !response.Success {
        dialog.ShowError(fmt.Errorf("Registration failed: %s", response.Error), lw.window)
        return
    }
    
    // Save session
    sessionManager := client.NewSessionManager()
    session := &client.Session{
        Token:    response.Token,
        User:     response.User,
        LastSeen: time.Now(),
    }
    sessionManager.SaveSession(session)
    
    // Show chat window
    chatWindow := NewChatWindow(lw.app)
    chatWindow.window.Show()
    lw.window.Hide()
}

func (lw *LoginWindow) showServerConfig() {
    configWindow := NewServerConfigWindow(lw.app)
    configWindow.window.Show()
}

func (lw *LoginWindow) showCertificateWindow() {
    certWindow := NewCertificateWindow(lw.app)
    certWindow.window.Show()
}
