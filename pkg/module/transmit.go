package module

import (
	"encoding/hex"
	"errors"
)

func (im920s *Im920s) Broadcast(data []byte) error {
	size := len(data)
	if size < 1 || size > 32 {
		return errors.New("invalid data size")
	}

	_, err := im920s.SendCommand("TXDA " + hex.EncodeToString(data))
	if err != nil {
		return err
	}

	return nil
}
