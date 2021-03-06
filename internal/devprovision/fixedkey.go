package devprovision

import (
	"crypto/aes"

	log "github.com/sirupsen/logrus"
)

var encFixedKey = []byte{
	0xfb, 0x4a, 0xe7, 0x76, 0x80, 0x3a, 0x54, 0xd7, 0xd9, 0xb5, 0x64, 0x2d, 0x25, 0x89, 0xdc, 0xbb}

var epromKey = []byte{
	0x54, 0x30, 0xc9, 0xe7, 0xdb, 0x00, 0x39, 0xcb,
	0xf6, 0x43, 0xd1, 0xdc, 0x9b, 0x54, 0xf3, 0x71}

func getFixedKey() []byte {
	output := make([]byte, 16)

	cipher, err := aes.NewCipher(epromKey)
	if err != nil {
		log.Error("aes.NewCipher() failed.")
	} else {
		cipher.Decrypt(output, encFixedKey)
	}
	return output
}
