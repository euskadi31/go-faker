// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
)

const (
	passwordUpperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	passwordLowerChars  = "abcdefghijklmnopqrstuvwxyz"
	passwordDigitChars  = "0123456789"
	defaultPasswordSpec = "!@#$%^&*"
)

// Common password validation errors.
var (
	ErrPasswordLength   = errors.New("password length must be > 0")
	ErrPasswordNegative = errors.New("password minimum counts must be >= 0")
	ErrPasswordTooShort = errors.New("password length is smaller than the sum of minimums")
	ErrPasswordSpecials = errors.New("password specials must not be empty when MinSpecial > 0")
)

// PasswordOptions describes the policy a generated password must satisfy.
type PasswordOptions struct {
	Length     int
	MinUpper   int
	MinLower   int
	MinDigit   int
	MinSpecial int
	Specials   string
}

// DefaultPasswordOptions returns a 16-character policy requiring at least one
// upper case letter, one lower case letter, one digit and one special.
func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{
		Length:     16,
		MinUpper:   1,
		MinLower:   1,
		MinDigit:   1,
		MinSpecial: 1,
		Specials:   defaultPasswordSpec,
	}
}

// Validate checks that the policy is internally consistent.
func (o PasswordOptions) Validate() error {
	if o.Length <= 0 {
		return ErrPasswordLength
	}

	if o.MinUpper < 0 || o.MinLower < 0 || o.MinDigit < 0 || o.MinSpecial < 0 {
		return ErrPasswordNegative
	}

	if o.MinUpper+o.MinLower+o.MinDigit+o.MinSpecial > o.Length {
		return ErrPasswordTooShort
	}

	if o.MinSpecial > 0 && o.Specials == "" {
		return ErrPasswordSpecials
	}

	return nil
}

// Password returns a password matching the supplied policy.
func (f *Faker) Password(opts PasswordOptions) (string, error) {
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("faker: invalid password options: %w", err)
	}

	out := make([]byte, 0, opts.Length)

	out = appendRandomChars(out, f.rand, passwordUpperChars, opts.MinUpper)
	out = appendRandomChars(out, f.rand, passwordLowerChars, opts.MinLower)
	out = appendRandomChars(out, f.rand, passwordDigitChars, opts.MinDigit)

	specials := opts.Specials
	if opts.MinSpecial == 0 {
		specials = ""
	}

	out = appendRandomChars(out, f.rand, specials, opts.MinSpecial)

	var pool strings.Builder

	if opts.MinUpper > 0 || (opts.MinUpper == 0 && opts.MinLower == 0 && opts.MinDigit == 0 && opts.MinSpecial == 0) {
		pool.WriteString(passwordUpperChars)
	}

	if opts.MinLower > 0 {
		pool.WriteString(passwordLowerChars)
	}

	if opts.MinDigit > 0 {
		pool.WriteString(passwordDigitChars)
	}

	if opts.MinSpecial > 0 {
		pool.WriteString(opts.Specials)
	}

	poolStr := pool.String()
	if poolStr == "" {
		poolStr = passwordUpperChars + passwordLowerChars + passwordDigitChars
	}

	out = appendRandomChars(out, f.rand, poolStr, opts.Length-len(out))

	shuffle(f.rand, out)

	return string(out), nil
}

// MustPassword returns a password matching the policy or panics on invalid
// options.
func (f *Faker) MustPassword(opts PasswordOptions) string {
	p, err := f.Password(opts)
	if err != nil {
		panic(err)
	}

	return p
}

// Password returns a generated password from the default faker.
func Password(opts PasswordOptions) (string, error) {
	return defaultFaker().Password(opts)
}

// MustPassword returns a generated password from the default faker.
func MustPassword(opts PasswordOptions) string {
	return defaultFaker().MustPassword(opts)
}

func appendRandomChars(dst []byte, rng *rand.Rand, pool string, count int) []byte {
	if count <= 0 || pool == "" {
		return dst
	}

	for range count {
		dst = append(dst, pool[rng.IntN(len(pool))])
	}

	return dst
}
