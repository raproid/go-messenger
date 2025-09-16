package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "secure-messenger/client"
    "strconv"
)

type ServerConfigWindow struct {
    app    fyne.App
    window fyne.Window
}

func NewServerConfigWindow(app fyne.App) *ServerConfigWindow {
    w := app.NewWindow("Server Configuration")
    w.Resize(fyne.NewSize(400, 300))
    w.CenterOnScreen()
    
    scw := &ServerConfigWindow{
        app:    app,
        window: w,
    }
    
    scw.setupUI()
    return scw
}

func (scw *ServerConfigWindow) setupUI() {
    // Load current config
    config, err := client.LoadConfig()
    if err != nil {
        dialog.ShowError(fmt.Errorf("Failed to load config: %v", err), scw.window)
        return
    }
    
    // Server address field
    addressEntry := widget.NewEntry()
    addressEntry.SetText(config.ServerAddress)
    addressEntry.SetPlaceHolder("Server Address")
    
    // Server port field
    portEntry := widget.NewEntry()
    portEntry.SetText(strconv.Itoa(config.ServerPort))
    portEntry.SetPlaceHolder("Server Port")
    
    // Server name field
    serverNameEntry := widget.NewEntry()
    serverNameEntry.SetText(config.ServerName)
    serverNameEntry.SetPlaceHolder("Server Name")
    
    // Use TLS checkbox
    useTLSCheck := widget.NewCheck("Use TLS", nil)
    useTLSCheck.SetChecked(config.UseTLS)
    
    // Certificate path field
    certPathEntry := widget.NewEntry()
    certPathEntry.SetText(config.CertPath)
    certPathEntry.SetPlaceHolder("Certificate Path")
    
    // Save button
    saveBtn := widget.NewButton("Save", func() {
        scw.saveConfig(addressEntry.Text, portEntry.Text, serverNameEntry.Text, useTLSCheck.Checked, certPathEntry.Text)
    })
    
    // Cancel button
    cancelBtn := widget.NewButton("Cancel", func() {
        scw.window.Close()
    })
    
    // Layout
    form := container.NewVBox(
        widget.NewLabel("Server Configuration"),
        widget.NewSeparator(),
        addressEntry,
        portEntry,
        serverNameEntry,
        useTLSCheck,
        certPathEntry,
        widget.NewSeparator(),
        container.NewHBox(saveBtn, cancelBtn),
    )
    
    scw.window.SetContent(form)
}

func (scw *ServerConfigWindow) saveConfig(address, port, serverName string, useTLS bool, certPath string) {
    // Parse port
    portInt, err := strconv.Atoi(port)
    if err != nil {
        dialog.ShowError(fmt.Errorf("Invalid port number"), scw.window)
        return
    }
    
    // Create config
    config := &client.Config{
        ServerAddress: address,
        ServerPort:    portInt,
        ServerName:    serverName,
        UseTLS:        useTLS,
        CertPath:      certPath,
    }
    
    // Save config
    if err := client.SaveConfig(config); err != nil {
        dialog.ShowError(fmt.Errorf("Failed to save config: %v", err), scw.window)
        return
    }
    
    dialog.ShowInformation("Success", "Configuration saved successfully", scw.window)
    scw.window.Close()
}
