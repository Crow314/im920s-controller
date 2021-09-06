package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"IM920s-controller-pi/pkg/connector"
	"IM920s-controller-pi/pkg/module"
)

func main() {
	im920s := module.NewIm920s()
	connector.InitConnector("COM5", im920s.UartTransmitter(), im920s.UartReceiver())

	im920s.Init()

	go func() {
		for {
			data := <-im920s.DataReceiver()
			fmt.Printf("node: %d\n", data.Node())
			fmt.Printf("RSSI: %ddb\n", data.Rssi())
			fmt.Println(data.Data())
		}
	}()

	for {
		var msg string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			msg = scanner.Text()
		}

		if !strings.HasSuffix(msg, "\r\n") { // 末尾がCRLFでない場合
			msg += "\r\n"
		}

		im920s.SendCommand(msg)
	}
}
