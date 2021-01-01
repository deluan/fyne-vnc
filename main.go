package main

import (
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	vnc "github.com/amitbet/vnc2video"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultWidth  = 1024
	defaultHeight = 768
)

func createVncConfig(username, password string) *vnc.ClientConfig {
	cchServer := make(chan vnc.ServerMessage, 1)
	cchClient := make(chan vnc.ClientMessage, 1)
	errorCh := make(chan error, 1)

	return &vnc.ClientConfig{
		SecurityHandlers: []vnc.SecurityHandler{
			//&vnc.ClientAuthATEN{Username: []byte(username), Password: []byte(password)},
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

	config := createVncConfig(os.Args[3], os.Args[4])
	vncDisplay := NewVncDisplay(os.Args[1], os.Args[2], config)

	vncApp := app.New()
	title := fmt.Sprintf("%s (%s:%s)", strings.ToUpper(os.Args[1]), os.Args[2], os.Args[3])
	w := vncApp.NewWindow(title)

	w.CenterOnScreen()
	w.SetContent(widget.NewVBox(
		widget.NewHBox(
			widget.NewButton("Quit", func() {
				vncApp.Quit()
			}),
		),
		vncDisplay,
	))

	w.ShowAndRun()
}
