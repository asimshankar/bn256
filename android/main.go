// +build android

// An "app" that runs the parings benchmarks and dumps the results to the
// logger. Flashes red/black when running and turns green when done.
package main

import (
	"crypto/rand"
	"flag"
	"log"
	"math/big"
	"testing"
	"time"

	thislib "github.com/asimshankar/bn256"
	stdlib "golang.org/x/crypto/bn256"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/gl"
)

var benchmarkA, benchmarkB *big.Int

func BenchmarkPairGo(b *testing.B) {
	pa := new(stdlib.G1).ScalarBaseMult(benchmarkA)
	qb := new(stdlib.G2).ScalarBaseMult(benchmarkB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdlib.Pair(pa, qb)
	}
}

func BenchmarkPairCGO(b *testing.B) {
	pa := new(thislib.G1).ScalarBaseMult(benchmarkA)
	qb := new(thislib.G2).ScalarBaseMult(benchmarkB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		thislib.Pair(pa, qb)
	}
}

func sanityTest() bool {
	var (
		a, _          = rand.Int(rand.Reader, stdlib.Order)
		b, _          = rand.Int(rand.Reader, stdlib.Order)
		stdlibPa      = new(stdlib.G1).ScalarBaseMult(a)
		stdlibQb      = new(stdlib.G2).ScalarBaseMult(b)
		stdlibPaired  = stdlib.Pair(stdlibPa, stdlibQb)
		thislibPa     = new(thislib.G1).ScalarBaseMult(a)
		thislibQb     = new(thislib.G2).ScalarBaseMult(b)
		thislibPaired = thislib.Pair(thislibPa, thislibQb)
	)
	if got, want := stdlibPaired.String(), thislibPaired.String(); got != want {
		log.Printf("ERROR: Got %q, want %q", got, want)
		return false
	}
	return true
}

func runBenchmarks(done chan<- bool) {
	if !sanityTest() {
		done <- false
		return
	}
	log.Printf("Pairing sanity check passed")
	flag.Set("test.benchtime", "2s")
	benchmarkA, _ = rand.Int(rand.Reader, stdlib.Order)
	benchmarkB, _ = rand.Int(rand.Reader, stdlib.Order)
	std := testing.Benchmark(BenchmarkPairGo)
	log.Printf("BenchmarkPairGo:  %v", std)
	this := testing.Benchmark(BenchmarkPairCGO)
	log.Printf("BenchmarkPairCGO: %v", this)
	done <- true
}

func main() {
	done := make(chan bool)
	go runBenchmarks(done)
	app.Main(func(a app.App) {
		var glctx gl.Context
		var success bool
		ticks := time.Tick(time.Second / 2)
		black := false
		for {
			select {
			case success = <-done:
				done = nil
				a.Send(paint.Event{})
			case <-ticks:
				black = !black
				a.Send(paint.Event{})
			case e := <-a.Events():
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					glctx, _ = e.DrawContext.(gl.Context)
				case paint.Event:
					if glctx == nil {
						continue
					}
					// solid green:         success
					// solid red:           failure
					// flashing black/blue: working
					if done == nil && success {
						glctx.ClearColor(0, 1, 0, 1)
					} else if done == nil && !success {
						glctx.ClearColor(1, 0, 0, 1)
					} else if black {
						glctx.ClearColor(0, 0, 0, 1)
					} else {
						glctx.ClearColor(0, 0, 1, 1)
					}
					glctx.Clear(gl.COLOR_BUFFER_BIT)
					a.Publish()
				}
			}
		}
	})
}
