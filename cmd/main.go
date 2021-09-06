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
	connector.InitConnector("COM5", im920s.TransmitChannel, im920s.ReceiveChannel)

	go func() {
		for {
			msg := <-im920s.ReceiveChannel
			msg = strings.Replace(msg, "\r\n", "", -1)
			fmt.Println(msg)
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
