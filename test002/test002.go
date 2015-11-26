package main

/**

 */

import (
	//"github.com/mattn/go-gtk/gdk"
	//"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("test002 - clicker game stuff thingy")
	window.Connect("destroy", gtk.MainQuit)
	window.SetSizeRequest(800, 600)

	// ...

	window.ShowAll()
	gtk.Main()
}
