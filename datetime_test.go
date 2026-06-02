// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dateLayout = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// --- Date ---------------------------------------------------------------------

func TestDate(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	date := f.MustDate()

	assert.Regexp(t, dateLayout, date)

	parsed, err := time.Parse("2006-01-02", date)
	require.NoError(t, err)

	assert.False(t, parsed.Before(defaultDateFrom), "date %q before default from", date)
	assert.False(t, parsed.After(defaultDateTo), "date %q after default to", date)
}

func TestDateCustomRange(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	from := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2006, 12, 31, 0, 0, 0, 0, time.UTC)

	for range 100 {
		date := f.MustDate(WithDateRange(from, to))

		parsed, err := time.Parse("2006-01-02", date)
		require.NoError(t, err)

		assert.False(t, parsed.Before(from), "date %q before from", date)
		assert.False(t, parsed.After(to), "date %q after to", date)
	}
}

func TestDateInclusiveSingleDay(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	day := time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, "2024-07-18", f.MustDate(WithDateRange(day, day)))
}

func TestDateNonUTCLocationConvertedToUTC(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	est := time.FixedZone("EST", -5*3600)

	// 2024-01-01 23:00 -05:00 == 2024-01-02 04:00 UTC, so the UTC date is the 2nd.
	day := time.Date(2024, 1, 1, 23, 0, 0, 0, est)

	assert.Equal(t, "2024-01-02", f.MustDate(WithDateRange(day, day)))
}

func TestDateDeterministicSameSeed(t *testing.T) {
	from := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2006, 12, 31, 0, 0, 0, 0, time.UTC)

	a := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))
	b := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	for i := range 50 {
		assert.Equal(t, a.MustDate(WithDateRange(from, to)), b.MustDate(WithDateRange(from, to)), "iter %d", i)
	}
}

func TestDateDifferentSeed(t *testing.T) {
	from := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2100, 12, 31, 0, 0, 0, 0, time.UTC)

	a := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	b := New(WithRand(rand.New(rand.NewPCG(3, 4))))

	var seqA, seqB []string

	for range 20 {
		seqA = append(seqA, a.MustDate(WithDateRange(from, to)))
		seqB = append(seqB, b.MustDate(WithDateRange(from, to)))
	}

	assert.NotEqual(t, seqA, seqB)
}

func TestDateInvalidRange(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	from := time.Date(2006, 12, 31, 0, 0, 0, 0, time.UTC)
	to := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := f.Date(WithDateRange(from, to))
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidRange))

	assert.Panics(t, func() { f.MustDate(WithDateRange(from, to)) })
}

// --- DateTime -----------------------------------------------------------------

func TestDateTime(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	value := f.MustDateTime()

	parsed, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)

	assert.True(t, value[len(value)-1] == 'Z', "datetime %q must end with Z", value)
	assert.False(t, parsed.Before(defaultDateTimeFrom), "datetime %q before default from", value)
	assert.False(t, parsed.After(defaultDateTimeTo), "datetime %q after default to", value)
}

func TestDateTimeCustomRange(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	for range 100 {
		value := f.MustDateTime(WithDateTimeRange(from, to))

		parsed, err := time.Parse(time.RFC3339, value)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, parsed.Unix(), from.Unix())
		assert.LessOrEqual(t, parsed.Unix(), to.Unix())
	}
}

func TestDateTimeInclusiveSingleInstant(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	instant := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, "2024-01-01T00:00:00Z", f.MustDateTime(WithDateTimeRange(instant, instant)))
}

func TestDateTimeDeterministicSameSeed(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	a := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))
	b := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	for i := range 50 {
		assert.Equal(t,
			a.MustDateTime(WithDateTimeRange(from, to)),
			b.MustDateTime(WithDateTimeRange(from, to)),
			"iter %d", i)
	}
}

func TestDateTimeDifferentSeed(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	b := New(WithRand(rand.New(rand.NewPCG(3, 4))))

	var seqA, seqB []string

	for range 20 {
		seqA = append(seqA, a.MustDateTime())
		seqB = append(seqB, b.MustDateTime())
	}

	assert.NotEqual(t, seqA, seqB)
}

