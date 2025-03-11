// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressGenerator(t *testing.T) {

	{
		g := NewAddressGenerator()

		address := g.Generate()

		assert.True(t, (address.Location.Latitude > -90 && address.Location.Latitude < 90))
		assert.True(t, (address.Location.Longitude > -180 && address.Location.Longitude < 180))
	}

}

func BenchmarkAddressGenerator(b *testing.B) {
	g := NewAddressGenerator()

	for i := 0; i < b.N; i++ {
		g.Generate()
	}
}
