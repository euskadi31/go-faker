// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimezone(t *testing.T) {
	tz := Timezone()

	assert.NotEmpty(t, tz)
	assert.True(t, slices.Contains(dataTZ, tz), "%q is not a known time zone", tz)
}

func TestFakerTimezone(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	tz := f.Timezone()

	assert.NotEmpty(t, tz)
	assert.True(t, slices.Contains(dataTZ, tz), "%q is not a known time zone", tz)
}

func TestTimezoneDeterministic(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(7, 11))))
	b := New(WithRand(rand.New(rand.NewPCG(7, 11))))

	assert.Equal(t, a.Timezone(), b.Timezone())
}

func BenchmarkTimezone(b *testing.B) {
	for b.Loop() {
		Timezone()
	}
}
