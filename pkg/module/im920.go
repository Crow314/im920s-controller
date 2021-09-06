package module

type Im920s struct {
	UartChannel UartChannel
}

type UartChannel struct {
	Transmitter chan string
	Receiver    chan string
}

func NewIm920s() *Im920s {
	im920s := new(Im920s)

	im920s.UartChannel.Transmitter = make(chan string)
	im920s.UartChannel.Receiver = make(chan string)

	return im920s
}

func (im920s *Im920s) SendCommand(msg string) {
	im920s.UartChannel.Transmitter <- msg
}
