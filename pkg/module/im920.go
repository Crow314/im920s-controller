package module

type Im920s struct {
	TransmitChannel chan string
	ReceiveChannel  chan string
}

func NewIm920s() *Im920s {
	im920s := new(Im920s)

	im920s.TransmitChannel = make(chan string)
	im920s.ReceiveChannel = make(chan string)

	return im920s
}

func (im920s *Im920s) SendCommand(msg string) {
	im920s.TransmitChannel <- msg
}
