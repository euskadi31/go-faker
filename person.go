// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

// Gender enumerates the genders supported when generating a Person.
type Gender string

// Supported Gender values.
const (
	// GenderAny lets the generator pick male or female at random.
	GenderAny Gender = ""
	// GenderMale forces a male first name.
	GenderMale Gender = "Male"
	// GenderFemale forces a female first name.
	GenderFemale Gender = "Female"
)

// PersonOptions controls how Person values are generated.
type PersonOptions struct {
	Gender Gender
}

// PersonOption mutates PersonOptions.
type PersonOption func(*PersonOptions)

// WithGender forces the generated person's gender.
func WithGender(g Gender) PersonOption {
	return func(o *PersonOptions) {
		o.Gender = g
	}
}

// PersonInfo describes a generated person.
type PersonInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
}

// String returns "FirstName LastName" when both are set, otherwise whichever
// field is populated.
func (p PersonInfo) String() string {
	switch {
	case p.FirstName != "" && p.LastName != "":
		return p.FirstName + " " + p.LastName
	case p.FirstName != "":
		return p.FirstName
	default:
		return p.LastName
	}
}

// Person returns a randomly generated person.
func (f *Faker) Person(opts ...PersonOption) PersonInfo {
	o := PersonOptions{Gender: GenderAny}

	for _, opt := range opts {
		opt(&o)
	}

	gender := o.Gender
	if gender == GenderAny {
		if f.rand.IntN(2) == 0 {
			gender = GenderMale
		} else {
			gender = GenderFemale
		}
	}

	var firstName string

	switch gender {
	case GenderMale, GenderAny:
		firstName = pickManFirstName(f.rand)
	case GenderFemale:
		firstName = pickWomanFirstName(f.rand)
	}

	return PersonInfo{
		FirstName: firstName,
		LastName:  pickLastName(f.rand),
		Gender:    string(gender),
	}
}

// Person returns a generated person using the default faker.
func Person(opts ...PersonOption) PersonInfo {
	return defaultFaker().Person(opts...)
}

// PersonFlags controls the legacy PersonGenerator behavior.
//
// Deprecated: use Faker.Person and WithGender.
type PersonFlags uint16

// Deprecated person flag values.
const (
	// Deprecated: included in Faker.Person by default.
	PersonFirstName PersonFlags = 1 << iota
	// Deprecated: included in Faker.Person by default.
	PersonLastName
	// Deprecated: use WithGender(GenderMale).
	PersonMan
	// Deprecated: use WithGender(GenderFemale).
	PersonWoman

	// Deprecated: equivalent to default Faker.Person behavior.
	PersonAll = PersonFirstName | PersonLastName | PersonMan | PersonWoman
)

// PersonGenerator is the legacy generator.
//
// Deprecated: use Faker.Person.
type PersonGenerator struct {
	Flags PersonFlags
}

// NewPersonGenerator returns a legacy PersonGenerator initialized with
// PersonAll flags.
//
// Deprecated: use New() and Faker.Person.
func NewPersonGenerator() *PersonGenerator {
	return &PersonGenerator{Flags: PersonAll}
}

// Generate produces a PersonInfo using the legacy flag semantics. Fields not
// requested by the flags are left empty.
//
// Deprecated: use Faker.Person.
func (g *PersonGenerator) Generate() *PersonInfo {
	f := defaultFaker()
	person := &PersonInfo{}

	if g.Flags&PersonLastName != 0 {
		person.LastName = pickLastName(f.rand)
	}

	wantMan := g.Flags&PersonFirstName != 0 && g.Flags&PersonMan != 0
	wantWoman := g.Flags&PersonFirstName != 0 && g.Flags&PersonWoman != 0

	switch {
	case wantMan && wantWoman:
		if f.rand.IntN(2) == 0 {
			person.FirstName = pickManFirstName(f.rand)
			person.Gender = string(GenderMale)
		} else {
			person.FirstName = pickWomanFirstName(f.rand)
			person.Gender = string(GenderFemale)
		}
	case wantMan:
		person.FirstName = pickManFirstName(f.rand)
		person.Gender = string(GenderMale)
	case wantWoman:
		person.FirstName = pickWomanFirstName(f.rand)
		person.Gender = string(GenderFemale)
	}

	return person
}
