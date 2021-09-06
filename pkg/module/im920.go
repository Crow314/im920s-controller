package module

type Im920s struct {
	uartChannel      UartChannel
	dataReceiver     chan ReceivedData
	responseReceiver chan string
}

type UartChannel struct {
	transmitter chan string
	receiver    chan string
}

func NewIm920s() *Im920s {
	im920s := new(Im920s)

	im920s.uartChannel.transmitter = make(chan string)
	im920s.uartChannel.receiver = make(chan string)

	return im920s
}

func (im920s *Im920s) SendCommand(msg string) {
	im920s.uartChannel.transmitter <- msg
}

// Getter

func (im920s Im920s) UartTransmitter() chan string {
	return im920s.uartChannel.transmitter
}

func (im920s Im920s) UartReceiver() chan string {
	return im920s.uartChannel.receiver
}

func (im920s Im920s) DataReceiver() chan ReceivedData {
	return im920s.dataReceiver
}
