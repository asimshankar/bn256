// +build arm 386

package bn256

import (
	"C"
	"fmt"
	"math/big"
)

func big2scalar(out *[4]C.ulonglong, in *big.Int) error {
	b := in.Bits()
	if len(b) > 8 {
		return fmt.Errorf("big.Int needs %d words, cannot be converted to scalar_t", len(b))
	}
	max := len(b) >> 1
	for i := 0; i < max; i++ {
		out[i] = C.ulonglong(b[i<<1]) | (C.ulonglong(b[i<<1+1]) << 32)
	}
	if len(b)&0x1 == 1 {
		out[max] = C.ulonglong(b[len(b)-1])
	}
	return nil
}

func scalar2big(out *big.Int, in *[4]C.ulonglong) {
	bits := make([]big.Word, 8)
	for i := 0; i < len(bits); i += 2 {
		bits[i] = big.Word(in[i>>1] & 0x00000000ffffffff)
		bits[i+1] = big.Word((in[i>>1] & 0xffffffff00000000) >> 32)
	}
	out.SetBits(bits)
}
