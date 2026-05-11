// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand/v2"
	"sync"
)

// Faker holds the random source used by every generator method.
//
// A Faker is not safe for concurrent use by multiple goroutines. Callers that
// need deterministic generation across goroutines must use one Faker per
// goroutine, each with its own *rand.Rand.
type Faker struct {
	rand *rand.Rand
}

// Option configures a Faker at construction time.
type Option func(*Faker)

// New returns a Faker. Without options, it is seeded from crypto/rand so each
// call to New produces an independent, non-deterministic source.
func New(opts ...Option) *Faker {
	//nolint:gosec // pseudo-random by design; callers needing crypto use CryptoBytes.
	f := &Faker{
		rand: rand.New(rand.NewPCG(cryptoSeed(), cryptoSeed())),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

// WithRand injects the random source used by the Faker. Pass a deterministic
// source (e.g. rand.New(rand.NewPCG(seed1, seed2))) for reproducible output.
// A nil argument is ignored and the default source is kept.
func WithRand(rng *rand.Rand) Option {
	return func(f *Faker) {
		if rng != nil {
			f.rand = rng
		}
	}
}

// Rand returns the underlying random source. Useful when callers need to share
// the same RNG with their own generation code.
func (f *Faker) Rand() *rand.Rand {
	return f.rand
}

func cryptoSeed() uint64 {
	var b [8]byte

	if _, err := cryptorand.Read(b[:]); err != nil {
		return 0
	}

	return binary.LittleEndian.Uint64(b[:])
}

var (
	defaultFakerOnce sync.Once
	defaultFakerInst *Faker
)

// defaultFaker returns the package-level Faker used by helper functions.
func defaultFaker() *Faker {
	defaultFakerOnce.Do(func() {
		defaultFakerInst = New()
	})

	return defaultFakerInst
}
