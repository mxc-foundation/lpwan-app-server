// Crypto using elliptic curves defined over the finite binary field GF(2^m) where m is prime.
// The curves used are the anomalous binary curves (ABC-curves) or also called Koblitz curves.
// This class of curves was chosen because it yields efficient implementation of operations.
//
// NIST      SEC Group     strength
// K-233     sect233k1     112 bit
//
// Port from C code (https://github.com/kokke/tiny-ECDH-c)

package ecdhk223

import (
	"fmt"
)

const (
	curveDegree   = 233
	eccPrvKeySize = 32
	eccPubKeySize = (2 * eccPrvKeySize)
)

// margin for overhead needed in intermediate calculations
const (
	bitVectorMargin   = 3
	bitVectorNumBits  = (curveDegree + bitVectorMargin)
	bitVectorNumWords = ((bitVectorNumBits + 31) / 32)

//	bitVectorNumBytes = (4 * bitVectorNumWords)
)

type bitVector [bitVectorNumWords]uint32

var (
	// NIST K-233
	polynomial = bitVector{0x00000001, 0x00000000, 0x00000400, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000200}
	coeffB     = bitVector{0x00000001, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000}
	baseX      = bitVector{0xefad6126, 0x0a4c9d6e, 0x19c26bf5, 0x149563a4, 0x29f22ff4, 0x7e731af1, 0x32ba853a, 0x00000172}
	baseY      = bitVector{0x56fae6a3, 0x56e0c110, 0xf18aeb9b, 0x27a8cd9b, 0x555a67c4, 0x19b7f70f, 0x537dece8, 0x000001db}
	baseOrder  = bitVector{0xf173abdf, 0x6efb1ad5, 0xb915bcd4, 0x00069d5b, 0x00000000, 0x00000000, 0x00000000, 0x00000080}
)

func (p bitVector) String() string {
	var str string
	for i := 0; i < bitVectorNumWords; i++ {
		strValue := fmt.Sprintf(" 0x%08X", p[i])
		str += strValue
	}
	return str
}

//==========================================================================
// some basic bit-manipulation routines that act on bit-vectors follow
//==========================================================================
func bitVectorGetBit(x *bitVector, idx int) uint32 {
	return (x[idx/32] >> (idx & 31) & 1)
}

func bitVectorClearBit(x *bitVector, idx int) {
	x[idx/32] &= ^(1 << (idx & 31))
}

func bitVectorCopy(x *bitVector, y *bitVector) {
	for i := 0; i < bitVectorNumWords; i++ {
		x[i] = y[i]
	}
}

func bitVectorSwap(x *bitVector, y *bitVector) {
	var tmp bitVector
	bitVectorCopy(&tmp, x)
	bitVectorCopy(x, y)
	bitVectorCopy(y, &tmp)
}

