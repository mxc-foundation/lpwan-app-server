package ecdhk223

import (
	"testing"
)

type keyPair struct {
	publicKey  [eccPubKeySize]byte
	privateKey [eccPrvKeySize]byte
	sharedKey  [eccPubKeySize]byte
}

func TestBitManipulation(t *testing.T) {
	var arrayTmp bitVector
	arrayA := bitVector{
		0x00000000, 0x00000001, 0x00000002, 0x00000003, 0x00000004, 0x00000005, 0x00000006, 0x00000007,
	}
	arrayB := bitVector{
		0x12345678, 0x0000000A, 0x0000000B, 0x0000000C, 0x0000000D, 0x0000000E, 0x0000000F, 0x55555555,
	}
	arrayZero := bitVector{
		0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
	}
	arrayAShift1 := bitVector{
		0x00000000, 0x00000002, 0x00000004, 0x00000006, 0x00000008, 0x0000000a, 0x0000000c, 0x0000000e,
	}
	arrayAShift31 := bitVector{
		0x00000000, 0x80000000, 0x00000000, 0x80000001, 0x00000001, 0x80000002, 0x00000002, 0x80000003,
	}
	arrayAShift32 := bitVector{
		0x00000000, 0x00000000, 0x00000001, 0x00000002, 0x00000003, 0x00000004, 0x00000005, 0x00000006,
	}
	arrayAShift40 := bitVector{
		0x00000000, 0x00000000, 0x00000100, 0x00000200, 0x00000300, 0x00000400, 0x00000500, 0x00000600,
	}

	// Get Bit
	var bitValue uint32
	bitValue = bitVectorGetBit(&arrayA, 31)
	if bitValue != 0 {
		t.Errorf("bitVectorGetBit() failed at bit 31, expect=0, got=%d", bitValue)
	}
	bitValue = bitVectorGetBit(&arrayA, 32)
	if bitValue != 1 {
		t.Errorf("bitVectorGetBit() failed at bit 32, expect=1, got=%d", bitValue)
	}

	// Clear bit
	arrayTmp = arrayA
	bitVectorClearBit(&arrayTmp, 32)
	if (arrayTmp[1] & 0x01) != 0 {
		t.Errorf("bitVectorClearBit() clear bit 32 failed")
	}

	// Copy
	bitVectorCopy(&arrayTmp, &arrayA)
	if arrayA != arrayTmp {
		t.Error("bitVectorCopy() failed.")
	}

	// Swap
	arrayA2 := arrayA
	arrayB2 := arrayB
	bitVectorSwap(&arrayA2, &arrayB2)
	if (arrayA2 != arrayB) || (arrayB2 != arrayA) {
		t.Error("bitVectorSwap() failed.")
	}

	// Check equal
	arrayTmp = arrayA
	if !bitVectorIsEqual(&arrayTmp, &arrayA) {
		t.Error("bitVectorIsEqual() failed.")
	}
	if bitVectorIsEqual(&arrayTmp, &arrayB) {
		t.Error("bitVectorIsEqual() failed.")
	}

	// Zero
	arrayTmp = arrayA
	bitVectorSetZero(&arrayTmp)
	if arrayTmp != arrayZero {
		t.Error("bitVectorSetZero() failed.")
	}

	// Check zero
	arrayTmp = arrayZero
	if !bitVectorIsZero(&arrayTmp) {
		t.Error("bitVectorIsZero() failed.")
	}
	arrayTmp[7] = 1
	if bitVectorIsZero(&arrayTmp) {
		t.Error("bitVectorIsZero() failed.")
	}

	// Check highest bit
	var bitDegree int
	arrayTmp = arrayZero
	arrayTmp[1] = 0x80000000
	bitDegree = bitVectorDegree(&arrayTmp)
	if bitDegree != 64 {
		t.Errorf("bitVectorDegree() failed. expect=64, got=%d", bitDegree)
	}
	arrayTmp = arrayZero
	bitDegree = bitVectorDegree(&arrayTmp)
	if bitDegree != 0 {
		t.Errorf("bitVectorDegree() failed. expect=0, got=%d", bitDegree)
	}
	arrayTmp = arrayZero
	arrayTmp[0] = 0x00000001
	bitDegree = bitVectorDegree(&arrayTmp)
	if bitDegree != 1 {
		t.Errorf("bitVectorDegree() failed. expect=1, got=%d", bitDegree)
	}

	// Left shift bit
	bitVectorLeftShift(&arrayTmp, &arrayA, 1)
	if arrayTmp != arrayAShift1 {
		t.Error("bitVectorLeftShift() failed on 1.")
		t.Logf("  arrayA: %s", arrayA)
		t.Logf("arrayTmp: %s", arrayTmp)
	}
	bitVectorLeftShift(&arrayTmp, &arrayA, 31)
	if arrayTmp != arrayAShift31 {
		t.Error("bitVectorLeftShift() failed on 31.")
		t.Logf("  arrayA: %s", arrayA)
		t.Logf("arrayTmp: %s", arrayTmp)
	}
	bitVectorLeftShift(&arrayTmp, &arrayA, 32)
	if arrayTmp != arrayAShift32 {
		t.Error("bitVectorLeftShift() failed on 32.")
		t.Logf("  arrayA: %s", arrayA)
		t.Logf("arrayTmp: %s", arrayTmp)
	}
	bitVectorLeftShift(&arrayTmp, &arrayA, 40)
	if arrayTmp != arrayAShift40 {
		t.Error("bitVectorLeftShift() failed on 40.")
		t.Logf("  arrayA: %s", arrayA)
		t.Logf("arrayTmp: %s", arrayTmp)
	}

}

