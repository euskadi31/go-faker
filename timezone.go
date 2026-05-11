// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

// Timezone returns a randomly selected IANA time zone name (e.g.
// "Europe/Paris", "America/New_York").
func (f *Faker) Timezone() string {
	return pickTZ(f.rand)
}

// Timezone returns a generated time zone name using the default faker.
func Timezone() string {
	return defaultFaker().Timezone()
}