func bitVectorIsEqual(x *bitVector, y *bitVector) bool {
	for i := 0; i < bitVectorNumWords; i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func bitVectorSetZero(x *bitVector) {
	for i := 0; i < bitVectorNumWords; i++ {
		x[i] = 0
	}
}

func bitVectorIsZero(x *bitVector) bool {
	i := uint32(0)
	for ; i < bitVectorNumWords; i++ {
		if x[i] != 0 {
			break
		}
	}
	return (i == bitVectorNumWords)
}

// return the number of the highest one-bit + 1
func bitVectorDegree(x *bitVector) int {
	i := int(bitVectorNumWords * 32)

	/* Start at the back of the vector (MSB) */
	xIdx := bitVectorNumWords - 1
	//  x += bitVectorNumWords;

	/* Skip empty / zero words */
	for (i > 0) && (x[xIdx] == 0) {
		i -= 32
		xIdx--
	}

	/* Run through rest if count is not multiple of bitsize of DTYPE */
	if i != 0 {
		u32mask := uint32(1) << 31
		for (x[xIdx] & u32mask) == 0 {
			u32mask >>= 1
			i--
		}
	}
	return i
}

// left-shift by 'count' digits
func bitVectorLeftShift(x *bitVector, y *bitVector, nbits int) {
	nwords := int(nbits / 32)

	/* Shift whole words first if nwords > 0 */
	var i, j int
	for i = 0; i < nwords; i++ {
		/* Zero-initialize from least-significant word until offset reached */
		x[i] = 0
	}
	j = 0
	/* Copy to x output */
	for i < bitVectorNumWords {
		x[i] = y[j]
		i++
		j++
	}

	/* Shift the rest if count was not multiple of bitsize of DTYPE */
	nbits &= 31
	if nbits != 0 {
		/* Left shift rest */
		for idx := (bitVectorNumWords - 1); idx > 0; idx-- {
			x[idx] = (x[idx] << nbits) | (x[idx-1] >> (32 - nbits))
		}
		x[0] <<= nbits
	}
}

//==========================================================================
// Code that does arithmetic on bit-vectors in the Galois Field GF(2^curveDegree).
//==========================================================================
func gf2FieldSetOne(x *bitVector) {
	/* Set first word to one */
	x[0] = 1
	/* .. and the rest to zero */
	for i := 1; i < bitVectorNumWords; i++ {
		x[i] = 0
	}
}

// fastest check if x == 1
func gf2FieldIsOne(x *bitVector) bool {
	/* Check if first word == 1 */
	if x[0] != 1 {
		return false
	}
	/* ...and if rest of words == 0 */
	i := 1
	for ; i < bitVectorNumWords; i++ {
		if x[i] != 0 {
			break
		}
	}
	return (i == bitVectorNumWords)
}

// galois field(2^m) addition is modulo 2, so XOR is used instead - 'z := a + b'
func gf2FieldIsAdd(z *bitVector, x *bitVector, y *bitVector) {
	for i := 0; i < bitVectorNumWords; i++ {
		z[i] = (x[i] ^ y[i])
	}
}

// increment element
// func gf2FieldIsInc(x *bitVector) {
// 	x[0] ^= 1
// }

// field multiplication 'z := (x * y)'
func gf2FieldMul(z *bitVector, x *bitVector, y *bitVector) {
	var i int
	var tmp bitVector

	if z == y {
		return
	}

	bitVectorCopy(&tmp, x)

	/* LSB set? Then start with x */
	if bitVectorGetBit(y, 0) != 0 {
		bitVectorCopy(z, x)
	} else {
		/* .. or else start with zero */
		bitVectorSetZero(z)
	}

	/* Then add 2^i * x for the rest */
	for i = 1; i < curveDegree; i++ {
		/* lshift 1 - doubling the value of tmp */
		bitVectorLeftShift(&tmp, &tmp, 1)

		/* Modulo reduction polynomial if degree(tmp) > curveDegree */
		if bitVectorGetBit(&tmp, curveDegree) == 1 {
			gf2FieldIsAdd(&tmp, &tmp, &polynomial)
		}

		/* Add 2^i * tmp if this factor in y is non-zero */
		if bitVectorGetBit(y, i) == 1 {
			gf2FieldIsAdd(z, z, &tmp)
		}
	}
}

// field inversion 'z := 1/x'
func gf2FieldInvert(z *bitVector, x *bitVector) {
	var u, v, g, h bitVector
	var i int

	bitVectorCopy(&u, x)
	bitVectorCopy(&v, &polynomial)
	bitVectorSetZero(&g)
	gf2FieldSetOne(z)

	for !gf2FieldIsOne(&u) {
		i = (bitVectorDegree(&u) - bitVectorDegree(&v))

		if i < 0 {
			bitVectorSwap(&u, &v)
			bitVectorSwap(&g, z)
			i = -i
		}
		bitVectorLeftShift(&h, &v, i)
		gf2FieldIsAdd(&u, &u, &h)
		bitVectorLeftShift(&h, &g, i)
		gf2FieldIsAdd(z, z, &h)
	}
}

//==========================================================================
// The following code takes care of Galois-Field arithmetic.
// Elliptic curve points are represented  by pairs (x,y) of bitVector.
// It is assumed that curve coefficient 'a' is {0,1}
// This is the case for all NIST binary curves.
// Coefficient 'b' is given in 'coeffB'.
// '(baseX, baseY)' is a point that generates a large prime order group.
//==========================================================================
func gf2PointCopy(x1 *bitVector, y1 *bitVector, x2 *bitVector, y2 *bitVector) {
	bitVectorCopy(x1, x2)
	bitVectorCopy(y1, y2)
}

func gf2PointSetZero(x *bitVector, y *bitVector) {
	bitVectorSetZero(x)
	bitVectorSetZero(y)
}

func gf2PointIsZero(x *bitVector, y *bitVector) bool {
	return (bitVectorIsZero(x) && bitVectorIsZero(y))
}

/* double the point (x,y) */
func gf2PointDouble(x *bitVector, y *bitVector) {
	/* iff P = O (zero or infinity): 2 * P = P */
	if bitVectorIsZero(x) {
		bitVectorSetZero(y)
	} else {
		var l bitVector

		gf2FieldInvert(&l, x)
		gf2FieldMul(&l, &l, y)
		gf2FieldIsAdd(&l, &l, x)
		gf2FieldMul(y, x, x)
		gf2FieldMul(x, &l, &l)
		gf2FieldIsAdd(x, x, &l)
		gf2FieldMul(&l, &l, x)
		gf2FieldIsAdd(y, y, &l)
	}
}

/* add two points together (x1, y1) := (x1, y1) + (x2, y2) */
func gf2PointAdd(x1 *bitVector, y1 *bitVector, x2 *bitVector, y2 *bitVector) {
	if !gf2PointIsZero(x2, y2) {
		if gf2PointIsZero(x1, y1) {
			gf2PointCopy(x1, y1, x2, y2)
		} else {
			if bitVectorIsEqual(x1, x2) {
				if bitVectorIsEqual(y1, y2) {
					gf2PointDouble(x1, y1)
				} else {
					gf2PointSetZero(x1, y1)
				}
			} else {
				/* Arithmetic with temporary variables */
				var a, b, c, d bitVector

				gf2FieldIsAdd(&a, y1, y2)
				gf2FieldIsAdd(&b, x1, x2)
				gf2FieldInvert(&c, &b)
				gf2FieldMul(&c, &c, &a)
				gf2FieldMul(&d, &c, &c)
				gf2FieldIsAdd(&d, &d, &c)
				gf2FieldIsAdd(&d, &d, &b)
				gf2FieldIsAdd(x1, x1, &d)
				gf2FieldMul(&a, x1, &c)
				gf2FieldIsAdd(&a, &a, &d)
				gf2FieldIsAdd(y1, y1, &a)
				bitVectorCopy(x1, &d)
			}
		}
	}
}

/* point multiplication via double-and-add algorithm */
func gf2PointMul(x *bitVector, y *bitVector, exp *bitVector) {
	var tmpx, tmpy bitVector
	var i int
	nbits := bitVectorDegree(exp)

	gf2PointSetZero(&tmpx, &tmpy)

	for i = (nbits - 1); i >= 0; i-- {
		gf2PointDouble(&tmpx, &tmpy)
		if bitVectorGetBit(exp, i) == 1 {
			gf2PointAdd(&tmpx, &tmpy, x, y)
		}
	}
	gf2PointCopy(x, y, &tmpx, &tmpy)
}

/* check if y^2 + x*y = x^3 + a*x^2 + coeffB holds */
func gf2PointIsOnCurve(x *bitVector, y *bitVector) bool {
	var a, b bitVector

	if gf2PointIsZero(x, y) {
		return true
	}
	gf2FieldMul(&a, x, x)
	gf2FieldMul(&a, &a, x)
	gf2FieldIsAdd(&a, &a, &coeffB)
	gf2FieldMul(&b, y, y)
	gf2FieldIsAdd(&a, &a, &b)
	gf2FieldMul(&b, x, y)

	return bitVectorIsEqual(&a, &b)
}

// Copy between uint32 and byte
func bytesToBitVector(aDest *bitVector, aSrc *[eccPrvKeySize]byte) {
	for i := 0; i < bitVectorNumWords; i++ {
		value := uint32(0)
		offset := i * 4
		value = uint32(aSrc[offset+0]) | (uint32(aSrc[offset+1]) << 8) | (uint32(aSrc[offset+2]) << 16) | (uint32(aSrc[offset+3]) << 24)
		aDest[i] = value
	}
}

func bitVectorToBytes(aDest *[eccPrvKeySize]byte, aSrc *bitVector) {
	for i := 0; i < bitVectorNumWords; i++ {
		value := aSrc[i]
		offset := i * 4

		aDest[offset+0] = uint8(value >> 0)
		aDest[offset+1] = uint8(value >> 8)
		aDest[offset+2] = uint8(value >> 16)
		aDest[offset+3] = uint8(value >> 24)
	}
}

//==========================================================================
//==========================================================================
// func dumpBitVector(aPrefix string, aSrc *bitVector) {
// 	str := aPrefix
// 	for i := range aSrc {
// 		str += fmt.Sprintf(" %08X", aSrc[i])
// 	}
// 	fmt.Println(str)
// }

//==========================================================================
// Elliptic Curve Diffie-Hellman key exchange protocol.
//==========================================================================
// NOTE: private should contain random data a-priori!
func ecdhGenerateKeysK223(publickey *[eccPubKeySize]byte, privatekey *[eccPrvKeySize]byte) int {
	var pub1, pub2, priv bitVector
	bytesToBitVector(&priv, privatekey)

	// dumpBitVector("priv:", &priv)

	/* Get copy of "base" point 'G' */
	gf2PointCopy(&pub1, &pub2, &baseX, &baseY)

	/* Abort key generation if random number is too small */
	if bitVectorDegree(&priv) < (curveDegree / 2) {
		return 0
	}
	/* Clear bits > curveDegree in highest word to satisfy constraint 1 <= exp < n. */
	nbits := bitVectorDegree(&baseOrder)

	for i := (nbits - 1); i < (bitVectorNumWords * 32); i++ {
		bitVectorClearBit(&priv, i)
	}

	/* Multiply base-point with scalar (private-key) */
	gf2PointMul(&pub1, &pub2, &priv)

	// dumpBitVector("pub1:", &pub1)
	// dumpBitVector("pub2:", &pub2)

	// Copy result
	var tmp [eccPrvKeySize]byte
	bitVectorToBytes(&tmp, &pub1)
	copy(publickey[0:], tmp[:])
	bitVectorToBytes(&tmp, &pub2)
	copy(publickey[eccPrvKeySize:], tmp[0:eccPrvKeySize])
	bitVectorToBytes(&tmp, &priv)
	copy(privatekey[:], tmp[:])

	return 1
}

func ecdhSharedSecretK223(privatekey *[eccPrvKeySize]byte, otherspub *[eccPubKeySize]byte, output *[eccPubKeySize]byte) int {
	var otherPub1, otherPub2, priv bitVector
	var output1, output2 bitVector
	var tmp [eccPrvKeySize]byte

	bytesToBitVector(&priv, privatekey)
	copy(tmp[:], otherspub[0:])
	bytesToBitVector(&otherPub1, &tmp)
	copy(tmp[:], otherspub[eccPrvKeySize:])
	bytesToBitVector(&otherPub2, &tmp)

	/* Do some basic validation of other party's public key */
	if !gf2PointIsZero(&otherPub1, &otherPub2) && gf2PointIsOnCurve(&otherPub1, &otherPub2) {
		/* Copy other side's public key to output */
		// for i := 0; i < (bitVectorNumBytes * 2); i++ {
		// 	output[i] = others_pub[i]
		// }
		output1 = otherPub1
		output2 = otherPub2

		// dumpBitVector("priv:", &priv)
		// dumpBitVector("output1:", &output1)
		// dumpBitVector("output2:", &output2)

		/* Multiply other side's public key with own private key */
		gf2PointMul(&output1, &output2, &priv)

		//
		bitVectorToBytes(&tmp, &output1)
		copy(output[0:], tmp[:])
		bitVectorToBytes(&tmp, &output2)
		copy(output[eccPrvKeySize:], tmp[0:eccPrvKeySize])

		return 1
	}

	return 0
}
