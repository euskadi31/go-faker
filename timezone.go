// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

// Timezone returns a randomly selected IANA tzdb name or alias (e.g.
// "Europe/Paris", "America/New_York", "UTC", "Eire"). The returned value
// is loadable via time.LoadLocation.
func (f *Faker) Timezone() string {
	return pickTZ(f.rand)
}

// Timezone returns a generated IANA tzdb name or alias using the default faker.
// The returned value is loadable via time.LoadLocation.
func Timezone() string {
	return defaultFaker().Timezone()
}
