// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
)

const enterpriseAlnumPool = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// SIREN returns a 9-digit French SIREN with a valid Luhn check digit.
func (f *Faker) SIREN() string {
	return generateSIREN(f.rand)
}

// SIREN returns a generated SIREN using the default faker.
func SIREN() string {
	return defaultFaker().SIREN()
}

// SIRET returns a 14-digit French SIRET with a valid Luhn check digit. The
// first 9 digits form a valid SIREN; the next 4 digits are the random NIC body
// and the last digit is the SIRET check.
func (f *Faker) SIRET() string {
	siren := generateSIREN(f.rand)
	prefix := siren + randomDigits(f.rand, 4)

	return prefix + strconv.Itoa(luhnCheckDigit(prefix))
}

// SIRET returns a generated SIRET using the default faker.
func SIRET() string {
	return defaultFaker().SIRET()
}

// einPrefixes lists the two-digit IRS campus codes assigned to EIN numbers.
var einPrefixes = []string{
	"01", "02", "03", "04", "05", "06", "10", "11", "12", "13", "14", "15",
	"16", "20", "21", "22", "23", "24", "25", "26", "27", "30", "31", "32",
	"33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44",
	"45", "46", "47", "48", "50", "51", "52", "53", "54", "55", "56", "57",
	"58", "59", "60", "61", "62", "63", "64", "65", "66", "67", "68", "71",
	"72", "73", "74", "75", "76", "77", "80", "81", "82", "83", "84", "85",
	"86", "87", "88", "90", "91", "92", "93", "94", "95", "98", "99",
}

// EIN returns a US Employer Identification Number formatted as "NN-NNNNNNN".
// The two-digit prefix is drawn from the IRS published campus code list; the
// body is random with no published checksum.
func (f *Faker) EIN() string {
	return pick(f.rand, einPrefixes) + "-" + randomDigits(f.rand, 7)
}

// EIN returns a generated EIN using the default faker.
func EIN() string {
	return defaultFaker().EIN()
}

// DUNS returns a 9-digit DUNS number. The format has no published checksum,
// so this is format only.
func (f *Faker) DUNS() string {
	return randomDigits(f.rand, 9)
}

// DUNS returns a generated DUNS using the default faker.
func DUNS() string {
	return defaultFaker().DUNS()
}

// LEI returns a 20-character ISO 17442 Legal Entity Identifier. The layout is
// 4-character LOU prefix, two reserved zeroes, 12-character entity-specific
// code, and a 2-digit ISO 7064 mod 97-10 check.
func (f *Faker) LEI() string {
	payload := randomAlnum(f.rand, 4) + "00" + randomAlnum(f.rand, 12)

	return payload + iso7064Mod97_10CheckDigits(payload)
}

// LEI returns a generated LEI using the default faker.
func LEI() string {
	return defaultFaker().LEI()
}

// euCountries lists the 27 EU member-state codes used by VATNumber and EUID.
var euCountries = []string{
	"AT", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "EL", "ES", "FI", "FR",
	"HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO",
	"SE", "SI", "SK",
}

// VATNumberOptions controls how VAT numbers are generated.
type VATNumberOptions struct {
	Country string
}

// VATNumberOption mutates VATNumberOptions.
type VATNumberOption func(*VATNumberOptions)

// WithVATCountry forces the ISO-3166-1 alpha-2 country code. Unknown or empty
// values fall back to a random EU member state.
func WithVATCountry(country string) VATNumberOption {
	return func(o *VATNumberOptions) {
		o.Country = country
	}
}

// vatBuilders maps each EU country to its body builder. The country prefix is
// added by VATNumber; each builder returns only the part after the prefix.
var vatBuilders = map[string]func(*rand.Rand) string{
	"AT": vatAT, "BE": vatBE, "BG": vatBG, "CY": vatCY, "CZ": vatCZ,
	"DE": vatDE, "DK": vatDK, "EE": vatEE, "EL": vatEL, "ES": vatES,
	"FI": vatFI, "FR": vatFR, "HR": vatHR, "HU": vatHU, "IE": vatIE,
	"IT": vatIT, "LT": vatLT, "LU": vatLU, "LV": vatLV, "MT": vatMT,
	"NL": vatNL, "PL": vatPL, "PT": vatPT, "RO": vatRO, "SE": vatSE,
	"SI": vatSI, "SK": vatSK,
}

