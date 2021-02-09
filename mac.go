// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"strings"

	"github.com/euskadi31/go-reggen"
)

const macAddressPattern = "^[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}$"

// MacAddressGenerator struct
type MacAddressGenerator struct {
	Prefix    string
	Separator string
}

// NewMacAddressGenerator return MacAddressGenerator
func NewMacAddressGenerator() *MacAddressGenerator {
	return &MacAddressGenerator{
		Separator: ":",
	}
}

// Generate mac address
func (g *MacAddressGenerator) Generate() string {
	v, _ := reggen.Generate(macAddressPattern)

	if g.Prefix != "" {
		parts := strings.Split(v, ":")

		prefix := strings.ToUpper(g.Prefix)

		prefixParts := strings.Split(prefix, ":")

		for i := 0; i < len(prefixParts); i++ {
			parts[i] = prefixParts[i]
		}

		return strings.Join(parts, g.Separator)
	}

	return strings.Replace(v, ":", g.Separator, 6)
}

// MacAddress return an generated mac address
func MacAddress() string {
	return NewMacAddressGenerator().Generate()
}
