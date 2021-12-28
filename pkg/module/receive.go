package module

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ReceivedData struct {
	node uint16
	rssi int8
	data []byte
}

var receiveDataFormat *regexp.Regexp

func (im920s *Im920s) receiver(cmdResponseChan chan string) {
	for {
		str := <-im920s.uartChannel.receiver

		// onReceived的なsomething

		if isReceivedData(str) { // データ受信
			data, err := parseData(str)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: Failed to parse data\n%v\n", err)
				continue
			}

			select {
			case im920s.dataReceiver <- *data:
			case <-time.After(10 * time.Second): // Timeout
			}

		} else {
			if str == "GRNOREGD\r\n" { // STGNコマンド実行後 グループ番号設定パケット受信時
				// TODO config struct
				println("Info: Group number has been registered")
			} else { // コマンドに対するレスポンス
				cmdResponseChan <- str
			}
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
	// TODO ECIO
	body := strings.Split(tmpStrs[1], ",") // Only DCIO mode

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

	// Data
	data.data = make([]byte, len(body))
	for i, v := range body {
		res, err := parseByteHex(v)
		if err != nil {
			return nil, err
		}
		data.data[i] = res
	}

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
