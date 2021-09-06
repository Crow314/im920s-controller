package connector

import (
	"fmt"
	"log"
	"os"

	"github.com/tarm/serial"
)

type Connector struct {
	portName string
	port     *serial.Port
}

func InitConnector(portName string, transmitter <-chan string, receiver chan<- string) *Connector {
	conn := new(Connector)
	conn.portName = portName

	conf := &serial.Config{Name: portName, Baud: 19200}

	p, err := serial.OpenPort(conf)
	if err != nil {
		log.Fatalf("Failed to open COM port\n%v", err)
	}
	conn.port = p

	go conn.transmit(transmitter)
	go conn.receive(receiver)

	return conn
}

func (connector *Connector) transmit(transmitter <-chan string) {
	println("Info: Start Transmitter")

	for {
		msg := <-transmitter

		_, err := connector.port.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
		println("Debug: Transmit: " + msg)
	}
}

func (connector *Connector) receive(receiver chan<- string) {
	println("Info: Start Receiver")
	dataBuf := make([]byte, 128)
	msgBuf := make([]byte, 0, 128)

	// Arduinoのloop()的なアレ
	for {
		var msg string

		// LFまでを取得
	line:
		for {
			n, err := connector.port.Read(dataBuf)
			if err != nil {
				log.Fatal(err)
			}

			_, _ = fmt.Fprintf(os.Stderr, "Debug: Receive %q\n", string(dataBuf[:n]))

			for _, v := range dataBuf[:n] {
				msgBuf = append(msgBuf, v)

				if v == '\n' {
					msg = string(msgBuf)
					msgBuf = msgBuf[:0] // 要素を全て削除
					break line
				}
			}
		}
		receiver <- msg
		println("Debug: Receive Message: " + msg)
	}
}
