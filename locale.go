// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "strings"

// LocaleOptions controls how Locale values are generated.
type LocaleOptions struct {
	Language  string
	Country   string
	Separator string
}

// LocaleOption mutates LocaleOptions.
type LocaleOption func(*LocaleOptions)

// WithLocaleLanguage forces the language subtag. The value is normalized to
// lower case in the final locale.
func WithLocaleLanguage(language string) LocaleOption {
	return func(o *LocaleOptions) {
		o.Language = language
	}
}

// WithLocaleCountry forces the country subtag. The value is normalized to
// upper case in the final locale.
func WithLocaleCountry(country string) LocaleOption {
	return func(o *LocaleOptions) {
		o.Country = country
	}
}

// WithLocaleSeparator forces the subtag separator. Only "-" and "_" are
// accepted; any other value falls back to "-".
func WithLocaleSeparator(separator string) LocaleOption {
	return func(o *LocaleOptions) {
		o.Separator = separator
	}
}

// Locale returns a generated locale identifier such as "fr-FR", "en_US" or
// "en-FR". Language and country are picked independently so combinations that
// are uncommon in CLDR but valid on real devices (English language with French
// region, etc.) can occur.
func (f *Faker) Locale(opts ...LocaleOption) string {
	o := LocaleOptions{}

	for _, opt := range opts {
		opt(&o)
	}

	language := o.Language
	if language == "" {
		language = pick(f.rand, languageCodes)
	}

	country := o.Country
	if country == "" {
		country = pick(f.rand, countryCodes)
	}

	separator := o.Separator
	if separator == "" {
		separator = pick(f.rand, []string{"-", "_"})
	}

	if separator != "-" && separator != "_" {
		separator = "-"
	}

	return strings.ToLower(language) + separator + strings.ToUpper(country)
}

// Locale returns a generated locale using the default faker.
func Locale(opts ...LocaleOption) string {
	return defaultFaker().Locale(opts...)
}
