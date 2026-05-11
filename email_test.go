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

func TestEmailGenerator(t *testing.T) {
	g := NewEmailGenerator()
	g.Flags = 0

	email := g.Generate()

	assert.True(t, regexp.MustCompile(`^[a-z0-9.\-]+@[a-z0-9.\-]+\.[a-z0-9.\-]{2,}$`).MatchString(email), "Email %q is not valid", email)
}

func TestEmailGeneratorWithFixedDomain(t *testing.T) {
	g := NewEmailGenerator(WithEmailDomain("acme.tld"))

	email := g.Generate()

	assert.True(t, regexp.MustCompile(`^[a-z0-9\-.]+@acme\.tld`).MatchString(email), "Email %q is not valid", email)
}

func TestEmail(t *testing.T) {
	email := Email()

	assert.True(t, regexp.MustCompile(`^[a-z0-9.\-]+@[a-z0-9.\-]+\.[a-z0-9.\-]{2,}$`).MatchString(email), "Email %q is not valid", email)
}

func TestFakerEmailDefaultUsesExampleDomain(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	email := f.Email()

	assert.True(t, strings.HasSuffix(email, "@example.com"), "got %q", email)
	assert.Equal(t, strings.ToLower(email), email)
}

func TestFakerEmailWithFixedDomain(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	email := f.Email(WithEmailDomain("acme.tld"))

	assert.True(t, strings.HasSuffix(email, "@acme.tld"), "got %q", email)
}

func TestFakerEmailWithPrefix(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	email := f.Email(WithEmailPrefix("test-"), WithEmailDomain("acme.tld"))

	assert.True(t, strings.HasPrefix(email, "test-"), "got %q", email)
	assert.True(t, strings.HasSuffix(email, "@acme.tld"), "got %q", email)
}

func TestFakerEmailWithRealDomain(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	email := f.Email(WithRealEmailDomain())

	parts := strings.SplitN(email, "@", 2)

	if assert.Len(t, parts, 2) {
		assert.Contains(t, dataRealEmailDomains, parts[1])
	}
}

func TestFakerEmailWithFakeDomain(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	email := f.Email(WithFakeEmailDomain())

	assert.Regexp(t, `^[a-z0-9.\-]+@[a-z0-9\-]+\.[a-z0-9.\-]+$`, email)
}

func TestEmailDeterministic(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(7, 11))))
	b := New(WithRand(rand.New(rand.NewPCG(7, 11))))

	assert.Equal(t, a.Email(), b.Email())
}

func BenchmarkEmail(b *testing.B) {
	for b.Loop() {
		Email()
	}
}

func BenchmarkEmailGenerator(b *testing.B) {
	g := NewEmailGenerator()
	g.Flags = 0

	for b.Loop() {
		g.Generate()
	}
}
