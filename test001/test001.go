package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"unsafe"
)

// Ein einfacher point ;)
type point struct {
	x int
	y int
}

func main() {
	//--------------------------------------------------------
	//
	// GTK initialization + Window creation
	//
	//--------------------------------------------------------
	// Initialisiert gtk für den gebraucht und erzeugt das Fenster.
	//--------------------------------------------------------
	// gtk muss einmal initialisiert werden
	gtk.Init(nil)

	// Neues Fenster erstellen, rest ist trivial
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetSizeRequest(600, 600)
	window.SetTitle("GTK Go!")
	// man kann icons aus /usr/share/icons/<Dein-aktuelles-theme>/... angeben
	window.SetIconName("user-images")
	// das schlißene-event
	window.Connect("destroy", gtk.MainQuit)

	//--------------------------------------------------------
	//
	// Gtk VBox
	//
	//--------------------------------------------------------
	// Ein vertikales design. Hier ists relativ egal welches man nimmt.
	//--------------------------------------------------------

	vbox := gtk.NewVBox(true, 1)
	// funktioniert genauso gut:
	//	vbox := gtk.NewHBox(true, 1)

	// kleinen Rand drum rum
	// wenn man außerhalb des weißen Bereiches (also der drawingarea) klickt sieht man,
	// dass auch das Click-Event nicht ausgeführt wird. Dazu unten mehr
	vbox.SetBorderWidth(13)

	//--------------------------------------------------------
	//
	// Drawing Area + Events
	//
	//--------------------------------------------------------
	// Eine art Panel auf dem man zeichnen kann.
	// Erzeugt alle Event mittels lambdas (aka anonyme Funktionen aka anonymous functions)
	// Alle Events werden aus der Hauptschleife aufgerufen (s.u.)
	//--------------------------------------------------------
	drawingarea := gtk.NewDrawingArea()
	// trivial :D
	createEvents(drawingarea)

	//--------------------------------------------------------
	//
	// Final stuff
	//
	//--------------------------------------------------------
	// Fügt alles dem Fenster hinzu und zeigt es an.
	//--------------------------------------------------------

	// die drawingarea dem vbox hinzufügen
	vbox.Add(drawingarea)

	window.Add(vbox)
	window.ShowAll()
	// geht in die gtk Hauptschleife für alle Events etc.
	gtk.Main()
}

/*
 * Erzeugt die ganzen Events für das drawing area mittels lambdas.
 *
 * Parameter:
 *     drawingarea - Die drawing area auf der gezeichnet werden soll
 *
 * Zu den Schnittstellenkommentaren: Leider gibts das anscheinend nicht direkt
 *                                   in go (jedenfall kein @param etc.) :/
 */
func createEvents(drawingarea *gtk.DrawingArea) {
	// Wir brauchen ne pixmap in der wir die Pixeldaten speichern
	var pixmap *gdk.Pixmap
	// gdk.GC ist einfach eine Sammlung an Eigenschaften und Einstellungen zum zeichnen
	var gc *gdk.GC
	p := point{x: -1, y: -1}

	drawingarea.Connect("configure-event", func() {
		// wenns schon ne pixmap gibt, lösche diese
		if pixmap != nil {
			pixmap.Unref()
		}

		// hole pixmap und stuff
		allocation := drawingarea.GetAllocation()
		pixmap = gdk.NewPixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)

		// weißen Hintergrund zeichnen:
		gc = gdk.NewGC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.NewColor("white"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)

		// Vorder- und Hintergrundfarbe setzen
		gc.SetRgbFgColor(gdk.NewColor("black"))
		gc.SetRgbBgColor(gdk.NewColor("red"))
	})

	// dieses event wird ausgeführt wenn auf das widget geklickt wurde
	drawingarea.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		// Argumente holen
		arg := ctx.Args(0)
		// irgend son pointer auf ein objekt holen in dem mehr Infos stehen
		mev := *(**gdk.EventMotion)(unsafe.Pointer(&arg))
		// Position ausgeben
		fmt.Println("Geklickte Position: ", p)

		// hier ist meine eigene Logik:
		// Wenn noch kein Klick registriert wurde ist p.x == p.y == -1 und da soll noch nichts gezeichnet werden
		// erst wenn der zweite klick kommt soll zwischen diesem und dem alten eine Linie gezeichnet werden
		p_neu := point{x: int(mev.X), y: int(mev.Y)}
		if p.x != -1 && p.y != -1 {
			pixmap.GetDrawable().DrawLine(gc, p.x, p.y, p_neu.x, p_neu.y)
			drawingarea.GetWindow().Invalidate(nil, false)
		}
		// Position des neuen klicks speichern, es ist also der neue startpunkt der nächsten linie
		p = p_neu
	})

	// dieses event wird ausgeführt wenn drawingarea.GetWindow().Invalidate(nil, false) aufgerufen wird
	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
	})

	// Events müssen manuell spezifiziert werden.
	// Hier werden immer MASKs übergeben, hier also eben die BUTTON_PRESS_MASK,
	// welche erst in einen int gecastet werden muss
	drawingarea.SetEvents(int(gdk.BUTTON_PRESS_MASK))
}
