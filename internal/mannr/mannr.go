package mannr

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func Serial2Manufacturer(serial string) string {
	hash := sha256.Sum256([]byte(serial))
	str := strings.ToUpper(hex.EncodeToString(hash[:]))
	return str[0:24]
}
