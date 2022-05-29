// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesGenerator(t *testing.T) {
	g := NewBytesGenerator(WithBytesLength(8))

	b := g.Generate()

	assert.Len(t, b, 8)
}

func BenchmarkBytesGenerator(b *testing.B) {
	g := NewBytesGenerator(WithBytesLength(8))

	for i := 0; i < b.N; i++ {
		g.Generate()
	}
}
