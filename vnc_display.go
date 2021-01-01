package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"net"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	vnc "github.com/amitbet/vnc2video"
	"github.com/boz/go-throttle"
	"github.com/sirupsen/logrus"
)

// Custom Widget that represents the remote computer being controlled
type VncDisplay struct {
	widget.BaseWidget
	//keyboardHandler
	//mouseHandler

	lastUpdate  time.Time
	screenImage *vnc.VncCanvas
	Display     image.Image
	Client      vnc.Conn
}

// Creates a new VncDisplay and does all the heavy lifting, setting up all event handlers
func NewVncDisplay(hostname, port string, config *vnc.ClientConfig) *VncDisplay {
	client := createVncClient(hostname, port, config)
	width := int(client.Width())
	height := int(client.Height())
	empty := image.NewNRGBA(image.Rect(0, 0, width-1, height-1))

	b := &VncDisplay{
		Client: client,
	}
	//b.keyboardHandler.display = b
	//b.mouseHandler.display = b

	b.SetDisplay(empty)
	b.screenImage = client.Canvas

	go func() {
		framerate := 12
		period := time.Duration(1000/framerate) * time.Millisecond
		update := throttle.ThrottleFunc(period, true, func() {
			b.updateDisplay()
		})
		defer update.Stop()
		for {
			select {
			case err := <-config.ErrorCh:
				panic(err)
			case msg := <-config.ClientMessageCh:
				logrus.Tracef("Received client message type:%v msg:%v", msg.Type(), msg)
			case msg := <-config.ServerMessageCh:
				logrus.Tracef("Received server message type:%v msg:%v", msg.Type(), msg)
				if msg.Type() == vnc.FramebufferUpdateMsgType {
					logrus.Debug("Received FramebufferUpdateMsg:  %s", time.Now())
					reqMsg := vnc.FramebufferUpdateRequest{Inc: 1, X: 0, Y: 0, Width: client.Width(), Height: client.Height()}
					//client.ResetAllEncodings()
					reqMsg.Write(client)
					update.Trigger()
				}
			}
		}
	}()

	return b
}

func createVncClient(hostname, port string, ccfg *vnc.ClientConfig) *vnc.ClientConn {
	// Establish TCP connection to VNC server.
	nc, err := net.DialTimeout("tcp", hostname+":"+port, 5*time.Second)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to VNC host. %v", err))
	}

	cc, err := vnc.Connect(context.Background(), nc, ccfg)
	if err != nil {
		panic(fmt.Sprintf("Error negotiating connection to VNC host. %v", err))
	}
	screenImage := cc.Canvas
	for _, enc := range ccfg.Encodings {
		myRenderer, ok := enc.(vnc.Renderer)
		if ok {
			myRenderer.SetTargetImage(screenImage)
		}
	}

	cc.SetEncodings([]vnc.EncodingType{
		vnc.EncCursorPseudo,
		vnc.EncPointerPosPseudo,
		//vnc.EncTight,
		vnc.EncZRLE,
		vnc.EncCopyRect,
		vnc.EncHextile,
		vnc.EncZlib,
		vnc.EncRRE,
	})

	return cc
}

func (b *VncDisplay) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return fyne.Size{
		Width:  b.Display.Bounds().Dx(),
		Height: b.Display.Bounds().Dy(),
	}
}

func (b *VncDisplay) CreateRenderer() fyne.WidgetRenderer {
	return &vncDisplayRenderer{
		objects: []fyne.CanvasObject{},
		remote:  b,
	}
}

// This forces a display refresh if there are any pending updates, instead of waiting for the next
// OnSync event. Ex: when moving the mouse we want instant feedback of the new cursor position
func (b *VncDisplay) updateDisplay() {
	img := b.screenImage.Image
	b.SetDisplay(img)
	b.lastUpdate = time.Now()
	//f, err := os.Create(fmt.Sprintf("%d.png", b.lastUpdate.UnixNano()))
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//jpeg.Encode(f, img, nil)
}

func (b *VncDisplay) SetDisplay(img image.Image) {
	b.Display = img
	b.Refresh()
}

type vncDisplayRenderer struct {
	objects []fyne.CanvasObject
	remote  *VncDisplay
}

func (r *vncDisplayRenderer) MinSize() fyne.Size {
	return r.remote.MinSize()
}

func (r *vncDisplayRenderer) Layout(size fyne.Size) {
	if len(r.objects) == 0 {
		return
	}

	r.objects[0].Resize(size)
}

func (r *vncDisplayRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *vncDisplayRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *vncDisplayRenderer) Refresh() {
	if len(r.objects) == 0 {
		raster := canvas.NewImageFromImage(r.remote.Display)
		raster.FillMode = canvas.ImageFillContain
		r.objects = append(r.objects, raster)
	} else {
		r.objects[0].(*canvas.Image).Image = r.remote.Display
	}
	r.Layout(r.remote.Size())
	canvas.Refresh(r.remote)
}

func (r *vncDisplayRenderer) Destroy() {
}
