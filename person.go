// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "math/rand"

// PersonFlags control the behavior of the generator.
type PersonFlags uint16

// PersonFlags enums.
const (
	PersonFirstName PersonFlags = 1 << iota
	PersonLastName
	PersonMan
	PersonWoman

	PersonAll = PersonFirstName | PersonLastName | PersonMan | PersonWoman
)

// PersonInfo struct.
type PersonInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
}

func (p PersonInfo) String() string {
	if p.FirstName != "" && p.LastName != "" {
		return p.FirstName + " " + p.LastName
	}

	if p.FirstName != "" {
		return p.FirstName
	}

	return p.LastName
}

// PersonGenerator struct.
type PersonGenerator struct {
	Flags PersonFlags
}

// NewPersonGenerator return PersonGenerator.
func NewPersonGenerator() *PersonGenerator {
	return &PersonGenerator{
		Flags: PersonAll,
	}
}

// Generate person.
func (g *PersonGenerator) Generate() *PersonInfo {
	person := &PersonInfo{}

	if g.Flags&PersonLastName != 0 {
		person.LastName = getLastName()
	}

	genders := []string{}
	firstNames := []string{}

	if g.Flags&PersonFirstName != 0 && g.Flags&PersonMan != 0 {
		firstNames = append(firstNames, getManFirstName())
		genders = append(genders, "Male")
	}

	if g.Flags&PersonFirstName != 0 && g.Flags&PersonWoman != 0 {
		firstNames = append(firstNames, getWomanFirstName())
		genders = append(genders, "Female")
	}

	if n := len(firstNames); n > 1 {
		//nolint:gosec
		i := rand.Intn(n - 1)

		person.FirstName = firstNames[i]
		person.Gender = genders[i]
	} else if n == 1 {
		person.FirstName = firstNames[0]
		person.Gender = genders[0]
	}

	return person
}

// Person return an generated person.
func Person() *PersonInfo {
	return NewPersonGenerator().Generate()
}
