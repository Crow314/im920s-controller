package connector

import (
	"github.com/tarm/serial"
	"log"
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
	for {
		msg := <-conn.txChan

		_, err := conn.port.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (conn *Connector) receive() {
	dataBuf := make([]byte, 128)
	msgBuf := make([]byte, 0, 128)

	for {
		n, err := conn.port.Read(dataBuf)
		if err != nil {
			log.Fatal(err)
		}

		// 1回での受信データ (!= 1行)
		for _, v := range dataBuf[:n] {
			msgBuf = append(msgBuf, v)

			if v == '\n' { // LFで一区切り
				msg := string(msgBuf)
				msgBuf = msgBuf[:0] // 要素を全て削除

				conn.rxChan <- msg
			}
		}
	}
}

// Getter

func (conn Connector) TransmitChannel() chan<- string {
	return conn.txChan
}

func (conn Connector) ReceiveChannel() <-chan string {
	return conn.rxChan
}
