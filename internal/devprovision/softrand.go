package devprovision

import "time"

type softRand struct {
	currentValue uint32
}

func (r *softRand) Get() uint32 {
	for r.currentValue == 0 {
		x := time.Now().UnixNano()
		r.currentValue = uint32((x >> 32) ^ x)
	}

	// Reference: https://en.wikipedia.org/wiki/Xorshift
	x := r.currentValue
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.currentValue = x
	return x
}
