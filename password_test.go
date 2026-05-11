// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"errors"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func countAny(s, set string) int {
	n := 0

	for i := 0; i < len(s); i++ {
		if strings.IndexByte(set, s[i]) >= 0 {
			n++
		}
	}

	return n
}

func TestPasswordDefault(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	p, err := f.Password(DefaultPasswordOptions())
	require.NoError(t, err)
	assert.Len(t, p, 16)

	assert.GreaterOrEqual(t, countAny(p, passwordUpperChars), 1)
	assert.GreaterOrEqual(t, countAny(p, passwordLowerChars), 1)
	assert.GreaterOrEqual(t, countAny(p, passwordDigitChars), 1)
	assert.GreaterOrEqual(t, countAny(p, defaultPasswordSpec), 1)
}

func TestPasswordCustomPolicy(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(7, 11))))

	opts := PasswordOptions{
		Length:     24,
		MinUpper:   3,
		MinLower:   3,
		MinDigit:   3,
		MinSpecial: 3,
		Specials:   "!@",
	}

	p, err := f.Password(opts)
	require.NoError(t, err)
	assert.Len(t, p, 24)
	assert.GreaterOrEqual(t, countAny(p, passwordUpperChars), 3)
	assert.GreaterOrEqual(t, countAny(p, passwordLowerChars), 3)
	assert.GreaterOrEqual(t, countAny(p, passwordDigitChars), 3)
	assert.GreaterOrEqual(t, countAny(p, "!@"), 3)
}

func TestPasswordValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		opts PasswordOptions
		err  error
	}{
		{"zero length", PasswordOptions{Length: 0, MinUpper: 1}, ErrPasswordLength},
		{"negative min", PasswordOptions{Length: 8, MinUpper: -1}, ErrPasswordNegative},
		{"length too short", PasswordOptions{Length: 2, MinUpper: 2, MinLower: 2}, ErrPasswordTooShort},
		{"empty specials", PasswordOptions{Length: 8, MinSpecial: 1, Specials: ""}, ErrPasswordSpecials},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New().Password(tc.opts)
			require.Error(t, err)
			assert.True(t, errors.Is(err, tc.err), "got %v, want %v", err, tc.err)
		})
	}
}

func TestMustPasswordPanicsOnInvalid(t *testing.T) {
	f := New()

	assert.Panics(t, func() {
		f.MustPassword(PasswordOptions{Length: 0})
	})
}

func TestPasswordDeterministic(t *testing.T) {
	rng1 := rand.New(rand.NewPCG(42, 42))
	rng2 := rand.New(rand.NewPCG(42, 42))

	p1, err := New(WithRand(rng1)).Password(DefaultPasswordOptions())
	require.NoError(t, err)
	p2, err := New(WithRand(rng2)).Password(DefaultPasswordOptions())
	require.NoError(t, err)

	assert.Equal(t, p1, p2)
}

func TestPackagePassword(t *testing.T) {
	p, err := Password(DefaultPasswordOptions())
	require.NoError(t, err)
	assert.Len(t, p, 16)
}

func BenchmarkPassword(b *testing.B) {
	f := New()
	opts := DefaultPasswordOptions()

	for b.Loop() {
		_, _ = f.Password(opts)
	}
}
