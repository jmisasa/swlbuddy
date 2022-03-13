package main

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/reiver/go-telnet"
	"io"
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

	err, frequency := sendCommand(conn, "f")

	if err != nil {
		fmt.Printf("Couldn't run command: %v", err)
		panic("couldn't couldn't")
	}

	fmt.Printf("err = %v string = %s\n", err, frequency)

	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}

func sendCommand(conn *telnet.Conn, command string) (error, string) {
	_, err := conn.Write([]byte(command))
	_, err = conn.Write([]byte("\n"))

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
