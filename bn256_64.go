// +build amd64

package bn256

import (
	"C"
	"fmt"
	"math/big"
)

func big2scalar(out *[4]C.ulonglong, in *big.Int) error {
	b := in.Bits()
	if len(b) > 4 {
		return fmt.Errorf("big.Int needs %d words, cannot be converted to scalar_t", len(b))
	}
	for i, w := range b {
		out[i] = C.ulonglong(w)
	}
	return nil
}

func scalar2big(out *big.Int, in *[4]C.ulonglong) {
	bits := make([]big.Word, 4)
	for i := range bits {
		bits[i] = big.Word(in[i])
	}
	out.SetBits(bits)
}
