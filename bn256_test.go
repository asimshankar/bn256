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
	onetest := func(k *big.Int) error {
		var (
			got   = new(G1).ScalarBaseMult(k)
			gotS  = got.String()
			gotB  = got.Marshal()
			want  = new(bn256.G1).ScalarBaseMult(k)
			wantS = want.String()
			wantB = want.Marshal()
		)
		// TODO: Minor implementation difference causes String for
		// golang.org/x/crypto/bn256.G1 to return (1, -2) for k=1,
		// while (1, 65000549695646603732796438742359905742825358107623003571877145026864184071781)
		// for this package. The two are identical since
		// (-2 mod p) == 65000549695646603732796438742359905742825358107623003571877145026864184071781
		// So, ignore that difference.
		if k.Cmp(big.NewInt(1)) == 0 {
			wantS = "bn256.G1(1, 65000549695646603732796438742359905742825358107623003571877145026864184071781)"
		}
		if gotS != wantS {
			return fmt.Errorf("k=%v: String: Got %q, want %q", k, gotS, wantS)
		}
		if !bytes.Equal(gotB, wantB) {
			return fmt.Errorf("k=%v: Marshal: Got %v, want %v", k, gotB, wantB)
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
	onetest := func(k *big.Int) error {
		var (
			got   = new(G2).ScalarBaseMult(k)
			gotS  = got.String()
			gotB  = got.Marshal()
			want  = new(bn256.G2).ScalarBaseMult(k)
			wantS = want.String()
			wantB = want.Marshal()
		)
		if gotS != wantS {
			return fmt.Errorf("k=%v: String: Got %q, want %q", k, gotS, wantS)
		}
		if !bytes.Equal(gotB, wantB) {
			return fmt.Errorf("k=%v: Marshal: Got %v, want %v", k, gotB, wantB)
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
}

func TestBadUnmarshal(t *testing.T) {
	var (
		k  = big.NewInt(10) // Anything random
		g1 = new(G1).ScalarBaseMult(k)
		g2 = new(G2).ScalarBaseMult(k)
		b1 = g1.Marshal()
		b2 = g2.Marshal()
	)
	// nil, empty, one byte less, one byte more
	for _, test := range [][]byte{
		nil,
		make([]byte, 0),
		b1[0 : len(b1)-1],
		append(b1, 0),
		b2[0 : len(b2)-1],
		append(b2, 0),
	} {
		if _, ok := g1.Unmarshal(test); ok {
			t.Errorf("G1.Unmarshal succeeded on a %d byte slice", len(test))
		}
		if _, ok := g2.Unmarshal(test); ok {
			t.Errorf("G2.Unmarshal succeeded on a %d byte slice", len(test))
		}
	}
}

var (
	benchmarkK             *big.Int
	benchmarkA, benchmarkB *big.Int
)

func init() {
	// Randomly chose one
	var ok bool
	if benchmarkK, ok = new(big.Int).SetString("55957183647262293325367359498614325417154459764697977524189246266898748271344", 10); !ok {
		panic("failed to set value for benchmarkK")
	}
	benchmarkA, _ = rand.Int(rand.Reader, bn256.Order)
	benchmarkB, _ = rand.Int(rand.Reader, bn256.Order)
}

// Ultimately, the only benchmark that will matter is the one for Pair, but
// some other tests in the meantime.
func BenchmarkG1_ScalarBaseMult_Baseline(b *testing.B) {
	g := new(bn256.G1)
	for i := 0; i < b.N; i++ {
		g.ScalarBaseMult(benchmarkK)
	}
}
func BenchmarkG1_ScalarBaseMult(b *testing.B) {
	g := new(G1)
	for i := 0; i < b.N; i++ {
		g.ScalarBaseMult(benchmarkK)
	}
}
func BenchmarkG1_ScalarBaseMult_C(b *testing.B) {
	benchmarkG1ScalarBaseMult(b.N, benchmarkK)
}
func BenchmarkG2_ScalarBaseMult_Baseline(b *testing.B) {
	g := new(bn256.G2)
	for i := 0; i < b.N; i++ {
		g.ScalarBaseMult(benchmarkK)
	}
}
func BenchmarkG2_ScalarBaseMult(b *testing.B) {
	g := new(G2)
	for i := 0; i < b.N; i++ {
		g.ScalarBaseMult(benchmarkK)
	}
}
func BenchmarkG2_ScalarBaseMult_C(b *testing.B) {
	benchmarkG2ScalarBaseMult(b.N, benchmarkK)
}

func BenchmarkPair_Baseline(b *testing.B) {
	pa := new(bn256.G1).ScalarBaseMult(benchmarkA)
	qb := new(bn256.G2).ScalarBaseMult(benchmarkB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bn256.Pair(pa, qb)
	}
}
