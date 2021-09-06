package module

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ReceivedData struct {
	node uint16
	rssi int8
	data []byte
}

var receiveDataFormat *regexp.Regexp

func (im920s *Im920s) receiver() {
	for {
		str := <-im920s.uartChannel.receiver

		// onReceived的なsomething

		if isReceivedData(str) {
			data, err := parseData(str)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to parse data\n%v\n", err)
				continue
			}

			im920s.dataReceiver <- *data
		} else {
			im920s.responseReceiver <- str
		}
	}
}

func isReceivedData(receivedStr string) bool {
	// 正規表現定義
	if receiveDataFormat == nil {
		receiveDataFormat = regexp.MustCompile("^00,[0-9A-F]{4},[0-9A-F]{2}:.+\r\n$")
	}

	return receiveDataFormat.MatchString(receivedStr)
}

func parseData(dataStr string) (*ReceivedData, error) {
	// 末尾CRLF除去
	if strings.HasSuffix(dataStr, "\r\n") {
		dataStr = dataStr[:len(dataStr)-2]
	}

	data := new(ReceivedData)

	tmpStrs := strings.Split(dataStr, ":")

	if len(tmpStrs) != 2 {
		return nil, errors.New("invalid format\nCan't separate header and body.")
	}

	header := strings.Split(tmpStrs[0], ",")

	if header[0] != "00" {
		return nil, errors.New("invalid format\nDummy is not 0x00.")
	}

	// Node
	hex, err := strconv.ParseInt(header[1], 16, 16)
	if err != nil {
		return nil, err
	}
	data.node = uint16(hex)

	// RSSI
	tmpByte, err := parseByteHex(header[2])
	if err != nil {
		return nil, err
	}
	data.rssi = int8(tmpByte)

	return data, nil
}

// Getter

func (data *ReceivedData) Node() uint16 {
	return data.node
}

func (data *ReceivedData) Rssi() int8 {
	return data.rssi
}

func (data *ReceivedData) Data() []byte {
	return data.data
}