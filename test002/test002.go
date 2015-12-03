package main

import (
	//	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type haus struct {
	state int
}

type spielstand struct {
	money, amountUpgrades int
	spielfeld             [5][5]*haus
}

func main() {

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("test002 - clicker game stuff thingy")
	window.Connect("destroy", gtk.MainQuit)
	window.SetSizeRequest(800, 600)

	layout := gtk.NewLayout(nil, nil)
	layout.ModifyBG(gtk.STATE_NORMAL, gdk.NewColorRGB(200, 200, 200))

	menubar := createMenu()

	spielstand := spielstand{money: 10, spielfeld: [5][5]*haus{}, amountUpgrades: 0}

	label := gtk.NewLabel("")
	labelUpdate(label, &spielstand)

	for i := 0; i < len(spielstand.spielfeld); i++ {
		for j := 0; j < len(spielstand.spielfeld[0]); j++ {
			spielHaus := haus{state: -1}
			spielstand.spielfeld[i][j] = &spielHaus

			bild := gtk.NewImage()
			if i == 2 && j == 3 {
				spielHaus.state = 0
			}
			setImage(bild, spielHaus.state)

			// lambda für die weiterschaltung
			// Erhöhe den state nur wenn es kein upgrade gibt oder wenn wenn wir genug geld haben.
			// Und nur bei letzterem kaufe das upgrade.
			weiter := func(ctx *glib.CallbackContext) {
				if spielHaus.state == -1 && spielstand.amountUpgrades != 0 && spielstand.money >= spielstand.amountUpgrades*10 {
					spielHaus.state = 0
					setImage(bild, 0)
					spielstand.money -= spielstand.amountUpgrades * 10
				} else {
					if spielHaus.state < 6*5 && spielHaus.state != -1 {
						if (spielHaus.state+1)%5 != 0 {
							spielHaus.state += 1
						} else {
							if spielstand.money > 9 {
								spielHaus.state += 1
								setImage(bild, spielHaus.state/5)
								spielstand.money -= 10
								spielstand.amountUpgrades++
							}
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
			layout.Put(eventBox, 100+i*100, 100+j*100)
		}
	}

	// start update routine
	go update(&spielstand, label)

	layout.Put(label, 300, 70)
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
		time.Sleep(3 * time.Second)
		stand.money += stand.amountUpgrades
		labelUpdate(label, stand)
	}
}

func labelUpdate(label *gtk.Label, spielstand *spielstand) {
	label.SetText("Geld: " + strconv.Itoa(spielstand.money) + "\nUpgrades: " + strconv.Itoa(spielstand.amountUpgrades))
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
