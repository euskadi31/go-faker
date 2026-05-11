// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"strings"
)

const domainAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

// generateFakeDomain builds a syntactically plausible but non-real domain. The
// label length is in [3, 10]; characters come from [a-z0-9-] with a '-' never
// at the start or end of the label.
func generateFakeDomain(rng *rand.Rand) string {
	labelLen := 3 + rng.IntN(8)

	var b strings.Builder

	b.Grow(labelLen)

	b.WriteByte(domainAlphabet[rng.IntN(len(domainAlphabet))])

	for i := 1; i < labelLen-1; i++ {
		// Inner characters may be '-' with low probability.
		if rng.IntN(8) == 0 {
			b.WriteByte('-')

			continue
		}

		b.WriteByte(domainAlphabet[rng.IntN(len(domainAlphabet))])
	}

	if labelLen > 1 {
		b.WriteByte(domainAlphabet[rng.IntN(len(domainAlphabet))])
	}

	return b.String() + "." + pickTLD(rng)
}
