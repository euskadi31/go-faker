// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"regexp"
	"strings"
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

func TestFakerMACAddressDefault(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	mac := f.MACAddress()

	assert.Regexp(t, `^[A-F0-9]{2}(:[A-F0-9]{2}){5}$`, mac)
}

func TestFakerMACAddressLowercase(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	mac := f.MACAddress(WithMACLowercase())

	assert.Regexp(t, `^[a-f0-9]{2}(:[a-f0-9]{2}){5}$`, mac)
}

func TestFakerMACAddressSeparator(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	mac := f.MACAddress(WithMACSeparator("-"))

	assert.Regexp(t, `^[A-F0-9]{2}(-[A-F0-9]{2}){5}$`, mac)
}

func TestFakerMACAddressPrefix(t *testing.T) {
	tests := []struct {
		prefix string
		want   string
	}{
		{"FF", "FF:"},
		{"FF:FF", "FF:FF:"},
		{"ff-ff", "FF:FF:"},
	}

	for _, tc := range tests {
		t.Run(tc.prefix, func(t *testing.T) {
			f := New(WithRand(rand.New(rand.NewPCG(1, 2))))
			mac := f.MACAddress(WithMACPrefix(tc.prefix))
			assert.True(t, strings.HasPrefix(mac, tc.want), "got %q", mac)
			assert.Len(t, mac, 17)
		})
	}
}

func TestFakerMACAddressInvalidPrefixIgnored(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	mac := f.MACAddress(WithMACPrefix("not-hex"))

	assert.Regexp(t, `^[A-F0-9]{2}(:[A-F0-9]{2}){5}$`, mac)
}

func TestMACAddressDeterministic(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(7, 11))))
	b := New(WithRand(rand.New(rand.NewPCG(7, 11))))

	assert.Equal(t, a.MACAddress(), b.MACAddress())
}

func TestDeprecatedMacAddressAlias(t *testing.T) {
	mac := MacAddress()

	assert.Len(t, mac, 17)
}

func BenchmarkMacAddress(b *testing.B) {
	for b.Loop() {
		MacAddress()
	}
}

func BenchmarkMacAddressGenerator(b *testing.B) {
	g := NewMacAddressGenerator()
	g.Prefix = "FF"

	for b.Loop() {
		g.Generate()
	}
}
