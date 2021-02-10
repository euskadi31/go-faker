// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersonGenerator(t *testing.T) {

	{
		g := NewPersonGenerator()
		g.Flags = PersonFirstName

		person := g.Generate()

		assert.Empty(t, person.FirstName)
		assert.Empty(t, person.LastName)
		assert.Empty(t, person.Gender)
	}

	{
		g := NewPersonGenerator()
		g.Flags = PersonFirstName | PersonMan

		person := g.Generate()

		assert.NotEmpty(t, person.FirstName)
		assert.Empty(t, person.LastName)
		assert.Equal(t, "Male", person.Gender)
	}

	{
		g := NewPersonGenerator()
		g.Flags = PersonFirstName | PersonWoman

		person := g.Generate()

		assert.NotEmpty(t, person.FirstName)
		assert.Empty(t, person.LastName)
		assert.Equal(t, "Female", person.Gender)

		assert.Equal(t, person.FirstName, person.String())
	}

	{
		g := NewPersonGenerator()
		g.Flags = PersonLastName

		person := g.Generate()

		assert.Empty(t, person.FirstName)
		assert.NotEmpty(t, person.LastName)
		assert.Empty(t, person.Gender)

		assert.Equal(t, person.LastName, person.String())
	}

	{
		g := NewPersonGenerator()
		g.Flags = PersonLastName | PersonFirstName | PersonMan

		person := g.Generate()

		assert.NotEmpty(t, person.FirstName)
		assert.NotEmpty(t, person.LastName)
		assert.Equal(t, "Male", person.Gender)

		assert.Equal(t, person.FirstName+" "+person.LastName, person.String())
	}
}

func TestPerson(t *testing.T) {
	person := Person()

	assert.True(t, regexp.MustCompile("(Male|Female)").MatchString(person.Gender))
}

func BenchmarkPerson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Person()
	}
}

func BenchmarkPersonGenerator(b *testing.B) {
	g := NewPersonGenerator()
	g.Flags = PersonFirstName | PersonLastName | PersonWoman

	for i := 0; i < b.N; i++ {
		g.Generate()
	}
}
