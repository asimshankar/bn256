// Package bn256 is a drop in replacement for golang.org/x/crypto/bn256.
//
// It uses cgo to wrap over the fast C implementation by Michael Naehrig, Ruben
// Niederhagen, and Peter Schwabe (See
// https://www.cryptojedi.org/crypto/#dclxvi and "New software speed records
// for cryptographic pairings." published in LATINCRYPT 2010).
package bn256

import (
	"fmt"
	"math/big"
)

// #cgo CFLAGS: -I./dclxvi-20130329
// #cgo LDFLAGS: -L./dclxvi-20130329 -ldclxvi -lm
//
/*
#include <string.h>
#include "curvepoint_fp.h"

// Forward declaration of constants defined in parameters.c
const scalar_t bn_v_scalar;
const curvepoint_fp_t bn_curvegen;

void fpe_v_memcpy(fpe_t *out, const double in[12]) {
	memcpy(out, in, sizeof(in)*12);
}

// TODO: Remove this.
void runBenchmarkScalarMult(int N, const curvepoint_fp_t op, const scalar_t s) {
	curvepoint_fp_t rop;
	int i;
	for (i = 0; i < N; i++) {
		curvepoint_fp_scalarmult_vartime(rop, op, s);
	}
}
*/
import "C"

var (
	v      = new(big.Int)
	p      *big.Int // The prime: 36u⁴+36u³+24u³+6u+1, where u=v³
	big6v  *big.Int
	baseG1 *G1
)

const numBytes = 32

func init() {
	scalar2big(v, &C.bn_v_scalar)
	p, _ = new(big.Int).SetString("65000549695646603732796438742359905742825358107623003571877145026864184071783", 10)
	baseG1 = new(G1)
	C.curvepoint_fp_set(&baseG1.p, &C.bn_curvegen[0])
	big6v = new(big.Int).Mul(big.NewInt(6), v)
}

type G1 struct {
	p C.struct_curvepoint_fp_struct
}

func (e *G1) Add(a, b *G1) *G1 {
	// TODO: Does b.p need to be in affine coordinates?
	// The commented out curvepoint_fp_mixadd requires that, but I'm not sure if that comment extends to
	// curvepoint_fp_add_vartime
	C.curvepoint_fp_add_vartime(&e.p, &a.p, &b.p)
	return e
}

func (e *G1) Neg(a *G1) *G1 {
	C.curvepoint_fp_neg(&e.p, &a.p)
	return e
}

func (e *G1) ScalarBaseMult(k *big.Int) *G1 {
	return e.ScalarMult(baseG1, k)
}

func (e *G1) ScalarMult(base *G1, k *big.Int) *G1 {
	if k.BitLen() == 0 {
		return e
	}
	var ck [4]C.ulonglong
	big2scalar(&ck, k)
	C.curvepoint_fp_makeaffine(&base.p)
	C.curvepoint_fp_scalarmult_vartime(&e.p, &base.p, &ck[0])
	return e
}

func (e *G1) String() string {
	C.curvepoint_fp_makeaffine(&e.p)
	tmp := new(big.Int)
	return "bn256.G1(" + fp2big(tmp, e.p.m_x).String() + ", " + fp2big(tmp, e.p.m_y).String() + ")"
}

func (e *G1) Marshal() []byte {
	C.curvepoint_fp_makeaffine(&e.p)
	tmp := new(big.Int)
	x, y := fp2big(tmp, e.p.m_x).Bytes(), fp2big(tmp, e.p.m_y).Bytes()
	ret := make([]byte, numBytes*2)
	copy(ret[0+(numBytes-len(x)):numBytes], x)
	copy(ret[numBytes+(numBytes-len(y)):], y)
	return ret
}

func (e *G1) Unmarshal(m []byte) (*G1, bool) {
	if len(m) != numBytes*2 {
		return e, false
	}
	tmp := new(big.Int)
	big2fp(&e.p.m_x, tmp.SetBytes(m[0:numBytes]))
	big2fp(&e.p.m_y, tmp.SetBytes(m[numBytes:]))
	C.fpe_setone(&e.p.m_z[0])
	C.fpe_setone(&e.p.m_t[0])
	return e, true
}

func big2scalar(out *[4]C.ulonglong, in *big.Int) error {
	// TODO: Not portable on 32-bit architectures. Use in.Bytes there?
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

// Section 4.1 of https://cryptojedi.org/papers/dclxvi-20100714.pdf
func fp2big(out *big.Int, in C.fpe_t) *big.Int {
	var (
		vx   = new(big.Int).Set(v)
		tmp  = new(big.Int)
		dbls = &(in[0].v)
	)
	out.SetInt64(int64(dbls[0]))
	for i, f := range []int64{6, 6, 6, 6, 6, 6, 36, 36, 36, 36, 36} {
		tmp.SetInt64(int64(dbls[i+1]) * f)
		tmp.Mul(tmp, vx)
		out.Add(out, tmp)
		vx.Mul(vx, v)
	}
	return out.Mod(out, p)
}
func big2fp(out *C.fpe_t, in *big.Int) {
	var (
		dividend  = new(big.Int).Set(in)
		remainder = new(big.Int)
		dbls      [12]C.double
	)
	dividend.DivMod(dividend, big6v, remainder)
	dbls[0] = C.double(remainder.Int64())
	// 6v, 6v², 6v³ etc. till v⁶ and then 36v⁷ till 36v¹¹
	for i := 1; i < 6; i++ {
		dividend.DivMod(dividend, v, remainder)
		dbls[i] = C.double(remainder.Int64())
	}
	dividend.DivMod(dividend, big6v, remainder)
	dbls[6] = C.double(remainder.Int64())
	for i := 7; i < 11; i++ {
		dividend.DivMod(dividend, v, remainder)
		dbls[i] = C.double(remainder.Int64())
	}
	dbls[11] = C.double(dividend.Int64())
	C.fpe_v_memcpy(out, &dbls[0])
}

// Helper function to measure "Go-overhead"
// TODO: Remove?
func benchmarkScalarBaseMultC(N int, k *big.Int) {
	var ck [4]C.ulonglong
	big2scalar(&ck, k)
	C.runBenchmarkScalarMult(C.int(N), &baseG1.p, &ck[0])
}