// func dumpPrivateKey(t *testing.T, prefix string, key *[eccPrvKeySize]byte) {
// 	str := prefix
// 	for i := range key {
// 		str += fmt.Sprintf(" %02X", key[i])
// 	}
// 	t.Logf(str)
// }

// func dumpPublicKey(t *testing.T, prefix string, key *[eccPubKeySize]byte) {
// 	str := prefix
// 	for i := range key {
// 		str += fmt.Sprintf(" %02X", key[i])
// 	}
// 	t.Logf(str)
// }

func TestGenerateKey(t *testing.T) {
	var genKey keyPair
	expectedKey := keyPair{
		privateKey: [eccPrvKeySize]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
			0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00},
		publicKey: [eccPubKeySize]byte{0xF5, 0xDD, 0xD2, 0xC7, 0x04, 0x92, 0xE0, 0xD6, 0xF2, 0x1F, 0x8D, 0xEC, 0xE0, 0x2D, 0x0A, 0xAF,
			0x75, 0x64, 0x78, 0xE1, 0x02, 0x09, 0x72, 0x75, 0x19, 0x5A, 0xFB, 0x9B, 0xB8, 0x01, 0x00, 0x00,
			0xB3, 0x29, 0x00, 0x02, 0x9A, 0xB4, 0xD6, 0x84, 0x1C, 0xC5, 0x2B, 0x51, 0x72, 0xEE, 0x2F, 0x3C,
			0x5A, 0x66, 0xBC, 0x6F, 0x03, 0x25, 0x3A, 0x92, 0x43, 0x9E, 0x14, 0x2F, 0x82, 0x00, 0x00, 0x00},
	}

	// Init keys
	for i := range genKey.privateKey {
		genKey.privateKey[i] = 0x01
	}
	for i := range genKey.publicKey {
		genKey.publicKey[i] = 0
	}
	ret := ecdhGenerateKeysK223(&genKey.publicKey, &genKey.privateKey)
	if ret == 0 {
		t.Error("ecdhGenerateKeysK223() failed.")
	}

	// dumpPrivateKey(t, "privateKey:", &genKey.privateKey)
	// dumpPublicKey(t, "publicKey:", &genKey.publicKey)
	if genKey.privateKey != expectedKey.privateKey {
		t.Error("Unexpected private key generated.")
	}
	if genKey.publicKey != expectedKey.publicKey {
		t.Error("Unexpected public key generated.")
	}
}

func TestSharedSecret(t *testing.T) {
	var alice keyPair
	var bob keyPair
	var ret int

	expectedSharedKey := [eccPubKeySize]byte{
		0x57, 0x57, 0x3A, 0x81, 0xE2, 0x7E, 0x48, 0x26, 0xFA, 0x8E, 0x18, 0x70, 0xCD, 0x6B, 0x66, 0x40,
		0xF3, 0x90, 0x5D, 0x98, 0x40, 0xF4, 0x12, 0xFA, 0xAE, 0x74, 0x0B, 0x12, 0xE0, 0x01, 0x00, 0x00,
		0xC4, 0xD8, 0x27, 0xA9, 0x37, 0x49, 0xEE, 0x44, 0xEA, 0x1B, 0xAC, 0x1C, 0x18, 0x8C, 0x03, 0xAA,
		0x6B, 0x02, 0xDA, 0x1C, 0x68, 0xE9, 0xE8, 0xE6, 0xCA, 0xB9, 0xD1, 0xED, 0x91, 0x01, 0x00, 0x00}

	// Set private key
	for i := range alice.privateKey {
		alice.privateKey[i] = 0x01
	}
	for i := range alice.publicKey {
		alice.publicKey[i] = 0
	}
	for i := range bob.privateKey {
		bob.privateKey[i] = 0x02
	}
	for i := range bob.publicKey {
		bob.publicKey[i] = 0
	}

	// Gen public key
	ret = ecdhGenerateKeysK223(&alice.publicKey, &alice.privateKey)
	if ret == 0 {
		t.Error("ecdhGenerateKeysK223() failed.")
	}
	ret = ecdhGenerateKeysK223(&bob.publicKey, &bob.privateKey)
	if ret == 0 {
		t.Error("ecdhGenerateKeysK223() failed.")
	}

	// Gen shared key
	ret = ecdhSharedSecretK223(&alice.privateKey, &bob.publicKey, &alice.sharedKey)
	if ret == 0 {
		t.Error("ecdhSharedSecretK223() failed.")
	}
	ret = ecdhSharedSecretK223(&bob.privateKey, &alice.publicKey, &bob.sharedKey)
	if ret == 0 {
		t.Error("ecdhSharedSecretK223() failed.")
	}

	// dumpPublicKey(t, "alice:", &alice.sharedKey)
	// dumpPublicKey(t, "bob:", &bob.sharedKey)
	if alice.sharedKey != bob.sharedKey {
		t.Error("The output shared key are mismatch.")
	}
	if alice.sharedKey != expectedSharedKey {
		t.Error("Unexpected shared key generated.")
	}

}
