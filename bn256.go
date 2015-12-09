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
#include "curvepoint_fp.h"
#include "twistpoint_fp2.h"

// Forward declaration of constants defined in parameters.c
const scalar_t bn_v_scalar;
const curvepoint_fp_t bn_curvegen;
const twistpoint_fp2_t bn_twistgen;

// TODO: Remove this.
void runBenchmarkG1ScalarMult(int N, const curvepoint_fp_t op, const scalar_t s) {
	curvepoint_fp_t rop;
	int i;
	for (i = 0; i < N; i++) {
		curvepoint_fp_scalarmult_vartime(rop, op, s);
	}
}
void runBenchmarkG2ScalarMult(int N, const twistpoint_fp2_t op, const scalar_t s) {
	twistpoint_fp2_t rop;
	int i;
	for (i = 0; i < N; i++) {
		twistpoint_fp2_scalarmult_vartime(rop, op, s);
	}
}
*/
import "C"

var (
	v      = new(big.Int)
	p      *big.Int // The prime: 36u⁴+36u³+24u³+6u+1, where u=v³
	big6v  *big.Int // 6v
	baseG1 = new(G1)
	baseG2 = new(G2)
)

const numBytes = 32

func init() {
	scalar2big(v, &C.bn_v_scalar)
	p, _ = new(big.Int).SetString("65000549695646603732796438742359905742825358107623003571877145026864184071783", 10)
	C.curvepoint_fp_set(&baseG1.p, &C.bn_curvegen[0])
	C.twistpoint_fp2_set(&baseG2.p, &C.bn_twistgen[0])
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
	b, tmp := e.Marshal(), new(big.Int)
	return fmt.Sprintf("bn256.G1(%v, %v)",
		tmp.SetBytes(b[:numBytes]).String(),
		tmp.SetBytes(b[numBytes:]).String())
}

func (e *G1) Marshal() []byte {
	C.curvepoint_fp_makeaffine(&e.p)
	tmp := new(big.Int)
	ret := make([]byte, numBytes*2)
	putBigBytes(ret, 0, fpe2big(tmp, e.p.m_x))
	putBigBytes(ret, 1, fpe2big(tmp, e.p.m_y))
	return ret
}

func (e *G1) Unmarshal(m []byte) (*G1, bool) {
	if len(m) != numBytes*2 {
		return e, false
	}
	tmp := new(big.Int)
	big2fpe(&e.p.m_x, tmp.SetBytes(m[0:numBytes]))
	big2fpe(&e.p.m_y, tmp.SetBytes(m[numBytes:]))
	C.fpe_setone(&e.p.m_z[0])
	C.fpe_setzero(&e.p.m_t[0])
	return e, true
}

type G2 struct {
	p C.struct_twistpoint_fp2_struct
}

func (e *G2) Add(a, b *G2) *G2 {
	C.twistpoint_fp2_add_vartime(&e.p, &a.p, &b.p)
	return e
}

func (e *G2) ScalarMult(base *G2, k *big.Int) *G2 {
	if k.BitLen() == 0 {
		C.fp2e_setzero(&e.p.m_z[0])
		return e
	}
	var ck [4]C.ulonglong
	big2scalar(&ck, k)
	C.twistpoint_fp2_makeaffine(&base.p)
	C.twistpoint_fp2_scalarmult_vartime(&e.p, &base.p, &ck[0])
	return e
}

func (e *G2) ScalarBaseMult(k *big.Int) *G2 {
	return e.ScalarMult(baseG2, k)
}

func (e *G2) Marshal() []byte {
	C.twistpoint_fp2_makeaffine(&e.p)
	ret := make([]byte, numBytes*4)
	a, b := new(big.Int), new(big.Int)
	fp2e2big(a, b, e.p.m_x)
	putBigBytes(ret, 0, a)
	putBigBytes(ret, 1, b)
	fp2e2big(a, b, e.p.m_y)
	putBigBytes(ret, 2, a)
	putBigBytes(ret, 3, b)
	return ret
}

func (e *G2) Unmarshal(m []byte) (*G2, bool) {
	if len(m) != numBytes*4 {
		return e, false
	}
	a, b := new(big.Int), new(big.Int)
	big2fp2e(&e.p.m_x, a.SetBytes(m[0:numBytes]), b.SetBytes(m[numBytes:2*numBytes]))
	big2fp2e(&e.p.m_y, a.SetBytes(m[2*numBytes:3*numBytes]), b.SetBytes(m[3*numBytes:4*numBytes]))
	C.fp2e_setone(&e.p.m_z[0])
	C.fp2e_setone(&e.p.m_t[0])
	return e, true
}

