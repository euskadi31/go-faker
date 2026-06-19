// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSIRENFormatAndLuhn(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for range 1000 {
		siren := f.SIREN()
		assert.Regexp(t, `^\d{9}$`, siren)
		assert.True(t, isLuhnValid(siren), "SIREN %q is not Luhn-valid", siren)
	}
}

func TestSIRETFormatAndLuhn(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for range 1000 {
		siret := f.SIRET()
		assert.Regexp(t, `^\d{14}$`, siret)
		assert.True(t, isLuhnValid(siret[:9]), "SIRET %q SIREN prefix not Luhn-valid", siret)
		assert.True(t, isLuhnValid(siret), "SIRET %q not Luhn-valid", siret)
	}
}

func TestEINFormatAndPrefix(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	allowed := make(map[string]struct{}, len(einPrefixes))

	for _, p := range einPrefixes {
		allowed[p] = struct{}{}
	}

	for range 200 {
		ein := f.EIN()
		assert.Regexp(t, `^\d{2}-\d{7}$`, ein)

		_, ok := allowed[ein[:2]]
		assert.True(t, ok, "EIN %q prefix not in published IRS list", ein)
	}
}

func TestDUNSFormat(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for range 50 {
		assert.Regexp(t, `^\d{9}$`, f.DUNS())
	}
}

func TestLEIFormatAndISO7064(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	pattern := regexp.MustCompile(`^[A-Z0-9]{4}00[A-Z0-9]{12}\d{2}$`)

	for range 200 {
		lei := f.LEI()
		assert.Regexp(t, pattern, lei)
		assert.True(t, isISO7064Mod97_10Valid(lei), "LEI %q failed ISO 7064 mod 97-10", lei)
	}
}

func TestVATNumberFR(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	pattern := regexp.MustCompile(`^FR\d{11}$`)

	for range 200 {
		vat := f.VATNumber(WithVATCountry("FR"))
		assert.Regexp(t, pattern, vat)

		siren := vat[4:]
		assert.True(t, isLuhnValid(siren), "FR VAT %q SIREN not Luhn-valid", vat)

		check, _ := strconv.Atoi(vat[2:4])
		n, _ := strconv.ParseInt(siren, 10, 64)
		assert.Equal(t, (12+3*(int(n)%97))%97, check, "FR VAT %q check mismatch", vat)
	}
}

func TestVATNumberBELuhnIT(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(3, 4))))

	for range 200 {
		be := f.VATNumber(WithVATCountry("BE"))
		assert.Regexp(t, `^BE0\d{9}$`, be)

		n, _ := strconv.ParseInt(be[2:10], 10, 64)
		check, _ := strconv.Atoi(be[10:])
		assert.Equal(t, 97-int(n%97), check, "BE VAT %q check mismatch", be)

		it := f.VATNumber(WithVATCountry("IT"))
		assert.Regexp(t, `^IT\d{11}$`, it)
		assert.True(t, isLuhnValid(it[2:]), "IT VAT %q not Luhn-valid", it)
	}
}

func TestVATNumberSKDivisibleBy11(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(5, 6))))

	for range 200 {
		sk := f.VATNumber(WithVATCountry("SK"))
		assert.Regexp(t, `^SK\d{10}$`, sk)

		n, err := strconv.ParseInt(sk[2:], 10, 64)
		require.NoError(t, err)
		assert.Equal(t, int64(0), n%11, "SK VAT %q not divisible by 11", sk)
	}
}

func TestVATNumberPTMod11(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(7, 8))))

	for range 200 {
		pt := f.VATNumber(WithVATCountry("PT"))
		assert.Regexp(t, `^PT\d{9}$`, pt)

		body := pt[2:10]
		sum := 0

		for i, c := range body {
			sum += int(c-'0') * (9 - i)
		}

		expected := 11 - sum%11
		if expected >= 10 {
			expected = 0
		}

		got, _ := strconv.Atoi(pt[10:])
		assert.Equal(t, expected, got, "PT VAT %q check mismatch", pt)
	}
}

func TestVATNumberAllEUCountriesShapes(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(9, 10))))
	patterns := map[string]*regexp.Regexp{
		"AT": regexp.MustCompile(`^ATU\d{8}$`),
		"BE": regexp.MustCompile(`^BE0\d{9}$`),
		"BG": regexp.MustCompile(`^BG\d{9,10}$`),
		"CY": regexp.MustCompile(`^CY\d{8}[A-Z]$`),
		"CZ": regexp.MustCompile(`^CZ\d{8,10}$`),
		"DE": regexp.MustCompile(`^DE\d{9}$`),
		"DK": regexp.MustCompile(`^DK\d{8}$`),
		"EE": regexp.MustCompile(`^EE\d{9}$`),
		"EL": regexp.MustCompile(`^EL\d{9}$`),
		"ES": regexp.MustCompile(`^ES[A-Z]\d{7}[A-Z]$`),
		"FI": regexp.MustCompile(`^FI\d{8}$`),
		"FR": regexp.MustCompile(`^FR\d{11}$`),
		"HR": regexp.MustCompile(`^HR\d{11}$`),
		"HU": regexp.MustCompile(`^HU\d{8}$`),
		"IE": regexp.MustCompile(`^IE\d{7}[A-Z]{1,2}$`),
		"IT": regexp.MustCompile(`^IT\d{11}$`),
		"LT": regexp.MustCompile(`^LT(\d{9}|\d{12})$`),
		"LU": regexp.MustCompile(`^LU\d{8}$`),
		"LV": regexp.MustCompile(`^LV\d{11}$`),
		"MT": regexp.MustCompile(`^MT\d{8}$`),
		"NL": regexp.MustCompile(`^NL\d{9}B\d{2}$`),
		"PL": regexp.MustCompile(`^PL\d{10}$`),
		"PT": regexp.MustCompile(`^PT\d{9}$`),
		"RO": regexp.MustCompile(`^RO\d{2,10}$`),
		"SE": regexp.MustCompile(`^SE\d{10}01$`),
		"SI": regexp.MustCompile(`^SI\d{8}$`),
		"SK": regexp.MustCompile(`^SK\d{10}$`),
	}

	require.Len(t, patterns, len(euCountries))

	for _, country := range euCountries {
		vat := f.VATNumber(WithVATCountry(country))
		assert.Regexp(t, patterns[country], vat, "country %s", country)
	}
}

