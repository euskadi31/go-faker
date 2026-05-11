// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "math/rand/v2"

func pickFirstName(rng *rand.Rand) string {
	if rng.IntN(2) == 0 {
		return pickManFirstName(rng)
	}

	return pickWomanFirstName(rng)
}
