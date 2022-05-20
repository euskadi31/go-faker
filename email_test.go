// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailGenerator(t *testing.T) {
	g := NewEmailGenerator()
	g.Flags = 0

	email := g.Generate()

	assert.True(t, regexp.MustCompile("^[a-z0-9\\-.]+@[a-z0-9\\-.].[a-z0-9.]{2,52}").MatchString(email), "Email %q is not valid", email)
}

func TestEmailGeneratorWithFixedDomain(t *testing.T) {
	g := NewEmailGenerator(WithEmailDomain("acme.tld"))

	email := g.Generate()

	assert.True(t, regexp.MustCompile(`^[a-z0-9\-.]+@acme\.tld`).MatchString(email), "Email %q is not valid", email)
}

func TestEmail(t *testing.T) {
	email := Email()

	assert.True(t, regexp.MustCompile("^[a-z0-9\\-.]+@[a-z0-9\\-.].[a-z0-9.]{2,52}").MatchString(email), "Email %q is not valid", email)
}

func BenchmarkEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Email()
	}
}

func BenchmarkEmailGenerator(b *testing.B) {
	g := NewEmailGenerator()
	g.Flags = 0

	for i := 0; i < b.N; i++ {
		g.Generate()
	}
}
