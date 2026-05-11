// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPickEmptySliceReturnsZero(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	assert.Equal(t, "", pick(rng, []string{}))
	assert.Equal(t, 0, pick(rng, []int{}))
}

func TestPickCoversFirstAndLastElement(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	values := []int{1, 2, 3}
	seen := map[int]bool{}

	for range 200 {
		seen[pick(rng, values)] = true
	}

	for _, v := range values {
		assert.True(t, seen[v], "value %d never picked", v)
	}
}

func TestFillBytesDeterministic(t *testing.T) {
	a := make([]byte, 32)
	b := make([]byte, 32)

	fillBytes(rand.New(rand.NewPCG(42, 42)), a)
	fillBytes(rand.New(rand.NewPCG(42, 42)), b)

	assert.Equal(t, a, b)
}

func TestShuffleDeterministic(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	shuffle(rand.New(rand.NewPCG(7, 11)), a)
	shuffle(rand.New(rand.NewPCG(7, 11)), b)

	assert.Equal(t, a, b)
}
