package main

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/reiver/go-telnet"
	"io"
	"os"
	"time"
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

	freqLabel := gtk.NewLabel("Frequency")

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})

	go func(conn *telnet.Conn, label *gtk.Label) {
		for {
			select {
			case <-ticker.C:
				err, frequency := sendCommand(conn, "f")

				if err != nil {
					panic("error")
				}

				fmt.Println("Freq %s", frequency)
				label.SetText(frequency)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(conn, freqLabel)

	vpaned := gtk.NewVPaned()

	store := gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_BOOL, gdkpixbuf.GetType())
	treeview := gtk.NewTreeView()
	vpaned.Pack1(freqLabel, false, true)
	vpaned.Pack2(treeview, true, false)

	treeview.SetModel(store)
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Station", gtk.NewCellRendererText(), "text", 0))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Country", gtk.NewCellRendererText(), "text", 1))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Language", gtk.NewCellRendererText(), "text", 2))

	/*
		n := 0
		gtk.StockListIDs().ForEach(func(d unsafe.Pointer, v interface{}) {
			id := glib.GPtrToString(d)
			var iter gtk.TreeIter
			store.Append(&iter)
			store.Set(&iter,
				0, id,
				1, (n == 1),
				2, gtk.NewImage().RenderIcon(id, gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf,
			)
			n = 1 - n
		})
	*/
	notebook := gtk.NewNotebook()

	tabLabel := "By frequency"
	notebook.AppendPage(vpaned, gtk.NewLabel("By frequency"))

	tabLabel = "Currently transmitting"
	currentlyTransmitting := gtk.NewFrame(tabLabel)
	notebook.AppendPage(currentlyTransmitting, gtk.NewLabel(tabLabel))

	window.Add(notebook)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}

func sendCommand(conn *telnet.Conn, command string) (error, string) {
	_, err := conn.Write([]byte(command))
	_, err = conn.Write([]byte("\n"))

	fmt.Println("command sent\n")

	if err != nil {
		return err, ""
	}

	return readLine(conn)
}

func readLine(reader io.Reader) (error, string) {
	b := make([]byte, 1)

	line := ""
	var err error

	for {
		_, err := reader.Read(b)

		fmt.Printf("read byte %v\n", b)

		if bytes.Equal(b, []byte{10}) {
			break
		}

		if err != nil {
			break
		}

		line += string(b)
	}

	return err, line
}