func (e *G2) String() string {
	b, tmp := e.Marshal(), new(big.Int)
	za, zb := new(big.Int), new(big.Int)
	// For compatibility with golang.org/x/cryto/bn2567.G2.String,
	// use (0,0) for e.p.m_z if m_x and m_y are (0,0)
	if C.fp2e_iszero(&e.p.m_x[0]) == 0 || C.fp2e_iszero(&e.p.m_y[0]) == 0 {
		fp2e2big(za, zb, e.p.m_z)
	}
	return fmt.Sprintf("bn256.G2((%v,%v), (%v,%v), (%v,%v))",
		tmp.SetBytes(b[0:numBytes]).String(),
		tmp.SetBytes(b[numBytes:2*numBytes]).String(),
		tmp.SetBytes(b[2*numBytes:3*numBytes]).String(),
		tmp.SetBytes(b[3*numBytes:4*numBytes]).String(),
		za,
		zb)
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

func fpe2big(out *big.Int, in C.fpe_t) *big.Int { return doubles2big(out, &(in[0].v)) }
func big2fpe(out *C.fpe_t, in *big.Int)         { big2doubles(&(out[0].v), in) }

func fp2e2big(a, b *big.Int, in C.fp2e_t) {
	// As per fp2e.h: Arrangement in memory: (b0, a0, b1, a1, ... b11,a11)
	var dbls [12]C.double
	for i := 1; i < 24; i += 2 {
		dbls[i>>1] = in[0].v[i]
	}
	doubles2big(a, &dbls)
	for i := 0; i < 24; i += 2 {
		dbls[i>>1] = in[0].v[i]
	}
	doubles2big(b, &dbls)
}
func big2fp2e(out *C.fp2e_t, a, b *big.Int) {
	var dbls [12]C.double
	big2doubles(&dbls, a)
	for i := 1; i < 24; i += 2 {
		out[0].v[i] = dbls[i>>1]
	}
	big2doubles(&dbls, b)
	for i := 0; i < 24; i += 2 {
		out[0].v[i] = dbls[i>>1]
	}
}

// Section 4.1 of https://cryptojedi.org/papers/dclxvi-20100714.pdf
func doubles2big(out *big.Int, in *[12]C.double) *big.Int {
	var (
		vx  = new(big.Int).Set(v)
		tmp = new(big.Int)
	)
	out.SetInt64(int64(in[0]))
	for i, f := range []int64{6, 6, 6, 6, 6, 6, 36, 36, 36, 36, 36} {
		tmp.SetInt64(int64(in[i+1]) * f)
		tmp.Mul(tmp, vx)
		out.Add(out, tmp)
		vx.Mul(vx, v)
	}
	return out.Mod(out, p)
}
func big2doubles(out *[12]C.double, in *big.Int) {
	var (
		dividend  = new(big.Int).Set(in)
		remainder = new(big.Int)
	)
	dividend.DivMod(dividend, big6v, remainder)
	out[0] = C.double(remainder.Int64())
	// 6v, 6v², 6v³ etc. till v⁶ and then 36v⁷ till 36v¹¹
	for i := 1; i < 6; i++ {
		dividend.DivMod(dividend, v, remainder)
		out[i] = C.double(remainder.Int64())
	}
	dividend.DivMod(dividend, big6v, remainder)
	out[6] = C.double(remainder.Int64())
	for i := 7; i < 11; i++ {
		dividend.DivMod(dividend, v, remainder)
		out[i] = C.double(remainder.Int64())
	}
	out[11] = C.double(dividend.Int64())
}

func putBigBytes(dst []byte, idx int, n *big.Int) {
	var (
		b     = n.Bytes()
		start = idx*numBytes + numBytes - len(b)
		limit = (idx + 1) * numBytes
	)
	copy(dst[start:limit], b)
}

// Helper function to measure "Go-overhead"
// TODO: Remove?
func benchmarkG1ScalarBaseMult(N int, k *big.Int) {
	var ck [4]C.ulonglong
	big2scalar(&ck, k)
	C.runBenchmarkG1ScalarMult(C.int(N), &baseG1.p, &ck[0])
}
func benchmarkG2ScalarBaseMult(N int, k *big.Int) {
	var ck [4]C.ulonglong
	big2scalar(&ck, k)
	C.runBenchmarkG2ScalarMult(C.int(N), &baseG2.p, &ck[0])
}
