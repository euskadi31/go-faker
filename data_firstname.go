// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "math/rand"

func getFirstName() string {
	firstNames := []string{}

	firstNames = append(firstNames, getManFirstName())
	firstNames = append(firstNames, getWomanFirstName())

	//nolint:gosec
	i := rand.Intn(len(firstNames) - 1)

	return firstNames[i]
}
