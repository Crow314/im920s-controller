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
	txChan   chan string
	rxChan   chan string
}

func NewConnector(portName string) *Connector {
	conn := new(Connector)
	conn.portName = portName

	conf := &serial.Config{Name: portName, Baud: 19200}

	p, err := serial.OpenPort(conf)
	if err != nil {
		log.Fatalf("Failed to open COM port\n%v", err)
	}
	conn.port = p

	conn.txChan = make(chan string)
	conn.rxChan = make(chan string, 5) // 損失対策 / RTSを受け取ってくれないので

	go conn.transmit()
	go conn.receive()

	return conn
}

func (conn *Connector) transmit() {
	println("Info: Start Transmitter")

	for {
		msg := <-conn.txChan

		_, err := conn.port.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
		println("Debug: Transmit: " + msg)
	}
}

func (conn *Connector) receive() {
	println("Info: Start Receiver")
	dataBuf := make([]byte, 128)
	msgBuf := make([]byte, 0, 128)

	// Arduinoのloop()的なアレ
	for {
		var msg string

		// LFまでを取得
	line:
		for {
			n, err := conn.port.Read(dataBuf)
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
		conn.rxChan <- msg
		println("Debug: Receive Message: " + msg)
	}
}

// Getter

func (conn Connector) TransmitChannel() chan<- string {
	return conn.txChan
}

func (conn Connector) ReceiveChannel() <-chan string {
	return conn.rxChan
}