// VATNumber returns an EU VAT registration number. When no country is forced,
// a random EU member state is picked. Per-country algorithms implement the
// published checksum where one is publicly defined; otherwise the value is
// format-only (correct length and charset).
func (f *Faker) VATNumber(opts ...VATNumberOption) string {
	o := VATNumberOptions{}

	for _, opt := range opts {
		opt(&o)
	}

	country := strings.ToUpper(o.Country)

	builder, ok := vatBuilders[country]
	if !ok {
		country = pick(f.rand, euCountries)
		builder = vatBuilders[country]
	}

	return country + builder(f.rand)
}

// VATNumber returns a generated VAT number using the default faker.
func VATNumber(opts ...VATNumberOption) string {
	return defaultFaker().VATNumber(opts...)
}

// EUIDOptions controls how EUID values are generated.
type EUIDOptions struct {
	Country string
}

// EUIDOption mutates EUIDOptions.
type EUIDOption func(*EUIDOptions)

// WithEUIDCountry forces the EU member state used as the EUID country prefix.
// Unknown or empty values fall back to a random EU member state.
func WithEUIDCountry(country string) EUIDOption {
	return func(o *EUIDOptions) {
		o.Country = country
	}
}

// EUID returns a BRIS-style EU Unique Identifier per Regulation (EU) 2015/884
// of the form "<CC>.<RegisterID>.<EntityID>.<Check>". The check is an ISO 7064
// mod 97-10 over the unpunctuated payload.
func (f *Faker) EUID(opts ...EUIDOption) string {
	o := EUIDOptions{}

	for _, opt := range opts {
		opt(&o)
	}

	country := strings.ToUpper(o.Country)
	if !slices.Contains(euCountries, country) {
		country = pick(f.rand, euCountries)
	}

	register := randomAlnum(f.rand, 3+f.rand.IntN(3)) // 3-5 chars
	entity := randomDigits(f.rand, 6+f.rand.IntN(7))  // 6-12 digits
	check := iso7064Mod97_10CheckDigits(country + register + entity)

	return country + "." + register + "." + entity + "." + check
}

// EUID returns a generated EUID using the default faker.
func EUID(opts ...EUIDOption) string {
	return defaultFaker().EUID(opts...)
}

// --- per-country VAT builders ------------------------------------------------

// vatFR returns "<2-digit check><9-digit SIREN>". The check is
// (12 + 3 * (SIREN mod 97)) mod 97.
func vatFR(rng *rand.Rand) string {
	siren := generateSIREN(rng)

	n, _ := strconv.ParseInt(siren, 10, 64)
	check := (12 + 3*(int(n)%97)) % 97

	return fmt.Sprintf("%02d%s", check, siren)
}

// vatBE returns 10 digits: "0" + 7 random digits + 2 check digits.
// Check = 97 - (first 8 digits mod 97), formatted on 2 digits.
func vatBE(rng *rand.Rand) string {
	base := "0" + randomDigits(rng, 7)

	n, _ := strconv.ParseInt(base, 10, 64)
	check := 97 - int(n%97)

	return fmt.Sprintf("%s%02d", base, check)
}

// vatIT returns 11 digits where the last is the Luhn check digit.
func vatIT(rng *rand.Rand) string {
	body := randomDigits(rng, 10)

	return body + strconv.Itoa(luhnCheckDigit(body))
}

// vatSK returns 10 digits whose value is divisible by 11.
func vatSK(rng *rand.Rand) string {
	// Smallest multiplier whose 11-product has 10 digits: ceil(10^9 / 11).
	// Largest:                                            floor((10^10 - 1) / 11).
	const lo, hi = 90909091, 909090909

	n := randomInt64Inclusive(rng, lo, hi)

	return fmt.Sprintf("%010d", n*11)
}

// vatPT returns 9 digits where the last is a mod 11 weighted check
// (weights 9..2, check = (11 - sum%11) mod 10, mapped to 0 when >= 10).
func vatPT(rng *rand.Rand) string {
	body := randomDigits(rng, 8)

	sum := 0
	for i, c := range body {
		sum += int(c-'0') * (9 - i)
	}

	check := 11 - (sum % 11)
	if check >= 10 {
		check = 0
	}

	return body + strconv.Itoa(check)
}

// The builders below match the documented length and charset for each country
// but apply no checksum — none is publicly defined, or the algorithm is too
// peripheral to maintain. Each value is suitable for tests, not validators.

