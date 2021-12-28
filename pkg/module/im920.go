package module

import (
	"errors"
	"strings"
	"time"
)

type Im920s struct {
	uartChannel        uartChannel
	dataReceiver       chan ReceivedData
	commandSendChannel chan *command
}

type uartChannel struct {
	transmitter chan<- string
	receiver    <-chan string
}

type command struct {
	message       string
	response      string
	errorResponse error
	notificator   chan struct{}
}

func NewIm920s(transmitter chan<- string, receiver <-chan string) *Im920s {
	im920s := new(Im920s)

	im920s.uartChannel.transmitter = transmitter
	im920s.uartChannel.receiver = receiver

	im920s.dataReceiver = make(chan ReceivedData)
	im920s.commandSendChannel = make(chan *command)

	cmdResChan := make(chan string)

	go im920s.receiver(cmdResChan)
	go im920s.commandSender(cmdResChan)

	return im920s
}

func newCommand(msg string) *command {
	cmd := new(command)
	cmd.message = msg
	cmd.notificator = make(chan struct{})

	return cmd
}

func (im920s *Im920s) SendCommand(msg string) (string, error) {
	if !strings.HasSuffix(msg, "\r\n") { // 末尾がCRLFでない場合
		msg += "\r\n"
	}

	cmd := newCommand(msg)
	im920s.commandSendChannel <- cmd
	<-cmd.notificator // 処理待機

	return cmd.response, cmd.errorResponse
}

func (im920s *Im920s) commandSender(cmdResponseChan chan string) {
	for {
		cmd := <-im920s.commandSendChannel
		msg := cmd.message
		im920s.uartChannel.transmitter <- msg

		res := ""
		var err error

		select {
		case res = <-cmdResponseChan:
			if res == "NG\r\n" {
				err = errors.New("returned \"NG\" response")
			}

			if msg == "RPRM\r\n" || msg == "rprm\r\n" {
			loop:
				for { // 2行目以降に対応
					select {
					case res2 := <-cmdResponseChan:
						res += res2
					case <-time.After(1 * time.Second):
						break loop
					}
				}
			}
		case <-time.After(10 * time.Second):
			err = errors.New("returned no response")
		}

		cmd.response = res
		cmd.errorResponse = err

		close(cmd.notificator) // 処理完了通知

		time.Sleep(5 * time.Millisecond)
	}
}

// Getter

func (im920s Im920s) DataReceiver() <-chan ReceivedData {
	return im920s.dataReceiver
}
