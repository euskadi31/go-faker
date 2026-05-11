// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "math/rand/v2"

var dataLanguageCodes = []string{
	"af",
	"ar",
	"bg",
	"ca",
	"cs",
	"da",
	"de",
	"el",
	"en",
	"es",
	"et",
	"fi",
	"fr",
	"he",
	"hi",
	"hr",
	"hu",
	"id",
	"it",
	"ja",
	"ko",
	"lt",
	"lv",
	"nl",
	"no",
	"pl",
	"pt",
	"ro",
	"ru",
	"sk",
	"sl",
	"sv",
	"th",
	"tr",
	"uk",
	"vi",
	"zh",
}

func pickLanguageCode(rng *rand.Rand) string {
	return pick(rng, dataLanguageCodes)
}
