// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBytesGenerator(t *testing.T) {
	g := NewBytesGenerator(WithBytesLength(8))

	b := g.Generate()

	assert.Len(t, b, 8)
}

func TestFakerBytesDeterministic(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(7, 11))))
	b := New(WithRand(rand.New(rand.NewPCG(7, 11))))

	assert.Equal(t, a.Bytes(32), b.Bytes(32))
}

func TestFakerBytesLength(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	assert.Len(t, f.Bytes(16), 16)
	assert.Empty(t, f.Bytes(0))
	assert.Empty(t, f.Bytes(-1))
}

func TestCryptoBytes(t *testing.T) {
	b, err := CryptoBytes(16)
	require.NoError(t, err)
	assert.Len(t, b, 16)

	empty, err := CryptoBytes(0)
	require.NoError(t, err)
	assert.Empty(t, empty)
}

func TestPackageBytes(t *testing.T) {
	b := Bytes(8)

	assert.Len(t, b, 8)
}

func BenchmarkBytesGenerator(b *testing.B) {
	g := NewBytesGenerator(WithBytesLength(8))

	for b.Loop() {
		g.Generate()
	}
}
