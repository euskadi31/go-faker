// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "math/rand"

type Address struct {
	Location *AddressLocation
}

type AddressLocation struct {
	Latitude  float64
	Longitude float64
}

// AddressGenerator struct.
type AddressGenerator struct {
}

// NewAddressGenerator return AddressGenerator.
func NewAddressGenerator() *AddressGenerator {
	return &AddressGenerator{}
}

// Generate address.
func (g *AddressGenerator) Generate() *Address {
	addr := &Address{
		Location: &AddressLocation{
			Latitude:  (rand.Float64() * 180) - 90,
			Longitude: (rand.Float64() * 360) - 180,
		},
	}

	return addr
}
