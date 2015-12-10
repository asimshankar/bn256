package bn256

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"golang.org/x/crypto/bn256"
)

func TestP(t *testing.T) {
	// Test that p is indeed 36u⁴+36u³+24u³+6u+1
	// (where u is v³ and v is picked up from the C-library)
	u := new(big.Int).Exp(v, big.NewInt(3), nil)
	elems := []*big.Int{
		big.NewInt(1),
		new(big.Int).Mul(big.NewInt(6), u),
		new(big.Int).Mul(big.NewInt(24), new(big.Int).Exp(u, big.NewInt(2), nil)),
		new(big.Int).Mul(big.NewInt(36), new(big.Int).Exp(u, big.NewInt(3), nil)),
		new(big.Int).Mul(big.NewInt(36), new(big.Int).Exp(u, big.NewInt(4), nil)),
	}
	expected := new(big.Int)
	for _, e := range elems {
		expected.Add(expected, e)
	}
	if expected.Cmp(p) != 0 {
		t.Fatalf("Prime should be %v, is set to %v", expected, p)
	}
}

func TestG1(t *testing.T) {
	cmp := func(got *G1, want *bn256.G1) error {
		if gotB, wantB := got.Marshal(), want.Marshal(); !bytes.Equal(gotB, wantB) {
			return fmt.Errorf("Got %v want %v", got, want)
		}
		return nil
	}
	onetest := func(k *big.Int) error {
		var (
			got  = new(G1).ScalarBaseMult(k)
			gotB = got.Marshal()
			want = new(bn256.G1).ScalarBaseMult(k)
		)
		if g, w := got.String(), want.String(); g != w {
			// TODO: Minor implementation difference causes String for
			// golang.org/x/crypto/bn256.G1 to return (1, -2) for k=1,
			// while (1, 65000549695646603732796438742359905742825358107623003571877145026864184071781)
			// for this package. The two are identical since
			// (-2 mod p) == 65000549695646603732796438742359905742825358107623003571877145026864184071781
			// So, ignore that difference.
			if k.Cmp(big.NewInt(1)) == 0 {
				w = "bn256.G1(1, 65000549695646603732796438742359905742825358107623003571877145026864184071781)"
			}
			if g != w {
				return fmt.Errorf("k=%v: String: Got %q, want %q", k, g, w)
			}
		}
		if err := cmp(got, want); err != nil {
			return fmt.Errorf("k=%v: ScalarBaseMult: %v", k, err)
		}
		if err := cmp(
			new(G1).Add(got, new(G1).ScalarBaseMult(big.NewInt(3))),
			new(bn256.G1).Add(want, new(bn256.G1).ScalarBaseMult(big.NewInt(3))),
		); err != nil {
			return fmt.Errorf("k=%v: Add: %v", k, err)
		}
		if err := cmp(new(G1).Neg(got), new(bn256.G1).Neg(want)); err != nil {
			return fmt.Errorf("k=%v: Neg: %v", k, err)
		}
		// Unmarshal and Marshal again.
		unmarshaled, ok := new(G1).Unmarshal(gotB)
		if !ok {
			return fmt.Errorf("k=%v: Unmarshal failed", k)
		}
		again := unmarshaled.Marshal()
		if !bytes.Equal(gotB, again) {
			return fmt.Errorf("k=%v: Umarshal+Marshal: Got %v, want %v", k, again, gotB)
		}
		return nil
	}
	if err := onetest(big.NewInt(0)); err != nil {
		t.Error(err)
	}
	if err := onetest(big.NewInt(1)); err != nil {
		t.Error(err)
	}
	for i := 0; i < 100; i++ {
		k, err := rand.Int(rand.Reader, p)
		if err != nil {
			t.Fatal(err)
		}
		if err := onetest(k); err != nil {
			t.Errorf("%v (random test #%d)", err, i)
		}
	}
}

func TestG2(t *testing.T) {
	cmp := func(got *G2, want *bn256.G2) error {
		if gotB, wantB := got.Marshal(), want.Marshal(); !bytes.Equal(gotB, wantB) {
			return fmt.Errorf("Got %v want %v", got, want)
		}
		return nil
	}
	onetest := func(k *big.Int) error {
		var (
			got  = new(G2).ScalarBaseMult(k)
			want = new(bn256.G2).ScalarBaseMult(k)
			gotB = got.Marshal()
		)
		if g, w := got.String(), want.String(); g != w {
			return fmt.Errorf("k=%v: String: Got %q, want %q", k, g, w)
		}
		if err := cmp(got, want); err != nil {
			return fmt.Errorf("k=%v: ScalarBaseMult: %v", k, err)
		}
		if err := cmp(
			new(G2).Add(got, new(G2).ScalarBaseMult(big.NewInt(14141))),
			new(bn256.G2).Add(want, new(bn256.G2).ScalarBaseMult(big.NewInt(14141))),
		); err != nil {
			return fmt.Errorf("k=%v: Add: %v", k, err)
		}
		// Unmarshal and Marshal again.
		unmarshaled, ok := new(G2).Unmarshal(gotB)
		if !ok {
			return fmt.Errorf("k=%v: Unmarshal failed", k)
		}
		again := unmarshaled.Marshal()
		if !bytes.Equal(gotB, again) {
			return fmt.Errorf("k=%v: Umarshal+Marshal: Got %v, want %v", k, again, gotB)
		}
		return nil
	}
	if err := onetest(big.NewInt(0)); err != nil {
		t.Error(err)
	}
	if err := onetest(big.NewInt(1)); err != nil {
		t.Error(err)
	}
	// TODO: Enable this. As of this writing, the String method wouldn't produce
	// output identical to golang.org/x/crypto/bn256.G2.String
	/*
		for i := 0; i < 100; i++ {
			k, err := rand.Int(rand.Reader, p)
			if err != nil {
				t.Fatal(err)
			}
			if err := onetest(k); err != nil {
				t.Errorf("%v (random test #%d)", err, i)
			}
		}
	*/
}

