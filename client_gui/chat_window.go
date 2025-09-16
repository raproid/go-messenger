package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "secure-messenger/client"
    "secure-messenger/shared"
)

type ChatWindow struct {
    app           fyne.App
    window        fyne.Window
    client        *client.NetworkClient
    messageList   *widget.List
    messageEntry  *widget.Entry
    sendBtn       *widget.Button
    messages      []*shared.Message
    currentChat   string
    chatType      string // "user" or "channel"
}

func NewChatWindow(app fyne.App) *ChatWindow {
    w := app.NewWindow("Secure Messenger")
    w.Resize(fyne.NewSize(800, 600))
    w.CenterOnScreen()
    
    cw := &ChatWindow{
        app:    app,
        window: w,
        client: client.NewNetworkClient(),
    }
    
    cw.setupUI()
    cw.loadSession()
    return cw
}

func (cw *ChatWindow) setupUI() {
    // Message list
    cw.messageList = widget.NewList(
        func() int {
            return len(cw.messages)
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(id widget.ListItemID, obj fyne.CanvasObject) {
            if id < len(cw.messages) {
                msg := cw.messages[id]
                label := obj.(*widget.Label)
                label.SetText(msg.Content)
            }
        },
    )
    
    // Message entry
    cw.messageEntry = widget.NewMultiLineEntry()
    cw.messageEntry.SetPlaceHolder("Type your message...")
    
    // Send button
    cw.sendBtn = widget.NewButton("Send", func() {
        cw.sendMessage()
    })
    
    // Chat list (users and channels)
    chatList := widget.NewList(
        func() int {
            return 10 // Placeholder
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(id widget.ListItemID, obj fyne.CanvasObject) {
            label := obj.(*widget.Label)
            label.SetText("Chat " + string(rune(id)))
        },
    )
    
    // Layout
    chatPanel := container.NewBorder(
        nil,
        container.NewHBox(cw.messageEntry, cw.sendBtn),
        nil,
        nil,
        cw.messageList,
    )
    
    mainContent := container.NewHSplit(
        chatList,
        chatPanel,
    )
    mainContent.SetOffset(0.3)
    
    cw.window.SetContent(mainContent)
}

func (cw *ChatWindow) loadSession() {
    sessionManager := client.NewSessionManager()
    session, err := sessionManager.LoadSession()
    if err != nil {
        dialog.ShowError(fmt.Errorf("Failed to load session: %v", err), cw.window)
        return
    }
    
    cw.client.Session = session
    
    // Connect to server
    if err := cw.client.Connect(); err != nil {
        dialog.ShowError(fmt.Errorf("Failed to connect to server: %v", err), cw.window)
        return
    }
    
    // Load recent messages
    cw.loadRecentMessages()
}

func (cw *ChatWindow) sendMessage() {
    content := cw.messageEntry.Text
    if content == "" {
        return
    }
    
    if cw.currentChat == "" {
        dialog.ShowError(fmt.Errorf("Please select a chat"), cw.window)
        return
    }
    
    var err error
    if cw.chatType == "user" {
        err = cw.client.SendMessage(cw.currentChat, content)
    } else {
        err = cw.client.SendChannelMessage(cw.currentChat, content)
    }
    
    if err != nil {
        dialog.ShowError(fmt.Errorf("Failed to send message: %v", err), cw.window)
        return
    }
    
    cw.messageEntry.SetText("")
    cw.loadRecentMessages()
}

func (cw *ChatWindow) loadRecentMessages() {
    if cw.currentChat == "" {
        return
    }
    
    var err error
    if cw.chatType == "user" {
        cw.messages, err = cw.client.GetMessages(cw.currentChat, 50)
    } else {
        cw.messages, err = cw.client.GetChannelMessages(cw.currentChat, 50)
    }
    
    if err != nil {
        dialog.ShowError(fmt.Errorf("Failed to load messages: %v", err), cw.window)
        return
    }
    
    cw.messageList.Refresh()
}
