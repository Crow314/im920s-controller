package module

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
	im920s.responseReceiver = make(chan string)

	go im920s.receiver()

	return im920s
}

func (im920s *Im920s) SendCommand(msg string) {
	im920s.uartChannel.transmitter <- msg
}

// Getter

func (im920s Im920s) DataReceiver() <-chan ReceivedData {
	return im920s.dataReceiver
}

func (im920s Im920s) MessageReceiver() <-chan string {
	return im920s.responseReceiver
}