func TestPair(t *testing.T) {
	cmp := func(got *GT, want *bn256.GT) error {
		if gotB, wantB := got.Marshal(), want.Marshal(); !bytes.Equal(gotB, wantB) {
			return fmt.Errorf("Got %v want %v", got, want)
		}
		return nil
	}
	onetest := func(k1, k2 *big.Int) error {
		var (
			got  = Pair(new(G1).ScalarBaseMult(k1), new(G2).ScalarBaseMult(k2))
			want = bn256.Pair(new(bn256.G1).ScalarBaseMult(k1), new(bn256.G2).ScalarBaseMult(k2))
			gotB = got.Marshal()
		)
		if g, w := got.String(), want.String(); g != w {
			return fmt.Errorf("(%v, %v): String: Got %q, want %q", k1, k2, g, w)
		}
		if err := cmp(got, want); err != nil {
			return fmt.Errorf("(%v, %v): Pair: %v", k1, k2, err)
		}
		if err := cmp(new(GT).ScalarMult(got, k1), new(bn256.GT).ScalarMult(want, k1)); err != nil {
			return fmt.Errorf("(%v, %v): ScalarMult: %v", k1, k2, err)
		}
		if err := cmp(
			new(GT).Add(new(GT).ScalarMult(got, k1), new(GT).ScalarMult(got, k2)),
			new(bn256.GT).Add(new(bn256.GT).ScalarMult(want, k1), new(bn256.GT).ScalarMult(want, k2)),
		); err != nil {
			return fmt.Errorf("(%v, %v): Add: %v", k1, k2, err)
		}
		if err := cmp(new(GT).Neg(got), new(bn256.GT).Neg(want)); err != nil {
			return fmt.Errorf("(%v, %v): Neg: %v", k1, k2, err)
		}
		// Unmarshal and Marshal again.
		unmarshaled, ok := new(GT).Unmarshal(gotB)
		if !ok {
			return fmt.Errorf("(%v, %v): Unmarshal failed", k1, k2)
		}
		again := unmarshaled.Marshal()
		if !bytes.Equal(gotB, again) {
			return fmt.Errorf("(%v, %v): Umarshal+Marshal: Got %v, want %v", k1, k2, again, gotB)
		}
		return nil
	}
	big0, big1 := big.NewInt(0), big.NewInt(1)
	for _, test := range [][2]*big.Int{
		{big0, big0},
		{big0, big1},
		{big1, big0},
		{big1, big1},
	} {
		if err := onetest(test[0], test[1]); err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < 25; i++ {
		k1, err := rand.Int(rand.Reader, p)
		if err != nil {
			t.Fatal(err)
		}
		k2, err := rand.Int(rand.Reader, p)
		if err != nil {
			t.Fatal(err)
		}
		if err := onetest(k1, k2); err != nil {
			t.Errorf("%v (random test #%d)", err, i)
		}
	}
}

func TestBadUnmarshal(t *testing.T) {
	var (
		k  = big.NewInt(10) // Anything random
		g1 = new(G1).ScalarBaseMult(k)
		g2 = new(G2).ScalarBaseMult(k)
		gt = Pair(g1, g2)
		b1 = g1.Marshal()
		b2 = g2.Marshal()
		bt = gt.Marshal()
	)
	// nil, empty, one byte less, one byte more
	for _, test := range [][]byte{
		nil,
		make([]byte, 0),
		b1[0 : len(b1)-1],
		append(b1, 0),
		b2[0 : len(b2)-1],
		append(b2, 0),
		bt[0 : len(bt)-1],
		append(bt, 0),
	} {
		if _, ok := g1.Unmarshal(test); ok {
			t.Errorf("G1.Unmarshal succeeded on a %d byte slice", len(test))
		}
		if _, ok := g2.Unmarshal(test); ok {
			t.Errorf("G2.Unmarshal succeeded on a %d byte slice", len(test))
		}
		if _, ok := gt.Unmarshal(test); ok {
			t.Errorf("GT.Unmarshal succeeded on a %d byte slice", len(test))
		}
	}
}

func TestOrder(t *testing.T) {
	if !bytes.Equal(Order.Bytes(), bn256.Order.Bytes()) {
		t.Errorf("Got %v, want %v", Order, bn256.Order)
	}
}

var benchmarkA, benchmarkB *big.Int

func init() {
	benchmarkA, _ = rand.Int(rand.Reader, bn256.Order)
	benchmarkB, _ = rand.Int(rand.Reader, bn256.Order)
}

func BenchmarkPairGo(b *testing.B) {
	pa := new(bn256.G1).ScalarBaseMult(benchmarkA)
	qb := new(bn256.G2).ScalarBaseMult(benchmarkB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bn256.Pair(pa, qb)
	}
}

func BenchmarkPairCGO(b *testing.B) {
	pa := new(G1).ScalarBaseMult(benchmarkA)
	qb := new(G2).ScalarBaseMult(benchmarkB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pair(pa, qb)
	}
}
