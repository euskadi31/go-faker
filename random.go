// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "math/rand/v2"

// pick returns a random element from values. Returns the zero value of T when
// values is empty so callers do not have to special-case that path.
func pick[T any](rng *rand.Rand, values []T) T {
	if len(values) == 0 {
		var zero T

		return zero
	}

	return values[rng.IntN(len(values))]
}

// fillBytes writes len(b) pseudo-random bytes into b using rng.
func fillBytes(rng *rand.Rand, b []byte) {
	for i := range b {
		b[i] = byte(rng.IntN(256)) //nolint:gosec // value bounded to [0,255].
	}
}

// shuffle performs an in-place Fisher-Yates shuffle using rng.
func shuffle[T any](rng *rand.Rand, values []T) {
	for i := len(values) - 1; i > 0; i-- {
		j := rng.IntN(i + 1)
		values[i], values[j] = values[j], values[i]
	}
}

// randomInt64Inclusive returns a pseudo-random int64 in the inclusive range
// [lo, hi] using rng. Callers must ensure lo <= hi.
func randomInt64Inclusive(rng *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}

	return lo + rng.Int64N(hi-lo+1) // hi-lo+1 >= 2, so Int64N never panics.
}
