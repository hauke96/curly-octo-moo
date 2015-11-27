package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"path/filepath"
	"strconv"
)

type haus struct {
	x, y, state int
}

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("test002 - clicker game stuff thingy")
	window.Connect("destroy", gtk.MainQuit)
	window.SetSizeRequest(800, 600)

	layout := gtk.NewLayout(nil, nil)

	button := gtk.NewButtonWithLabel("Weiter")

	haus1 := haus{x: 2, y: 3, state: 0}
	bild := gtk.NewImage()
	setLabel(bild, haus1.state)

	// lambda für die weiterschaltung
	weiter := func(ctx *glib.CallbackContext) {
		fmt.Println("OK")
		if haus1.state < 6*5 {
			haus1.state += 1
		}
		if haus1.state%5 == 0 {
			setLabel(bild, haus1.state/5)
		}
	}

	// Ein bild kann ohne weiteres keine events, daher braucht man eine
	// EventBox drum rum um events zu haben
	eventBox := gtk.NewEventBox()
	eventBox.Add(bild)
	// das event kann nicht direkt auf das image gesetzt werden, deswegen
	// diese event box dafür.
	eventBox.Connect("button-press-event", weiter)
	eventBox.SetEvents(int(gdk.BUTTON_PRESS_MASK))
	button.Connect("button-press-event", weiter)
	button.SetEvents(int(gdk.BUTTON_PRESS_MASK))

	layout.Put(button, 50, 100)
	layout.Put(eventBox, 200, 100)
	window.Add(layout)
	window.ShowAll()
	gtk.Main()
}

func setLabel(label *gtk.Image, state int) {
	dir, _ := filepath.Split(os.Args[0])
	fmt.Println(strconv.Itoa(state))
	imagefile := filepath.Join(dir, "state"+strconv.Itoa(state)+".png")
	fmt.Println(imagefile)
	label.SetFromFile(imagefile)
}