func TestDateTimeInvalidRange(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	from := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := f.DateTime(WithDateTimeRange(from, to))
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidRange))

	assert.Panics(t, func() { f.MustDateTime(WithDateTimeRange(from, to)) })
}

// --- UnixTime -----------------------------------------------------------------

func TestUnixTime(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	value := f.MustUnixTime()

	assert.GreaterOrEqual(t, value, defaultDateTimeFrom.Unix())
	assert.LessOrEqual(t, value, defaultDateTimeTo.Unix())
}

func TestUnixTimeCustomRange(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	for range 100 {
		value := f.MustUnixTime(WithUnixTimeRange(from, to))

		assert.GreaterOrEqual(t, value, from.Unix())
		assert.LessOrEqual(t, value, to.Unix())
	}
}

func TestUnixTimeInclusiveSingleInstant(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	instant := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, int64(1704067200), f.MustUnixTime(WithUnixTimeRange(instant, instant)))
}

func TestUnixTimeDeterministicSameSeed(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	a := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))
	b := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	for i := range 50 {
		assert.Equal(t,
			a.MustUnixTime(WithUnixTimeRange(from, to)),
			b.MustUnixTime(WithUnixTimeRange(from, to)),
			"iter %d", i)
	}
}

func TestUnixTimeDifferentSeed(t *testing.T) {
	a := New(WithRand(rand.New(rand.NewPCG(1, 2))))
	b := New(WithRand(rand.New(rand.NewPCG(3, 4))))

	var seqA, seqB []int64

	for range 20 {
		seqA = append(seqA, a.MustUnixTime())
		seqB = append(seqB, b.MustUnixTime())
	}

	assert.NotEqual(t, seqA, seqB)
}

func TestUnixTimeInvalidRange(t *testing.T) {
	f := New(WithRand(rand.New(rand.NewPCG(1, 2))))

	from := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := f.UnixTime(WithUnixTimeRange(from, to))
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidRange))

	assert.Panics(t, func() { f.MustUnixTime(WithUnixTimeRange(from, to)) })
}

// --- Package-level helpers ----------------------------------------------------

func TestDateHelper(t *testing.T) {
	date := Date()

	assert.Regexp(t, dateLayout, date)

	parsed, err := time.Parse("2006-01-02", date)
	require.NoError(t, err)

	assert.False(t, parsed.Before(defaultDateFrom))
	assert.False(t, parsed.After(defaultDateTo))
}

func TestDateTimeHelper(t *testing.T) {
	value := DateTime()

	parsed, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)

	assert.True(t, value[len(value)-1] == 'Z', "datetime %q must end with Z", value)
	assert.False(t, parsed.Before(defaultDateTimeFrom))
	assert.False(t, parsed.After(defaultDateTimeTo))
}

func TestUnixTimeHelper(t *testing.T) {
	value := UnixTime()

	assert.GreaterOrEqual(t, value, defaultDateTimeFrom.Unix())
	assert.LessOrEqual(t, value, defaultDateTimeTo.Unix())
}

// --- Examples -----------------------------------------------------------------

func ExampleFaker_Date() {
	from := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2006, 12, 31, 0, 0, 0, 0, time.UTC)

	f := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	fmt.Println(f.MustDate(WithDateRange(from, to)))
	// Output: 1999-02-10
}

func ExampleFaker_DateTime() {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	f := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	fmt.Println(f.MustDateTime(WithDateTimeRange(from, to)))
	// Output: 2024-09-16T01:44:22Z
}

func ExampleFaker_UnixTime() {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	f := New(WithRand(rand.New(rand.NewPCG(1234, 5678))))

	fmt.Println(f.MustUnixTime(WithUnixTimeRange(from, to)))
	// Output: 1726451062
}

// --- Benchmarks ---------------------------------------------------------------

func BenchmarkDate(b *testing.B) {
	for b.Loop() {
		Date()
	}
}

func BenchmarkDateTime(b *testing.B) {
	for b.Loop() {
		DateTime()
	}
}

func BenchmarkUnixTime(b *testing.B) {
	for b.Loop() {
		UnixTime()
	}
}