func vatAT(rng *rand.Rand) string { return "U" + randomDigits(rng, 8) }
func vatBG(rng *rand.Rand) string { return randomDigits(rng, 9+rng.IntN(2)) }
func vatCY(rng *rand.Rand) string { return randomDigits(rng, 8) + string(randomUpperLetter(rng)) }
func vatCZ(rng *rand.Rand) string { return randomDigits(rng, 8+rng.IntN(3)) }
func vatDE(rng *rand.Rand) string { return randomDigits(rng, 9) }
func vatDK(rng *rand.Rand) string { return randomDigits(rng, 8) }
func vatEE(rng *rand.Rand) string { return randomDigits(rng, 9) }
func vatEL(rng *rand.Rand) string { return randomDigits(rng, 9) }
func vatES(rng *rand.Rand) string {
	return string(randomUpperLetter(rng)) + randomDigits(rng, 7) + string(randomUpperLetter(rng))
}
func vatFI(rng *rand.Rand) string { return randomDigits(rng, 8) }
func vatHR(rng *rand.Rand) string { return randomDigits(rng, 11) }
func vatHU(rng *rand.Rand) string { return randomDigits(rng, 8) }
func vatIE(rng *rand.Rand) string {
	body := randomDigits(rng, 7) + string(randomUpperLetter(rng))
	if rng.IntN(2) == 0 {
		body += string(randomUpperLetter(rng))
	}

	return body
}

func vatLT(rng *rand.Rand) string {
	if rng.IntN(2) == 0 {
		return randomDigits(rng, 9)
	}

	return randomDigits(rng, 12)
}
func vatLU(rng *rand.Rand) string { return randomDigits(rng, 8) }
func vatLV(rng *rand.Rand) string { return randomDigits(rng, 11) }
func vatMT(rng *rand.Rand) string { return randomDigits(rng, 8) }
func vatNL(rng *rand.Rand) string { return randomDigits(rng, 9) + "B" + randomDigits(rng, 2) }
func vatPL(rng *rand.Rand) string { return randomDigits(rng, 10) }
func vatRO(rng *rand.Rand) string { return randomDigits(rng, 2+rng.IntN(9)) }
func vatSE(rng *rand.Rand) string { return randomDigits(rng, 10) + "01" }
func vatSI(rng *rand.Rand) string { return randomDigits(rng, 8) }

// --- helpers -----------------------------------------------------------------

// generateSIREN returns a 9-digit SIREN whose Luhn check is valid.
func generateSIREN(rng *rand.Rand) string {
	body := randomDigits(rng, 8)

	return body + strconv.Itoa(luhnCheckDigit(body))
}

// randomDigits returns a string of n random decimal digits.
func randomDigits(rng *rand.Rand, n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = byte('0' + rng.IntN(10)) //nolint:gosec // '0'+[0,9] always fits in byte.
	}

	return string(out)
}

// randomAlnum returns a string of n uppercase A-Z / 0-9 characters.
func randomAlnum(rng *rand.Rand, n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = enterpriseAlnumPool[rng.IntN(len(enterpriseAlnumPool))]
	}

	return string(out)
}

// randomUpperLetter returns a random ASCII uppercase letter.
func randomUpperLetter(rng *rand.Rand) byte {
	return byte('A' + rng.IntN(26)) //nolint:gosec // 'A'+[0,25] always fits in byte.
}

// luhnCheckDigit returns the digit that, appended to prefix, makes the full
// numeric string Luhn-valid. prefix must contain only ASCII digits.
func luhnCheckDigit(prefix string) int {
	sum := 0
	n := len(prefix)

	for i := range n {
		d := int(prefix[i] - '0')
		if (n-i)%2 == 1 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
	}

	return (10 - sum%10) % 10
}

// iso7064Mod97_10CheckDigits returns the two-digit ISO 7064 mod 97-10 check
// for payload. Letters in payload are converted to their numeric value (A=10,
// …, Z=35) before the modular reduction. The returned string makes
// payload + result valid (mod 97 == 1).
func iso7064Mod97_10CheckDigits(payload string) string {
	rem := 0

	for i := range len(payload) {
		c := payload[i]

		switch {
		case c >= '0' && c <= '9':
			rem = (rem*10 + int(c-'0')) % 97
		case c >= 'A' && c <= 'Z':
			rem = (rem*100 + int(c-'A') + 10) % 97
		}
	}

	rem = (rem * 100) % 97

	return fmt.Sprintf("%02d", 98-rem)
}