func TestVATNumberRandomCountry(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(11, 12))))

	for range 100 {
		vat := f.VATNumber()
		require.GreaterOrEqual(t, len(vat), 4)
		assert.Contains(t, euCountries, vat[:2], "country %q not in EU set", vat[:2])
	}
}

func TestVATNumberUnknownCountryFallsBack(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(13, 14))))

	vat := f.VATNumber(WithVATCountry("ZZ"))
	require.GreaterOrEqual(t, len(vat), 4)
	assert.Contains(t, euCountries, vat[:2])
}

func TestVATNumberCountryIsCaseInsensitive(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(15, 16))))

	assert.True(t, strings.HasPrefix(f.VATNumber(WithVATCountry("fr")), "FR"))
}

func TestEUIDFormatAndCheck(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(17, 18))))
	pattern := regexp.MustCompile(`^[A-Z]{2}\.[A-Z0-9]{3,5}\.\d{6,12}\.\d{2}$`)

	for range 200 {
		euid := f.EUID()
		assert.Regexp(t, pattern, euid)

		parts := strings.Split(euid, ".")
		require.Len(t, parts, 4)
		payload := parts[0] + parts[1] + parts[2] + parts[3]
		assert.True(t, isISO7064Mod97_10Valid(payload), "EUID %q failed ISO 7064 mod 97-10", euid)
	}
}

func TestEUIDWithCountry(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(19, 20))))

	assert.True(t, strings.HasPrefix(f.EUID(WithEUIDCountry("FR")), "FR."))
	assert.True(t, strings.HasPrefix(f.EUID(WithEUIDCountry("de")), "DE."))
}

func TestEUIDUnknownCountryFallsBack(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(21, 22))))

	euid := f.EUID(WithEUIDCountry("ZZ"))
	assert.Contains(t, euCountries, euid[:2])
}

func TestEnterpriseIDReproducible(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))
	b := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	for range 25 {
		assert.Equal(t, a.SIREN(), b.SIREN())
		assert.Equal(t, a.SIRET(), b.SIRET())
		assert.Equal(t, a.EIN(), b.EIN())
		assert.Equal(t, a.DUNS(), b.DUNS())
		assert.Equal(t, a.LEI(), b.LEI())
		assert.Equal(t, a.VATNumber(WithVATCountry("FR")), b.VATNumber(WithVATCountry("FR")))
		assert.Equal(t, a.EUID(WithEUIDCountry("FR")), b.EUID(WithEUIDCountry("FR")))
	}
}

func TestEnterpriseIDPackageLevelHelpers(t *testing.T) {
	assert.Regexp(t, `^\d{9}$`, SIREN())
	assert.Regexp(t, `^\d{14}$`, SIRET())
	assert.Regexp(t, `^\d{2}-\d{7}$`, EIN())
	assert.Regexp(t, `^\d{9}$`, DUNS())
	assert.Regexp(t, `^[A-Z0-9]{20}$`, LEI())

	vat := VATNumber()
	require.GreaterOrEqual(t, len(vat), 4)
	assert.Contains(t, euCountries, vat[:2])

	euid := EUID()
	assert.Contains(t, euCountries, euid[:2])
	assert.Equal(t, 3, strings.Count(euid, "."))
}

func TestLuhnCheckDigitKnownValues(t *testing.T) {
	cases := []struct {
		prefix string
		want   int
	}{
		{"73282932", 0},      // SIREN 732829320
		{"7328293200007", 4}, // SIRET 73282932000074
		{"7992739871", 3},    // classic Luhn example 79927398713
	}

	for _, c := range cases {
		assert.Equal(t, c.want, luhnCheckDigit(c.prefix), "prefix %s", c.prefix)
	}
}

func TestISO7064Mod97_10Validates(t *testing.T) {
	// Self-consistency: computing the check then validating must return true.
	cases := []string{
		"ABCDEF1234567890XYZA",
		"5493001KJTIIGC8Y1R12",
		"0000000000000000",
	}

	for _, payload := range cases {
		check := iso7064Mod97_10CheckDigits(payload)
		assert.True(t, isISO7064Mod97_10Valid(payload+check), "payload %s", payload)
	}
}

func BenchmarkSIREN(b *testing.B) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for b.Loop() {
		f.SIREN()
	}
}

func BenchmarkLEI(b *testing.B) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for b.Loop() {
		f.LEI()
	}
}

func BenchmarkVATNumber(b *testing.B) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	for b.Loop() {
		f.VATNumber()
	}
}

// --- helpers -----------------------------------------------------------------

func isLuhnValid(digits string) bool {
	sum := 0
	n := len(digits)

	for i := range n {
		d := int(digits[i] - '0')
		if (n-1-i)%2 == 1 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
	}

	return sum%10 == 0
}

func isISO7064Mod97_10Valid(s string) bool {
	rem := 0

	for i := range len(s) {
		c := s[i]

		switch {
		case c >= '0' && c <= '9':
			rem = (rem*10 + int(c-'0')) % 97
		case c >= 'A' && c <= 'Z':
			rem = (rem*100 + int(c-'A') + 10) % 97
		default:
			return false
		}
	}

	return rem == 1
}
