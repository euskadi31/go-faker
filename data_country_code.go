// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "math/rand/v2"

var dataCountryCodes = []string{
	"AE",
	"AR",
	"AT",
	"AU",
	"BE",
	"BR",
	"CA",
	"CH",
	"CN",
	"DE",
	"DK",
	"ES",
	"FI",
	"FR",
	"GB",
	"IE",
	"IN",
	"IT",
	"JP",
	"KR",
	"MX",
	"NL",
	"NO",
	"NZ",
	"PL",
	"PT",
	"RO",
	"SE",
	"SG",
	"TR",
	"UA",
	"US",
	"ZA",
}

func pickCountryCode(rng *rand.Rand) string {
	return pick(rng, dataCountryCodes)
}
