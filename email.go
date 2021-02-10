// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand"
	"strings"
)

// EmailFlags control the behavior of the generator
type EmailFlags uint16

// EmailFlags enums
const (
	EmailRealDomain EmailFlags = 1 << iota
)

// EmailGenerator struct
type EmailGenerator struct {
	Flags EmailFlags
}

// NewEmailGenerator return EmailGenerator
func NewEmailGenerator() *EmailGenerator {
	return &EmailGenerator{
		Flags: EmailRealDomain,
	}
}

// Generate email
func (g *EmailGenerator) Generate() string {
	var domain string
	if g.Flags&EmailRealDomain != 0 {
		domain = getRealEmailDomain()
	} else {
		domain = getFakeDomain()
	}

	names := []string{}

	names = append(names, getEmailFirstNameAndLastNamePattern())
	names = append(names, getEmailFirstNameLetterAndLastNamePattern())

	//nolint:gosec
	i := rand.Intn(len(names) - 1)

	return names[i] + "@" + domain
}

func getEmailFirstNameLetterAndLastNamePattern() string {
	firstname := getFirstName()
	lastname := getLastName()

	value := firstname[0:1] + lastname

	return strings.ToLower(value)
}

func getEmailFirstNameAndLastNamePattern() string {
	firstname := getFirstName()
	lastname := getLastName()

	value := firstname + "." + lastname

	return strings.ToLower(value)
}

// Email return an generated email address
func Email() string {
	return NewEmailGenerator().Generate()
}
