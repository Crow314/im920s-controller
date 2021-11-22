package main

import (
	"bufio"
	"fmt"
	"github.com/Crow314/im920s-controller/pkg/connector"
	"github.com/Crow314/im920s-controller/pkg/module"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Please input COM port path/name: ")
	var portName string
	if scanner.Scan() {
		portName = scanner.Text()
	}

	conn := connector.NewConnector(portName)
	im920s := module.NewIm920s(conn.TransmitChannel(), conn.ReceiveChannel())

	go func() {
		for {
			data := <-im920s.DataReceiver()
			fmt.Printf("node: %d\n", data.Node())
			fmt.Printf("RSSI: %ddb\n", data.Rssi())

			for _, v := range data.Data() {
				fmt.Printf("%X, ", v)
			}
			fmt.Println()
		}
	}()

	for {
		var msg string
		if scanner.Scan() {
			msg = scanner.Text()
		}

		res, err := im920s.SendCommand(msg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			fmt.Println(res)
		}
	}
}
