// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMacAddressGenerator(t *testing.T) {
	tests := []struct {
		prefix string
		expect string
	}{
		{
			prefix: "FF",
			expect: "^FF:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}$",
		},
		{
			prefix: "FF:FF",
			expect: "^FF:FF:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}$",
		},
	}

	for _, test := range tests {
		g := NewMacAddressGenerator()
		g.Prefix = test.prefix

		mac := g.Generate()

		assert.Equal(t, 17, len(mac))
		assert.True(t, regexp.MustCompile(test.expect).MatchString(mac), "Mac %q not match %q", mac, test.expect)
	}
}

func TestMacAddress(t *testing.T) {
	mac := MacAddress()

	assert.Equal(t, 17, len(mac))
}

func BenchmarkMacAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MacAddress()
	}
}
