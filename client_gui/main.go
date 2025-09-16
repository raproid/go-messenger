package main

import (
    "fyne.io/fyne/v2/app"
    "secure-messenger/client"
)

func main() {
    myApp := app.New()
    // App metadata is set automatically by Fyne
    
    // Check if we have a valid session
    sessionManager := client.NewSessionManager()
    if sessionManager.HasValidSession() {
        // Show main chat window
        chatWindow := NewChatWindow(myApp)
        chatWindow.window.Show()
    } else {
        // Show login window
        loginWindow := NewLoginWindow(myApp)
        loginWindow.window.Show()
    }
    
    myApp.Run()
}
