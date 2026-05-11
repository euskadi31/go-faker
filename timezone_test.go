// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimezone(t *testing.T) {
	tz := Timezone()

	assert.NotEmpty(t, tz)
	assert.True(t, slices.Contains(dataTZ, tz), "%q is not a known time zone", tz)

	_, err := time.LoadLocation(tz)
	assert.NoError(t, err, "%q is not loadable via time.LoadLocation", tz)
}

func TestFakerTimezone(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	tz := f.Timezone()

	assert.NotEmpty(t, tz)
	assert.True(t, slices.Contains(dataTZ, tz), "%q is not a known time zone", tz)

	_, err := time.LoadLocation(tz)
	assert.NoError(t, err, "%q is not loadable via time.LoadLocation", tz)
}

func TestTimezoneDataLoadable(t *testing.T) {
	for _, tz := range dataTZ {
		_, err := time.LoadLocation(tz)
		assert.NoErrorf(t, err, "%q in dataTZ is not loadable via time.LoadLocation", tz)
	}
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
