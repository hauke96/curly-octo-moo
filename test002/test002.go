package main

import (
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

	menubar := createMenu()

	button := gtk.NewButtonWithLabel("Weiter")

	haus1 := haus{x: 2, y: 3, state: 0}

	spielstand := spielstand{money: 10, haus1: haus1}

	label := gtk.NewLabel("")
	labelUpdate(label, &spielstand)

	// start update routine
	// TODO Channels benutzen
	go update(&spielstand, label)

	bild := gtk.NewImage()
	setImage(bild, spielstand.haus1.state)

	// lambda für die weiterschaltung
	// Erhöhe den state nur wenn es kein upgrade gibt oder wenn wenn wir genug geld haben.
	// Und nur bei letzterem kaufe das upgrade.
	weiter := func(ctx *glib.CallbackContext) {
		if spielstand.haus1.state < 6*5 {
			if (spielstand.haus1.state+1)%5 != 0 {
				spielstand.haus1.state += 1
			} else {
				if spielstand.money > 9 {
					spielstand.haus1.state += 1
					setImage(bild, spielstand.haus1.state/5)
					spielstand.money -= 10
				}
			}
		}

		labelUpdate(label, &spielstand)
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
	layout.Put(menubar, 0, 0)

	window.Add(layout)
	window.ShowAll()
	gtk.Main()
}

func setImage(bild *gtk.Image, state int) {
	dir, _ := filepath.Split(os.Args[0])
	imagefile := dir + "/state" + strconv.Itoa(state) + ".png"
	bild.SetFromFile(imagefile)
}

func update(stand *spielstand, label *gtk.Label) {
	for true {
		time.Sleep(time.Second)
		stand.money += (stand.haus1.state / 5)
		labelUpdate(label, stand)
	}
}

func labelUpdate(label *gtk.Label, spielstand *spielstand) {
	label.SetText(strconv.Itoa(spielstand.money) + "\n" + strconv.Itoa(spielstand.haus1.state))
}

func createMenu() *gtk.MenuBar {
	menubar := gtk.NewMenuBar()
	vpaned := gtk.NewVPaned()

	//--------------------------------------------------------
	// GtkMenuItem
	//--------------------------------------------------------
	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	var menuitem *gtk.MenuItem
	menuitem = gtk.NewMenuItemWithMnemonic("E_xit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_View")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	checkmenuitem := gtk.NewCheckMenuItemWithMnemonic("_Disable")
	checkmenuitem.Connect("activate", func() {
		vpaned.SetSensitive(!checkmenuitem.GetActive())
	})
	submenu.Append(checkmenuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	return menubar
}
