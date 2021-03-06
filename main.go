package main

import (
	"fmt"
	"github.com/jmisasa/swlbuddy/rigctl"
	"github.com/jmisasa/swlbuddy/swldata"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"net"
	"os"
	"strconv"
	"time"
)

const TabIndexCurrentlyTransmitting = 1

func main() {
	var currentFrequency = ""

	gtk.Init(&os.Args)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("SWLBuddy")
	window.SetIconName("gtk-dialog-info")

	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		fmt.Println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	conn, err := rigctl.Connect("localhost", 7356)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to RigControl server: %v", err))
	}

	freqLabel := gtk.NewLabel("Frequency")
	freqLabel.ModifyFontEasy("18")

	vpaned := gtk.NewVPaned()

	byFrequencyStore := gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	treeview := gtk.NewTreeView()
	vpaned.Pack1(freqLabel, false, true)
	vpaned.Pack2(treeview, true, false)

	treeview.SetModel(byFrequencyStore)
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Station", gtk.NewCellRendererText(), "text", 0))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Country", gtk.NewCellRendererText(), "text", 1))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Language", gtk.NewCellRendererText(), "text", 2))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Target", gtk.NewCellRendererText(), "text", 3))

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})

	go func(conn net.Conn, label *gtk.Label, listStore *gtk.ListStore) {
		for {
			select {
			case <-ticker.C:
				err, frequency := rigctl.Command(conn, "f")

				if err != nil {
					panic("error")
				}

				if currentFrequency != frequency {
					label.SetText(frequency)
					currentFrequency = frequency

					byFrequencyStore.Clear()

					for _, line := range swldata.GetByFrequency(frequency) {
						var iter gtk.TreeIter
						byFrequencyStore.Append(&iter)
						byFrequencyStore.Set(&iter,
							0, line.Station,
							1, line.CountryName,
							2, line.Language,
							3, line.TargetAreaCode,
						)
					}
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(conn, freqLabel, byFrequencyStore)

	currentDateTimeLabel := gtk.NewLabel("Current hour")
	currentlyTxingVpaned := gtk.NewVPaned()

	currentlyTxingStore := gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	currentlyTxingTreeView := gtk.NewTreeView()
	currentlyTxingVpaned.Pack1(currentDateTimeLabel, false, true)
	currentlyTxingVpaned.Pack2(currentlyTxingTreeView, true, false)

	currentlyTxingTreeView.SetModel(currentlyTxingStore)
	currentlyTxingTreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Frequency", gtk.NewCellRendererText(), "text", 0))
	currentlyTxingTreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Station", gtk.NewCellRendererText(), "text", 1))
	currentlyTxingTreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Country", gtk.NewCellRendererText(), "text", 2))
	currentlyTxingTreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Language", gtk.NewCellRendererText(), "text", 3))
	currentlyTxingTreeView.Connect("row-activated", func(ctx *glib.CallbackContext) {
		var path *gtk.TreePath
		var column *gtk.TreeViewColumn
		currentlyTxingTreeView.GetCursor(&path, &column)

		treeModel := currentlyTxingStore.ToTreeModel()
		var iter gtk.TreeIter
		treeModel.GetIter(&iter, path)

		var value glib.GValue
		value.Init(glib.G_TYPE_STRING)
		treeModel.GetValue(&iter, 0, &value)

		khz, _ := strconv.ParseFloat(value.GetString(), 64)
		hz := fmt.Sprintf("%d\n", int(khz*1000))

		rigctl.Command(conn, "F", hz)
	}, "")

	notebook := gtk.NewNotebook()

	notebook.Connect("switch-page", func(ctx *glib.CallbackContext) {
		tabIndex := ctx.Args(1)

		if tabIndex == TabIndexCurrentlyTransmitting {
			currentlyTxingStore.Clear()

			for _, line := range swldata.GetCurrentlyTransmitting() {
				var iter gtk.TreeIter
				currentlyTxingStore.Append(&iter)
				currentlyTxingStore.Set(&iter,
					0, line.Frequency,
					1, line.Station,
					2, line.CountryName,
					3, line.Language,
				)
			}
		}
	}, "")

	tabLabel := "By frequency"
	notebook.AppendPage(vpaned, gtk.NewLabel("By frequency"))

	tabLabel = "Currently transmitting"
	notebook.AppendPage(currentlyTxingVpaned, gtk.NewLabel(tabLabel))

	window.Add(notebook)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}
