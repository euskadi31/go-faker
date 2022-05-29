// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "crypto/rand"

// BytesGenerator struct.
type BytesGenerator struct {
	length int
}

type BytesOption func(g *BytesGenerator)

func WithBytesLength(length int) BytesOption {
	return func(g *BytesGenerator) {
		g.length = length
	}
}

// NewBytesGenerator return BytesGenerator.
func NewBytesGenerator(opts ...BytesOption) *BytesGenerator {
	g := &BytesGenerator{}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Generate bytes.
func (g *BytesGenerator) Generate() []byte {
	b := make([]byte, g.length)

	if _, err := rand.Read(b); err != nil {
		return []byte{}
	}

	return b
}
