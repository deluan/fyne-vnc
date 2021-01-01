package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	vnc "github.com/amitbet/vnc2video"
)

func createVncConfig(password string) *vnc.ClientConfig {
	cchServer := make(chan vnc.ServerMessage, 1)
	cchClient := make(chan vnc.ClientMessage, 1)
	errorCh := make(chan error, 1)

	return &vnc.ClientConfig{
		SecurityHandlers: []vnc.SecurityHandler{
			&vnc.ClientAuthVNC{Password: []byte(password)},
			&vnc.ClientAuthNone{},
		},
		DrawCursor:      true,
		PixelFormat:     vnc.PixelFormat16bit,
		ClientMessageCh: cchClient,
		ServerMessageCh: cchServer,
		ErrorCh:         errorCh,
		Messages:        vnc.DefaultServerMessages,
		Encodings: []vnc.Encoding{
			&vnc.TightEncoding{},
			&vnc.HextileEncoding{},
			&vnc.ZRLEEncoding{},
			&vnc.CopyRectEncoding{},
			&vnc.CursorPseudoEncoding{},
			&vnc.CursorPosPseudoEncoding{},
			&vnc.ZLibEncoding{},
			&vnc.RREEncoding{},
			&vnc.RawEncoding{},
		},
	}
}

func main() {
	if len(os.Args) < 4 {
		cmd := filepath.Base(os.Args[0])
		fmt.Printf("Usage: %s <vnc|rdp> address port", cmd)
		os.Exit(1)
	}

	config := createVncConfig(os.Args[3])
	vncDisplay := NewVncDisplay(os.Args[1], os.Args[2], config)

	vncApp := app.New()
	title := fmt.Sprintf("VNC (%s:%s)", os.Args[1], os.Args[2])
	w := vncApp.NewWindow(title)

	w.CenterOnScreen()
	top := widget.NewHBox(
		widget.NewButton("Quit", func() {
			vncApp.Quit()
		}),
	)
	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(top, nil, nil, nil),
		top, vncDisplay)
	w.SetContent(content)

	w.Resize(fyne.NewSize(1024, 768))
	w.ShowAndRun()
}
