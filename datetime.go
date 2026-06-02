// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"errors"
	"time"
)

// secondsPerDay is the number of seconds in a 24-hour UTC day. UTC has no DST,
// so the gap between two UTC midnights is always an exact multiple of it.
const secondsPerDay = 86400

// ErrInvalidRange is returned by the date and time generators when the supplied
// from value is after the to value. The equal case (from == to) is valid and
// yields that exact instant or date.
var ErrInvalidRange = errors.New("invalid range: from is after to")

// Default generation ranges. All bounds are inclusive and expressed in UTC.
var (
	defaultDateFrom     = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultDateTo       = time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)
	defaultDateTimeFrom = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultDateTimeTo   = time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)
)

// truncateToUTCDate converts t to UTC and zeroes the clock, returning UTC
// midnight of the same calendar day.
func truncateToUTCDate(t time.Time) time.Time {
	t = t.UTC()

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// DateOptions controls how Date values are generated. Both bounds are inclusive
// and normalized to whole UTC days.
type DateOptions struct {
	From time.Time
	To   time.Time
}

// DateOption mutates DateOptions.
type DateOption func(*DateOptions)

// WithDateRange sets the inclusive [from, to] range for a generated date. Values
// carrying a non-UTC location are converted to UTC before the date is taken.
func WithDateRange(from, to time.Time) DateOption {
	return func(o *DateOptions) {
		o.From = from
		o.To = to
	}
}

// Date returns a random calendar date formatted as "2006-01-02". The range is
// inclusive and operates by whole UTC days; from and to are normalized to UTC
// midnight, so any clock component is ignored. It returns ErrInvalidRange when
// from is after to. Output is deterministic for a Faker built with WithRand.
func (f *Faker) Date(opts ...DateOption) (string, error) {
	o := DateOptions{From: defaultDateFrom, To: defaultDateTo}

	for _, opt := range opts {
		opt(&o)
	}

	from := truncateToUTCDate(o.From)
	to := truncateToUTCDate(o.To)

	if from.After(to) {
		return "", ErrInvalidRange
	}

	days := (to.Unix() - from.Unix()) / secondsPerDay
	offset := randomInt64Inclusive(f.rand, 0, days)

	//nolint:gosec // offset is bounded by the validated day count and fits int for practical ranges.
	return from.AddDate(0, 0, int(offset)).Format("2006-01-02"), nil
}

// MustDate is like Date but panics when the range is invalid.
func (f *Faker) MustDate(opts ...DateOption) string {
	s, err := f.Date(opts...)
	if err != nil {
		panic(err)
	}

	return s
}

// Date returns a random date in the default range using the default faker.
func Date() string {
	return defaultFaker().MustDate()
}

// DateTimeOptions controls how DateTime values are generated. Both bounds are
// inclusive and operate at second granularity.
type DateTimeOptions struct {
	From time.Time
	To   time.Time
}

// DateTimeOption mutates DateTimeOptions.
type DateTimeOption func(*DateTimeOptions)

// WithDateTimeRange sets the inclusive [from, to] range for a generated
// datetime. Values carrying a non-UTC location are converted to UTC.
func WithDateTimeRange(from, to time.Time) DateTimeOption {
	return func(o *DateTimeOptions) {
		o.From = from
		o.To = to
	}
}

// DateTime returns a random UTC instant formatted as RFC3339 (e.g.
// "2024-07-18T12:34:56Z"). The range is inclusive and operates by seconds. It
// returns ErrInvalidRange when from is after to. Output is deterministic for a
// Faker built with WithRand.
func (f *Faker) DateTime(opts ...DateTimeOption) (string, error) {
	o := DateTimeOptions{From: defaultDateTimeFrom, To: defaultDateTimeTo}

	for _, opt := range opts {
		opt(&o)
	}

	from := o.From.UTC()
	to := o.To.UTC()

	if from.Unix() > to.Unix() {
		return "", ErrInvalidRange
	}

	unix := randomInt64Inclusive(f.rand, from.Unix(), to.Unix())

	return time.Unix(unix, 0).UTC().Format(time.RFC3339), nil
}

// MustDateTime is like DateTime but panics when the range is invalid.
func (f *Faker) MustDateTime(opts ...DateTimeOption) string {
	s, err := f.DateTime(opts...)
	if err != nil {
		panic(err)
	}

	return s
}

// DateTime returns a random RFC3339 datetime in the default range using the
// default faker.
func DateTime() string {
	return defaultFaker().MustDateTime()
}

// UnixTimeOptions controls how UnixTime values are generated. Both bounds are
// inclusive and operate at second granularity.
type UnixTimeOptions struct {
	From time.Time
	To   time.Time
}

// UnixTimeOption mutates UnixTimeOptions.
type UnixTimeOption func(*UnixTimeOptions)

// WithUnixTimeRange sets the inclusive [from, to] range for a generated Unix
// timestamp. Values carrying a non-UTC location are converted to UTC.
func WithUnixTimeRange(from, to time.Time) UnixTimeOption {
	return func(o *UnixTimeOptions) {
		o.From = from
		o.To = to
	}
}

// UnixTime returns a random Unix timestamp in seconds. The range is inclusive
// and operates by seconds. It returns ErrInvalidRange when from is after to.
// Output is deterministic for a Faker built with WithRand.
func (f *Faker) UnixTime(opts ...UnixTimeOption) (int64, error) {
	o := UnixTimeOptions{From: defaultDateTimeFrom, To: defaultDateTimeTo}

	for _, opt := range opts {
		opt(&o)
	}

	from := o.From.UTC()
	to := o.To.UTC()

	if from.Unix() > to.Unix() {
		return 0, ErrInvalidRange
	}

	return randomInt64Inclusive(f.rand, from.Unix(), to.Unix()), nil
}

// MustUnixTime is like UnixTime but panics when the range is invalid.
func (f *Faker) MustUnixTime(opts ...UnixTimeOption) int64 {
	v, err := f.UnixTime(opts...)
	if err != nil {
		panic(err)
	}

	return v
}

// UnixTime returns a random Unix timestamp in the default range using the
// default faker.
func UnixTime() int64 {
	return defaultFaker().MustUnixTime()
}
