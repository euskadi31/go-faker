// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import cryptorand "crypto/rand"

// Bytes returns length pseudo-random bytes drawn from the Faker's RNG. The
// result is deterministic when WithRand was used at construction time. If
// length is <= 0, an empty slice is returned.
func (f *Faker) Bytes(length int) []byte {
	if length <= 0 {
		return []byte{}
	}

	b := make([]byte, length)
	fillBytes(f.rand, b)

	return b
}

// Bytes returns length pseudo-random bytes from the default faker.
func Bytes(length int) []byte {
	return defaultFaker().Bytes(length)
}

// CryptoBytes returns length cryptographically secure random bytes. It returns
// an empty slice when length is <= 0.
func CryptoBytes(length int) ([]byte, error) {
	if length <= 0 {
		return []byte{}, nil
	}

	b := make([]byte, length)

	if _, err := cryptorand.Read(b); err != nil {
		return nil, err //nolint:wrapcheck
	}

	return b, nil
}

// BytesOption is the legacy option type for BytesGenerator.
//
// Deprecated: use Faker.Bytes or CryptoBytes.
type BytesOption func(*BytesGenerator)

// WithBytesLength sets the BytesGenerator output length.
//
// Deprecated: pass the length directly to Faker.Bytes or CryptoBytes.
func WithBytesLength(length int) BytesOption {
	return func(g *BytesGenerator) {
		g.length = length
	}
}

// BytesGenerator is the legacy generator.
//
// Deprecated: use Faker.Bytes (deterministic) or CryptoBytes (crypto-secure).
type BytesGenerator struct {
	length int
}

// NewBytesGenerator returns a legacy BytesGenerator. The generator uses
// crypto/rand under the hood to preserve the prior behavior.
//
// Deprecated: use Faker.Bytes or CryptoBytes.
func NewBytesGenerator(opts ...BytesOption) *BytesGenerator {
	g := &BytesGenerator{}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Generate returns cryptographically secure random bytes. On error the legacy
// API returned an empty slice; this wrapper keeps that behavior.
//
// Deprecated: use CryptoBytes which surfaces errors.
func (g *BytesGenerator) Generate() []byte {
	b, err := CryptoBytes(g.length)
	if err != nil {
		return []byte{}
	}

	return b
}
