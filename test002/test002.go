package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type haus struct {
	x, y, state int
}

type spielstand struct {
	money int
	haus1 haus
}

func main() {

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("test002 - clicker game stuff thingy")
	window.Connect("destroy", gtk.MainQuit)
	window.SetSizeRequest(800, 600)

	layout := gtk.NewLayout(nil, nil)

	button := gtk.NewButtonWithLabel("Weiter")

	spielstand := spielstand{money: 10, haus1: haus{x: 2, y: 3, state: 0}}

	label := gtk.NewLabel(strconv.Itoa(spielstand.money) + "\n" + strconv.Itoa(spielstand.haus1.state))

	// start update routine
	go update(&spielstand, label)

	bild := gtk.NewImage()
	setImage(bild, spielstand.haus1.state)

	// lambda für die weiterschaltung
	weiter := func(ctx *glib.CallbackContext) {
		if spielstand.haus1.state < 6*5 {
			spielstand.haus1.state += 1
			if (spielstand.haus1.state)%10 == 0 && spielstand.money >= 10 {
				fmt.Println(spielstand.haus1.state, " - ", (spielstand.haus1.state)%5)
				setImage(bild, spielstand.haus1.state/5)
				spielstand.money -= 10
			} else if (spielstand.haus1.state)%10 == 0 && spielstand.money < 10 {
				spielstand.haus1.state -= 1
			}
		}
		label.SetText(strconv.Itoa(spielstand.money) + "\n" + strconv.Itoa(spielstand.haus1.state))
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
	layout.Put(label, 50, 60)
	layout.Put(eventBox, 200, 100)
	window.Add(layout)
	window.ShowAll()
	gtk.Main()
}

func setImage(label *gtk.Image, state int) {
	dir, _ := filepath.Split(os.Args[0])
	fmt.Println(strconv.Itoa(state))
	imagefile := filepath.Join(dir, "state"+strconv.Itoa(state)+".png")
	fmt.Println(imagefile)
	label.SetFromFile(imagefile)
}

func update(stand *spielstand, label *gtk.Label) {
	for true {
		time.Sleep(3 * time.Second)
		fmt.Print(stand.haus1.state)
		stand.money += (stand.haus1.state / 5)
		fmt.Print(" - ", stand.haus1.state, " - ", (stand.haus1.state / 5), "\n")
		label.SetText(strconv.Itoa(stand.money) + "\n" + strconv.Itoa(stand.haus1.state))
	}
}
