package devprovision

import (
	"crypto/aes"
	"time"

	"github.com/jacobsa/crypto/cmac"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/devprovision/ecdh"
)

// Device session

type deviceSession struct {
	rDevEui          []byte
	serverNonce      []byte
	devNonce         []byte
	devicePublicKey  []byte
	serverPublicKey  []byte
	serverPrivateKey []byte
	sharedKey        []byte
	assignedDevEui   []byte
	assignedAppEui   []byte
	appKey           []byte
	nwkKey           []byte
	provKey          []byte
	expireTime       time.Time
}

//
func makeDeviceSession() deviceSession {
	session := deviceSession{}
	session.rDevEui = make([]byte, 8)
	session.serverNonce = make([]byte, 4)
	session.devNonce = make([]byte, 4)
	session.devicePublicKey = make([]byte, ecdh.K233PubKeySize)
	session.serverPublicKey = make([]byte, ecdh.K233PubKeySize)
	session.serverPrivateKey = make([]byte, ecdh.K233PrvKeySize)
	session.sharedKey = make([]byte, ecdh.K233PubKeySize)
	session.expireTime = funcGetNow().Add(deviceSessionLifeTime)
	session.assignedDevEui = make([]byte, 8)
	session.assignedAppEui = make([]byte, 8)
	session.appKey = make([]byte, 16)
	session.nwkKey = make([]byte, 16)
	session.provKey = make([]byte, 16)

	return session
}

//
func (d *deviceSession) genServerKeys() {
	randbuf := funcGen128Rand()
	privateKey, publickey := ecdhK223.GenerateKeys(randbuf)
	if privateKey != nil {
		copy(d.serverPrivateKey[:], privateKey[:])
		copy(d.serverPublicKey[:], publickey[:])
	}
	copy(d.serverNonce[0:], randbuf[ecdh.K233PrvKeySize:])
}

func (d *deviceSession) genSharedKey() {
	newsharedkey := ecdhK223.SharedSecret(d.serverPrivateKey, d.devicePublicKey)
	copy(d.sharedKey[:], newsharedkey[:])
}

func (d *deviceSession) deriveKeys() {
	aeskey := make([]byte, 16) // AES-128 key size

	// AppKey
	copy(aeskey[:], d.sharedKey[0:])
	cipher, err := aes.NewCipher(aeskey)
	if err != nil {
		log.Error("aes.NewCipher() failed.")
		return
	}
	appkeyblock := make([]byte, cipher.BlockSize())
	fillByteArray(appkeyblock, 0x01)
	copy(appkeyblock[:], d.rDevEui[:])
	cipher.Encrypt(appkeyblock, appkeyblock)
	copy(d.appKey[:], appkeyblock[:])

	// NwkKey
	copy(aeskey[:], d.sharedKey[32:])
	cipher, err = aes.NewCipher(aeskey)
	if err != nil {
		log.Error("aes.NewCipher() failed.")
		return
	}
	nwkkeyblock := make([]byte, cipher.BlockSize())
	fillByteArray(nwkkeyblock, 0x02)
	copy(nwkkeyblock[:], d.rDevEui[:])
	cipher.Encrypt(nwkkeyblock, nwkkeyblock)
	copy(d.nwkKey[:], nwkkeyblock[:])

	// ProvKey
	copy(aeskey[0:], d.sharedKey[16:24])
	copy(aeskey[8:], d.sharedKey[48:56])
	cipher, err = aes.NewCipher(aeskey)
	if err != nil {
		log.Error("aes.NewCipher() failed.")
		return
	}
	provkeyblock := make([]byte, cipher.BlockSize())
	fillByteArray(provkeyblock, 0x03)
	copy(provkeyblock[:], d.rDevEui[:])
	cipher.Encrypt(provkeyblock, provkeyblock)
	copy(d.provKey[:], provkeyblock[:])
}

func (d *deviceSession) encryptAuthPayload(payload []byte, isUplink bool) []byte {
	output := payload
	ablock := make([]byte, 16)
	sblock := make([]byte, 16)
	blockcounter := 1

	ablock[0] = 0x02
	if isUplink {
		ablock[5] = 0
	} else {
		ablock[5] = 1
	}
	copy(ablock[6:14], d.rDevEui[:])

	cipher, err := aes.NewCipher(d.provKey)
	if err != nil {
		log.Error("aes.NewCipher() failed.")
		fillByteArray(output, 0)
		return output
	}

	for payloadidx := 0; payloadidx < len(payload); payloadidx += cipher.BlockSize() {
		ablock[15] = uint8(blockcounter)

		cipher.Encrypt(sblock, ablock)
		xorlen := len(payload) - payloadidx
		if xorlen > 16 {
			xorlen = 16
		}

		for i := 0; i < xorlen; i++ {
			output[payloadidx+i] ^= sblock[i]
		}

	}

	return output
}

func (d *deviceSession) calVerifyCode(serialnumber string, useservernonce bool) []byte {
	cmacbuf := make([]byte, 16)
	calbuf := []byte(serialnumber)

	if useservernonce {
		calbuf = append(calbuf, d.serverNonce...)
	} else {
		calbuf = append(calbuf, d.devNonce...)
	}

	hash, err := cmac.New(getFixedKey())
	if err != nil {
		return cmacbuf
	}
	if _, err = hash.Write(calbuf); err != nil {
		return cmacbuf
	}

	hb := hash.Sum([]byte{})

	copy(cmacbuf[:], hb[:])
	return cmacbuf
}
