package main

import (
	"fmt"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/reiver/go-telnet"
	"os"
)

func main() {
	gtk.Init(&os.Args)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("SWLBuddy")
	window.SetIconName("gtk-dialog-info")

	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		fmt.Println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	conn, err := telnet.DialTo("localhost:7356")

	if err != nil {
		panic("Couldn't connect to RigControl server")
	}

	_, err = conn.Write([]byte("f"))

	if err != nil {
		panic("Couldn't write to server")
	}

	//	b := make([]byte, 4096)
	//	n, err := conn.Read(b)

	if err != nil {
		panic("Couldn't read from server")
	}

	//	fmt.Printf("n = %v err = %v b = %v\n", n, err, b)

	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}
