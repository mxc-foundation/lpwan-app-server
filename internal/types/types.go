package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

type MD5SUM [16]byte

// MarshalText implements encoding.TextMarshaler.
func (e MD5SUM) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (e *MD5SUM) UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(e) != len(b) {
		return fmt.Errorf("lorawan: exactly %d bytes are expected", len(e))
	}
	copy(e[:], b)
	return nil
}

// String implement fmt.Stringer.
func (e MD5SUM) String() string {
	return hex.EncodeToString(e[:])
}

// Scan implements sql.Scanner.
func (e *MD5SUM) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan md5sum: []byte type expected")
	}
	var tmp = []uint8{0}
	if bytes.Equal(b[:], tmp[:]) == false && len(b) != len(e) {
		return fmt.Errorf("Scan md5sum: []byte must have length %d", len(e))
	}
	copy(e[:], b)
	return nil
}