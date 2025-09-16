package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "io/ioutil"
)

type CertificateWindow struct {
    app    fyne.App
    window fyne.Window
}

func NewCertificateWindow(app fyne.App) *CertificateWindow {
    w := app.NewWindow("Certificate Management")
    w.Resize(fyne.NewSize(500, 400))
    w.CenterOnScreen()
    
    cw := &CertificateWindow{
        app:    app,
        window: w,
    }
    
    cw.setupUI()
    return cw
}

func (cw *CertificateWindow) setupUI() {
    // Certificate info label
    infoLabel := widget.NewLabel("Certificate Information")
    infoLabel.TextStyle = fyne.TextStyle{Bold: true}
    
    // Certificate path display
    certPathLabel := widget.NewLabel("No certificate selected")
    
    // Select certificate button
    selectBtn := widget.NewButton("Select Certificate", func() {
        cw.selectCertificate(certPathLabel)
    })
    
    // Download certificate button
    downloadBtn := widget.NewButton("Download Certificate", func() {
        cw.downloadCertificate()
    })
    
    // Certificate content display
    certContent := widget.NewMultiLineEntry()
    certContent.SetPlaceHolder("Certificate content will appear here...")
    certContent.Disable()
    
    // Load certificate button
    loadBtn := widget.NewButton("Load Certificate", func() {
        cw.loadCertificate(certPathLabel, certContent)
    })
    
    // Close button
    closeBtn := widget.NewButton("Close", func() {
        cw.window.Close()
    })
    
    // Layout
    form := container.NewVBox(
        infoLabel,
        widget.NewSeparator(),
        certPathLabel,
        container.NewHBox(selectBtn, loadBtn),
        widget.NewSeparator(),
        certContent,
        widget.NewSeparator(),
        container.NewHBox(downloadBtn, closeBtn),
    )
    
    cw.window.SetContent(form)
}

func (cw *CertificateWindow) selectCertificate(certPathLabel *widget.Label) {
    dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
        if err != nil {
            dialog.ShowError(fmt.Errorf("Failed to open file: %v", err), cw.window)
            return
        }
        if reader == nil {
            return
        }
        defer reader.Close()
        
        certPathLabel.SetText(reader.URI().Path())
    }, cw.window)
}

func (cw *CertificateWindow) downloadCertificate() {
    dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
        if err != nil {
            dialog.ShowError(fmt.Errorf("Failed to save file: %v", err), cw.window)
            return
        }
        if writer == nil {
            return
        }
        defer writer.Close()
        
        // This would typically download from server
        dialog.ShowInformation("Info", "Certificate download not implemented yet", cw.window)
    }, cw.window)
}

func (cw *CertificateWindow) loadCertificate(certPathLabel *widget.Label, certContent *widget.Entry) {
    path := certPathLabel.Text
    if path == "No certificate selected" {
        dialog.ShowError(fmt.Errorf("Please select a certificate file first"), cw.window)
        return
    }
    
    content, err := ioutil.ReadFile(path)
    if err != nil {
        dialog.ShowError(fmt.Errorf("Failed to read certificate: %v", err), cw.window)
        return
    }
    
    certContent.SetText(string(content))
}
