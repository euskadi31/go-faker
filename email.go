// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"strings"
	"unicode"
)

// emailDomainMode controls how the email domain is chosen.
type emailDomainMode uint8

const (
	emailDomainExample emailDomainMode = iota
	emailDomainReal
	emailDomainFake
	emailDomainFixed
)

const exampleEmailDomain = "example.com"

// maxEmailSuffixDigits caps the random suffix so the local part stays well
// within the 64-octet limit of RFC 5321.
const maxEmailSuffixDigits = 12

// EmailOptions holds the settings used when generating an email address.
type EmailOptions struct {
	Domain string
	Prefix string
	// SuffixDigits is the number of random digits appended to the local part.
	// Zero (the default) means no suffix.
	SuffixDigits int
	mode         emailDomainMode
}

// EmailOption mutates EmailOptions.
type EmailOption func(*EmailOptions)

// WithEmailDomain forces the email domain to the given value.
func WithEmailDomain(domain string) EmailOption {
	return func(o *EmailOptions) {
		o.Domain = domain
		o.mode = emailDomainFixed
	}
}

// WithEmailPrefix prepends prefix to the local part of the generated address.
func WithEmailPrefix(prefix string) EmailOption {
	return func(o *EmailOptions) {
		o.Prefix = prefix
	}
}

// WithEmailSuffixDigits appends n random digits to the local part, as in
// jean.dupond8317@example.com. It multiplies the number of distinct addresses
// by 10^n, which matters for test suites generating many values from the same
// Faker. Values <= 0 disable the suffix; n is capped at 12.
func WithEmailSuffixDigits(n int) EmailOption {
	return func(o *EmailOptions) {
		if n > maxEmailSuffixDigits {
			n = maxEmailSuffixDigits
		}

		if n < 0 {
			n = 0
		}

		o.SuffixDigits = n
	}
}

// WithRealEmailDomain picks the domain from a list of well-known providers
// (gmail.com, yahoo.fr, …). Useful when test data must look realistic.
func WithRealEmailDomain() EmailOption {
	return func(o *EmailOptions) {
		o.mode = emailDomainReal
	}
}

// WithFakeEmailDomain generates a syntactically plausible but non-real domain.
func WithFakeEmailDomain() EmailOption {
	return func(o *EmailOptions) {
		o.mode = emailDomainFake
	}
}

// Email returns a generated email address. The default domain is example.com
// (RFC 2606 reserved) so generated values cannot collide with real accounts.
func (f *Faker) Email(opts ...EmailOption) string {
	o := EmailOptions{mode: emailDomainExample}

	for _, opt := range opts {
		opt(&o)
	}

	var domain string

	switch o.mode {
	case emailDomainFixed:
		domain = o.Domain
	case emailDomainReal:
		domain = pickRealEmailDomain(f.rand)
	case emailDomainFake:
		domain = generateFakeDomain(f.rand)
	case emailDomainExample:
		fallthrough
	default:
		domain = exampleEmailDomain
	}

	if domain == "" {
		domain = exampleEmailDomain
	}

	local := f.emailLocalPart()

	if o.Prefix != "" {
		local = strings.ToLower(o.Prefix) + local
	}

	if o.SuffixDigits > 0 {
		local += randomDigits(f.rand, o.SuffixDigits)
	}

	return local + "@" + domain
}

// Local-part shapes produced by emailLocalPart.
const (
	emailFirstDotLast uint8 = iota
	emailFirstUnderLast
	emailFirstLast
	emailLastDotFirst
	emailFirstDotInitial
	emailInitialLast
)

// emailLocalPatterns weights the local-part shapes by repetition. The two
// initial-based shapes keep a single letter of one name, so they span far fewer
// values (26 x len(dataLastNames)) than the full-name shapes; drawing them less
// often keeps the overall generated space large.
var emailLocalPatterns = []uint8{
	emailFirstDotLast, emailFirstDotLast, emailFirstDotLast, emailFirstDotLast,
	emailFirstUnderLast, emailFirstUnderLast, emailFirstUnderLast,
	emailFirstLast, emailFirstLast, emailFirstLast,
	emailLastDotFirst, emailLastDotFirst, emailLastDotFirst,
	emailInitialLast, emailInitialLast,
	emailFirstDotInitial,
}

// Modes accepted by sanitizeInto.
const (
	wholeName   = false
	initialOnly = true
)

