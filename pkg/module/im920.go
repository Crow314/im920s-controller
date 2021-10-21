package module

import (
	"errors"
	"strings"
	"time"
)

type Im920s struct {
	uartChannel      uartChannel
	dataReceiver     chan ReceivedData
	responseReceiver chan string
}

type uartChannel struct {
	transmitter chan<- string
	receiver    <-chan string
}

func NewIm920s(transmitter chan<- string, receiver <-chan string) *Im920s {
	im920s := new(Im920s)

	im920s.uartChannel.transmitter = transmitter
	im920s.uartChannel.receiver = receiver

	im920s.dataReceiver = make(chan ReceivedData)

	go im920s.receiver()

	return im920s
}

func (im920s *Im920s) SendCommand(msg string) (string, error) {
	if !strings.HasSuffix(msg, "\r\n") { // 末尾がCRLFでない場合
		msg += "\r\n"
	}

	im920s.responseReceiver = make(chan string)
	im920s.uartChannel.transmitter <- msg

	res := ""
	var err error

	select {
	case res = <-im920s.responseReceiver:
		if res == "NG\r\n" {
			err = errors.New("returned \"NG\" response")
		}

		if msg == "RPRM\r\n" || msg == "rprm\r\n" {
		loop:
			for { // 2行目以降に対応
				select {
				case res2 := <-im920s.responseReceiver:
					res += res2
				case <-time.After(1 * time.Second):
					break loop
				}
			}
		}
	case <-time.After(10 * time.Second):
		err = errors.New("returned no response")
	}

	close(im920s.responseReceiver)

	return res, err
}

// Getter

func (im920s Im920s) DataReceiver() <-chan ReceivedData {
	return im920s.dataReceiver
}
