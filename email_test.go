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

	assert.True(t, regexp.MustCompile(`^[a-z0-9._\-]+@[a-z0-9.\-]+\.[a-z0-9.\-]{2,}$`).MatchString(email), "Email %q is not valid", email)
}

func TestEmailGeneratorWithFixedDomain(t *testing.T) {
	g := NewEmailGenerator(WithEmailDomain("acme.tld"))

	email := g.Generate()

	assert.True(t, regexp.MustCompile(`^[a-z0-9\-._]+@acme\.tld`).MatchString(email), "Email %q is not valid", email)
}

func TestEmail(t *testing.T) {
	email := Email()

	assert.True(t, regexp.MustCompile(`^[a-z0-9._\-]+@[a-z0-9.\-]+\.[a-z0-9.\-]{2,}$`).MatchString(email), "Email %q is not valid", email)
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

	assert.Regexp(t, `^[a-z0-9._\-]+@[a-z0-9\-]+\.[a-z0-9.\-]+$`, email)
}

// sanitizeLocalPart returns the sanitized form of s. Production code writes
// straight into the result builder via sanitizeInto; this wrapper exists so the
// folding rules can be asserted in isolation.
func sanitizeLocalPart(s string) string {
	var b strings.Builder

	b.Grow(len(s))

	sanitizeInto(&b, s, wholeName)

	return b.String()
}

func TestSanitizeLocalPart(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
	}{
		{"Dupond", "dupond"},
		{"Muñoz", "munoz"},
		{"Peña", "pena"},
		{"Nuñez", "nunez"},
		{"Piñeda", "pineda"},
		{"Cortéz", "cortez"},
		{"Ávila", "avila"},
		{"Colón", "colon"},
		{"O’brien", "obrien"},
		{"O’connor", "oconnor"},
		{"O'neill", "oneill"},
		{"Le Blanc", "leblanc"},
		{"Straße", "strasse"},
		{"Æther", "aether"},
		{"", ""},
	} {
		assert.Equal(t, tc.want, sanitizeLocalPart(tc.in), "input %q", tc.in)
	}
}

// TestFakerEmailAlwaysValid guards against name data leaking non-ASCII
// characters into the local part: 13 of the surnames contain accents or a
// typographic apostrophe.
func TestFakerEmailAlwaysValid(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(3, 5))))
	re := regexp.MustCompile(`^[a-z0-9._\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	for range 100000 {
		email := f.Email()

		assert.Regexp(t, re, email)
	}
}

func TestFakerEmailEntropy(t *testing.T) {
	const runs = 50000

	f := New(WithRand(rand.New(rand.NewPCG(13, 17))))
	seen := make(map[string]struct{}, runs)

	for range runs {
		seen[f.Email()] = struct{}{}
	}

	assert.Greater(t, len(seen), runs*99/100, "too many duplicate addresses")
}

func TestFakerEmailEntropyWithSuffix(t *testing.T) {
	const runs = 200000

	f := New(WithRand(rand.New(rand.NewPCG(19, 23))))
	seen := make(map[string]struct{}, runs)

	for range runs {
		seen[f.Email(WithEmailSuffixDigits(4))] = struct{}{}
	}

	assert.Greater(t, len(seen), runs*99/100, "too many duplicate addresses")
}

func TestFakerEmailWithSuffixDigits(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for range 100 {
		email := f.Email(WithEmailSuffixDigits(4), WithEmailDomain("acme.tld"))

		local := strings.TrimSuffix(email, "@acme.tld")

		assert.Regexp(t, `^[a-z0-9._]+[0-9]{4}$`, local)
		assert.Regexp(t, `[a-z][0-9]{4}$`, local, "suffix must follow a letter")
	}
}

func TestFakerEmailSuffixDigitsBounds(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for _, n := range []int{0, -1, -100} {
		o := EmailOptions{}
		WithEmailSuffixDigits(n)(&o)

		assert.Equal(t, 0, o.SuffixDigits, "n = %d", n)
	}

	o := EmailOptions{}
	WithEmailSuffixDigits(1000)(&o)
	assert.Equal(t, maxEmailSuffixDigits, o.SuffixDigits)

	email := f.Email(WithEmailSuffixDigits(1000), WithEmailDomain("acme.tld"))
	local := strings.TrimSuffix(email, "@acme.tld")

	assert.Regexp(t, `[0-9]{12}$`, local)
}

func TestFakerEmailSuffixWithPrefix(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	email := f.Email(
		WithEmailPrefix("qa-"),
		WithEmailSuffixDigits(3),
		WithEmailDomain("acme.tld"),
	)

	assert.Regexp(t, `^qa-[a-z0-9._]+[0-9]{3}@acme\.tld$`, email)
}

func TestLastNameDataIsUsable(t *testing.T) {
	assert.GreaterOrEqual(t, len(dataLastNames), 5000)

	seen := make(map[string]struct{}, len(dataLastNames))

	for _, name := range dataLastNames {
		assert.NotEmpty(t, name)

		local := sanitizeLocalPart(name)

		assert.Regexp(t, `^[a-z0-9]+$`, local, "surname %q", name)

		_, dup := seen[local]
		assert.False(t, dup, "surname %q collides with an earlier entry", name)
		seen[local] = struct{}{}
	}
}

func TestFirstNameDataIsUsable(t *testing.T) {
	for _, names := range [][]string{dataManFirstNames, dataWomanFirstNames} {
		for _, name := range names {
			assert.NotEmpty(t, name)
			assert.Regexp(t, `^[a-z0-9]+$`, sanitizeLocalPart(name), "first name %q", name)
		}
	}
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
