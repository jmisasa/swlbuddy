package rigctl

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func Connect(host string, port int) (net.Conn, error) {
	return net.Dial("tcp", host+":"+strconv.Itoa(port))
}

func Command(conn net.Conn, command string, params ...string) (error, string) {
	finalCommand := command

	if len(params) > 0 {
		finalCommand = command + " " + strings.Join(params, " ")
	}

	_, err := conn.Write([]byte(finalCommand))
	_, err = conn.Write([]byte("\n"))

	if err != nil {
		return err, ""
	}

	time.Sleep(100 * time.Millisecond)
	rawResult, err := bufio.NewReader(conn).ReadBytes(10)

	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	commandResult := strings.TrimSuffix(string(rawResult), "\n")

	return nil, commandResult
}
