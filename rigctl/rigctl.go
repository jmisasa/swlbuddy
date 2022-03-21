package rigctl

import (
	"bytes"
	"github.com/reiver/go-telnet"
	"io"
	"strconv"
	"strings"
)

func Connect(host string, port int) (*telnet.Conn, error) {
	return telnet.DialTo(host + ":" + strconv.Itoa(port))
}

func Command(conn *telnet.Conn, command string, params ...string) (error, string) {
	finalCommand := command

	if len(params) > 0 {
		finalCommand = command + " " + strings.Join(params, " ")
	}

	_, err := conn.Write([]byte(finalCommand))
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
