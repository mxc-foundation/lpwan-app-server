package types

import (
	"encoding/hex"
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
