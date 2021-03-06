package module

import "strconv"

func parseByteHex(rssi string) (byte, error) {
	hex, err := strconv.ParseUint(rssi, 16, 8)
	if err != nil {
		return 0, err
	}
	return byte(hex), nil
}