// emailLocalPart builds the part before the '@'. Names are sanitized straight
// into the result to avoid intermediate allocations.
func (f *Faker) emailLocalPart() string {
	firstname := pickFirstName(f.rand)
	lastname := pickLastName(f.rand)

	var b strings.Builder

	b.Grow(len(firstname) + len(lastname) + 1)

	switch pick(f.rand, emailLocalPatterns) {
	case emailFirstUnderLast:
		sanitizeInto(&b, firstname, wholeName)
		b.WriteByte('_')
		sanitizeInto(&b, lastname, wholeName)
	case emailFirstLast:
		sanitizeInto(&b, firstname, wholeName)
		sanitizeInto(&b, lastname, wholeName)
	case emailLastDotFirst:
		sanitizeInto(&b, lastname, wholeName)
		b.WriteByte('.')
		sanitizeInto(&b, firstname, wholeName)
	case emailFirstDotInitial:
		sanitizeInto(&b, firstname, wholeName)
		b.WriteByte('.')
		sanitizeInto(&b, lastname, initialOnly)
	case emailInitialLast:
		sanitizeInto(&b, firstname, initialOnly)
		sanitizeInto(&b, lastname, wholeName)
	case emailFirstDotLast:
		fallthrough
	default:
		sanitizeInto(&b, firstname, wholeName)
		b.WriteByte('.')
		sanitizeInto(&b, lastname, wholeName)
	}

	return b.String()
}

// asciiFolding maps the non-ASCII runes present in the name data sets to their
// ASCII equivalent. Runes absent from this table and outside [a-z0-9] are
// dropped by sanitizeInto.
var asciiFolding = map[rune]string{
	'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a",
	'ç': "c",
	'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
	'ì': "i", 'í': "i", 'î': "i", 'ï': "i",
	'ñ': "n",
	'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o", 'ø': "o",
	'ù': "u", 'ú': "u", 'û': "u", 'ü': "u",
	'ý': "y", 'ÿ': "y",
	'æ': "ae", 'œ': "oe", 'ß': "ss",
}

// sanitizeInto appends s to b lower-cased, with accented characters folded to
// ASCII and everything outside [a-z0-9] dropped. Names such as "Muñoz" or
// "O’brien" (typographic apostrophe) would otherwise produce syntactically
// invalid addresses. When initial is true, only the first resulting character
// is written.
func sanitizeInto(b *strings.Builder, s string, initial bool) {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}

		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)

			if initial {
				return
			}

			continue
		}

		// Non-ASCII: fold to its ASCII equivalent. Unmapped characters are
		// dropped and do not count as the initial.
		folded, ok := asciiFolding[unicode.ToLower(r)]
		if !ok {
			continue
		}

		if initial {
			b.WriteString(folded[:1])

			return
		}

		b.WriteString(folded)
	}
}

// Email returns an email address generated by the default faker.
func Email(opts ...EmailOption) string {
	return defaultFaker().Email(opts...)
}

// EmailFlags is retained for backward compatibility with the pre-Faker API.
//
// Deprecated: use the option-based API on Faker.Email.
type EmailFlags uint16

// Deprecated email flag values.
const (
	// Deprecated: use WithRealEmailDomain().
	EmailRealDomain EmailFlags = 1 << iota
	// Deprecated: use WithEmailDomain(domain).
	EmailFixedDomain EmailFlags = 2 << iota
)

// EmailGenerator is the legacy generator type.
//
// Deprecated: use Faker.Email instead.
type EmailGenerator struct {
	domain string
	Flags  EmailFlags
}

// NewEmailGenerator constructs a legacy EmailGenerator.
//
// Deprecated: use New().Email(...) or the package-level Email helper.
func NewEmailGenerator(opts ...EmailOption) *EmailGenerator {
	g := &EmailGenerator{Flags: EmailRealDomain}

	tmp := EmailOptions{}

	for _, opt := range opts {
		opt(&tmp)
	}

	if tmp.mode == emailDomainFixed {
		g.domain = tmp.Domain
		g.Flags = EmailFixedDomain
	}

	return g
}

// Generate returns an email address following the legacy flag semantics.
//
// Deprecated: use Faker.Email.
func (g *EmailGenerator) Generate() string {
	switch {
	case g.Flags&EmailFixedDomain != 0 && g.domain != "":
		return defaultFaker().Email(WithEmailDomain(g.domain))
	case g.Flags&EmailRealDomain != 0:
		return defaultFaker().Email(WithRealEmailDomain())
	default:
		return defaultFaker().Email(WithFakeEmailDomain())
	}
}
