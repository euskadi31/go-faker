// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocaleDefaultFormat(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	assert.Regexp(t, `^[a-z]{2}[-_][A-Z]{2}$`, f.Locale())
}

func TestLocaleWithDashSeparator(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	assert.Regexp(t, `^[a-z]{2}-[A-Z]{2}$`, f.Locale(WithLocaleSeparator("-")))
}

func TestLocaleWithUnderscoreSeparator(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	assert.Regexp(t, `^[a-z]{2}_[A-Z]{2}$`, f.Locale(WithLocaleSeparator("_")))
}

func TestLocaleWithLanguageAndCountry(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	locale := f.Locale(
		WithLocaleLanguage("EN"),
		WithLocaleCountry("fr"),
		WithLocaleSeparator("-"),
	)

	assert.Equal(t, "en-FR", locale)
}

func TestLocaleAllowsMixedLanguageAndCountry(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	locale := f.Locale(
		WithLocaleLanguage("en"),
		WithLocaleCountry("FR"),
		WithLocaleSeparator("-"),
	)

	assert.Equal(t, "en-FR", locale)
}

func TestLocaleInvalidSeparatorFallback(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	locale := f.Locale(
		WithLocaleLanguage("fr"),
		WithLocaleCountry("FR"),
		WithLocaleSeparator("."),
	)

	assert.Equal(t, "fr-FR", locale)
}

func TestLocaleSameSeed(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))
	b := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	seqA := make([]string, 10)
	seqB := make([]string, 10)

	for i := range seqA {
		seqA[i] = a.Locale()
		seqB[i] = b.Locale()
	}

	assert.Equal(t, seqA, seqB)
}

func TestLocaleDifferentSeed(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	b := New(WithRand(rand.New(rand.NewPCG(3, 4))))

	seqA := make([]string, 10)
	seqB := make([]string, 10)

	for i := range seqA {
		seqA[i] = a.Locale()
		seqB[i] = b.Locale()
	}

	assert.NotEqual(t, seqA, seqB)
}

func TestLocalePackageLevelHelper(t *testing.T) {
	assert.Regexp(t, `^[a-z]{2}[-_][A-Z]{2}$`, Locale())
}

func TestLocaleDataIsNotEmpty(t *testing.T) {
	assert.NotEmpty(t, languageCodes)
	assert.NotEmpty(t, countryCodes)

	langRE := regexp.MustCompile(`^[a-z]{2}$`)
	for _, code := range languageCodes {
		assert.True(t, langRE.MatchString(code), "language code %q is not lowercase 2-letter", code)
	}

	countryRE := regexp.MustCompile(`^[A-Z]{2}$`)
	for _, code := range countryCodes {
		assert.True(t, countryRE.MatchString(code), "country code %q is not uppercase 2-letter", code)
	}
}

func BenchmarkLocale(b *testing.B) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for b.Loop() {
		f.Locale()
	}
}
