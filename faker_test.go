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

func TestNewDefault(t *testing.T) {
	f := New()

	require.NotNil(t, f)
	require.NotNil(t, f.Rand())
}

func TestWithRandNilIsIgnored(t *testing.T) {
	assert.NotPanics(t, func() {
		f := New(WithRand(nil))
		require.NotNil(t, f.Rand())
	})
}

func TestDeterminismSameSeed(t *testing.T) {
	mk := func() *Faker { return New(WithRand(rand.New(rand.NewPCG(1, 2)))) }

	f1, f2 := mk(), mk()

	for i := range 50 {
		assert.Equal(t, f1.Email(), f2.Email(), "iter %d", i)
		assert.Equal(t, f1.MACAddress(), f2.MACAddress(), "iter %d", i)
		assert.Equal(t, f1.Person(), f2.Person(), "iter %d", i)
		assert.Equal(t, f1.Bytes(8), f2.Bytes(8), "iter %d", i)
	}
}

func TestDeterminismIndependentInstances(t *testing.T) {
	f1 := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	f2 := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	// Burn some randomness on f1 only.
	for range 5 {
		_ = f1.Email()
	}

	// f2 was constructed identically; advance it the same way.
	for range 5 {
		_ = f2.Email()
	}

	assert.Equal(t, f1.Email(), f2.Email())
}

func TestDifferentSeedsLikelyDifferent(t *testing.T) {
	f1 := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	f2 := New(WithRand(rand.New(rand.NewPCG(3, 4))))

	differ := 0

	for range 50 {
		if f1.Email() != f2.Email() {
			differ++
		}
	}

	assert.Greater(t, differ, 40, "different seeds should yield mostly different output")
}
