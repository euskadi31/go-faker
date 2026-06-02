# go-faker

[![Go Reference](https://pkg.go.dev/badge/github.com/euskadi31/go-faker.svg)](https://pkg.go.dev/github.com/euskadi31/go-faker)
[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-faker)](https://goreportcard.com/report/github.com/euskadi31/go-faker)
[![Build Status](https://github.com/euskadi31/go-faker/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/euskadi31/go-faker/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/euskadi31/go-faker/badge.svg?branch=master)](https://coveralls.io/github/euskadi31/go-faker?branch=master)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

A small, dependency-light fake data generator for Go. It produces realistic test
fixtures (people, emails, passwords, dates, MAC addresses, locales, timezones,
random bytes) behind a functional-options API.

Output is non-deterministic by default (seeded from `crypto/rand`) and fully
reproducible when you inject your own `*rand.Rand`.

## Requirements

- Go 1.26 or later (uses `math/rand/v2`).

## Install

```bash
go get github.com/euskadi31/go-faker
```

```go
import faker "github.com/euskadi31/go-faker"
```

## Quick start

Every generator is available in two forms:

- A method on a `*Faker` instance, for full control over the random source.
- A package-level helper backed by a lazily-initialized default `Faker`.

```go
package main

import (
	"fmt"

	faker "github.com/euskadi31/go-faker"
)

func main() {
	// Package-level helpers (default faker, seeded from crypto/rand).
	fmt.Println(faker.Person())   // John Doe
	fmt.Println(faker.Email())    // j.doe@example.com
	fmt.Println(faker.Timezone()) // Europe/Paris

	// Or your own instance.
	f := faker.New()
	fmt.Println(f.MACAddress()) // 0A:1B:2C:3D:4E:5F
}
```

## Generators

### Person

```go
f := faker.New()

p := f.Person()                              // random gender
m := f.Person(faker.WithGender(faker.GenderMale))
w := f.Person(faker.WithGender(faker.GenderFemale))

fmt.Println(p.FirstName, p.LastName, p.Gender)
fmt.Println(p.String()) // "FirstName LastName"
```

`PersonInfo` is JSON-friendly:

```go
b, _ := json.Marshal(f.Person())
// {"first_name":"Jane","last_name":"Doe","gender":"Female"}
```

### Email

The default domain is `example.com` (RFC 2606 reserved), so generated addresses
cannot collide with real accounts.

```go
f := faker.New()

f.Email()                                        // j.doe@example.com
f.Email(faker.WithRealEmailDomain())             // pick a real provider (gmail.com, ...)
f.Email(faker.WithFakeEmailDomain())             // plausible but non-real domain
f.Email(faker.WithEmailDomain("acme.test"))      // fixed domain
f.Email(faker.WithEmailPrefix("qa-"))            // qa-j.doe@example.com
```

### Password

`Password` enforces a policy and returns an error when the policy is
inconsistent (for example, the sum of minimum counts exceeds the length).

```go
f := faker.New()

// Convenience defaults: 16 chars, at least 1 upper / lower / digit / special.
pwd := f.MustPassword(faker.DefaultPasswordOptions())

// Custom policy.
pwd, err := f.Password(faker.PasswordOptions{
	Length:     24,
	MinUpper:   2,
	MinLower:   2,
	MinDigit:   2,
	MinSpecial: 2,
	Specials:   "!@#$%",
})
if err != nil {
	// ErrPasswordLength, ErrPasswordTooShort, ErrPasswordSpecials, ...
}
```

### Date, DateTime and UnixTime

All bounds are inclusive and operate in UTC. `Date` works on whole UTC days
(clock components are ignored), `DateTime` and `UnixTime` work at second
granularity. Each generator returns `ErrInvalidRange` when `from` is after `to`,
and has a `Must` variant that panics instead.

```go
f := faker.New()

f.Date()      // "2007-03-14", default range 1970-01-01..2030-12-31
f.DateTime()  // "2024-07-18T12:34:56Z" (RFC3339)
f.UnixTime()  // 1721306096

from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)

d, err := f.Date(faker.WithDateRange(from, to))
ts, err := f.DateTime(faker.WithDateTimeRange(from, to))
n, err := f.UnixTime(faker.WithUnixTimeRange(from, to))

// Panic on invalid range instead of returning an error.
d = f.MustDate(faker.WithDateRange(from, to))
```

### Timezone

Returns an IANA tzdb name or alias loadable via `time.LoadLocation`.

```go
f := faker.New()

name := f.Timezone() // "America/New_York", "Europe/Paris", "UTC", ...
loc, _ := time.LoadLocation(name)
```

### Locale

Language and country subtags are picked independently, so valid but uncommon
combinations (English language with a French region, etc.) can occur.

```go
f := faker.New()

f.Locale()                                        // "fr-FR", "en_US", "en-FR", ...
f.Locale(faker.WithLocaleLanguage("fr"))          // force language (lowercased)
f.Locale(faker.WithLocaleCountry("US"))           // force country (uppercased)
f.Locale(faker.WithLocaleSeparator("_"))          // only "-" or "_" are accepted
```

### MAC address

```go
f := faker.New()

f.MACAddress()                                    // 0A:1B:2C:3D:4E:5F (uppercase, ":")
f.MACAddress(faker.WithMACLowercase())            // 0a:1b:2c:3d:4e:5f
f.MACAddress(faker.WithMACSeparator("-"))         // 0A-1B-2C-3D-4E-5F
f.MACAddress(faker.WithMACPrefix("0a:1b:2c"))     // fixes the first bytes
```

The prefix accepts `:` or `-` separators and is case-insensitive. Invalid
prefixes are ignored.

### Bytes

`Bytes` draws from the Faker's RNG (deterministic with `WithRand`). For
cryptographically secure output use `CryptoBytes`.

```go
f := faker.New()

b := f.Bytes(32)              // pseudo-random, reproducible with WithRand

secure, err := faker.CryptoBytes(32) // crypto/rand backed, surfaces errors
```

## Reproducible output

Pass your own `*rand.Rand` to make every generator deterministic. Seeding two
fakers with the same source produces identical sequences, which is ideal for
golden tests.

```go
import "math/rand/v2"

rng := rand.New(rand.NewPCG(1, 2))
f := faker.New(faker.WithRand(rng))

f.Email() // same value on every run with the same seed
```

`f.Rand()` exposes the underlying source if you need to share it with your own
generation code.

## Concurrency

A `Faker` is not safe for concurrent use by multiple goroutines. For
deterministic generation across goroutines, give each goroutine its own `Faker`
built with its own `*rand.Rand`.

The package-level helpers share a single lazily-initialized default `Faker` and
are subject to the same constraint.

## Testing

```bash
go test -race -cover ./...
go test -bench=. ./...
```

## License

go-faker is released under the MIT License. See the source file headers for
details.
