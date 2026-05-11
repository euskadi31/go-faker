// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import (
	"fmt"
	"strconv"
	"strings"
)

// macAddressBytes is the canonical length of a MAC address in bytes.
const macAddressBytes = 6

// MACAddressOptions controls how a MAC address is rendered.
type MACAddressOptions struct {
	Prefix    string
	Separator string
	Lowercase bool
}

// MACAddressOption mutates MACAddressOptions.
type MACAddressOption func(*MACAddressOptions)

// WithMACPrefix replaces the first bytes of the generated MAC with those
// described by prefix. The prefix may use ':' or '-' as a separator and is
// case-insensitive (e.g. "FF", "ff:ff", "0a-1b-2c"). Invalid prefixes are
// ignored.
func WithMACPrefix(prefix string) MACAddressOption {
	return func(o *MACAddressOptions) {
		o.Prefix = prefix
	}
}

// WithMACSeparator sets the byte separator (default ":").
func WithMACSeparator(sep string) MACAddressOption {
	return func(o *MACAddressOptions) {
		o.Separator = sep
	}
}

// WithMACLowercase renders the MAC address using lowercase hex digits.
func WithMACLowercase() MACAddressOption {
	return func(o *MACAddressOptions) {
		o.Lowercase = true
	}
}

// WithMACUppercase renders the MAC address using uppercase hex digits (the
// default).
func WithMACUppercase() MACAddressOption {
	return func(o *MACAddressOptions) {
		o.Lowercase = false
	}
}

// MACAddress returns a randomly generated MAC address.
func (f *Faker) MACAddress(opts ...MACAddressOption) string {
	o := MACAddressOptions{Separator: ":"}

	for _, opt := range opts {
		opt(&o)
	}

	if o.Separator == "" {
		o.Separator = ":"
	}

	bytes := make([]byte, macAddressBytes)
	fillBytes(f.rand, bytes)

	if prefix := parseMACPrefix(o.Prefix); len(prefix) > 0 {
		copy(bytes, prefix)
	}

	format := "%02X"
	if o.Lowercase {
		format = "%02x"
	}

	parts := make([]string, macAddressBytes)
	for i, b := range bytes {
		parts[i] = fmt.Sprintf(format, b)
	}

	return strings.Join(parts, o.Separator)
}

// MACAddress returns a generated MAC address using the default faker.
func MACAddress(opts ...MACAddressOption) string {
	return defaultFaker().MACAddress(opts...)
}

// parseMACPrefix decodes up to 6 hex bytes from prefix. Accepted separators
// inside the prefix are ':' and '-'. Returns nil for malformed input.
func parseMACPrefix(prefix string) []byte {
	if prefix == "" {
		return nil
	}

	normalized := strings.ReplaceAll(prefix, "-", ":")
	parts := strings.Split(normalized, ":")

	if len(parts) > macAddressBytes {
		return nil
	}

	out := make([]byte, 0, len(parts))

	for _, p := range parts {
		if len(p) != 2 {
			return nil
		}

		v, err := strconv.ParseUint(p, 16, 8)
		if err != nil {
			return nil
		}

		out = append(out, byte(v))
	}

	return out
}

// MacAddressGenerator is the legacy generator.
//
// Deprecated: use Faker.MACAddress.
type MacAddressGenerator struct {
	Prefix    string
	Separator string
}

// NewMacAddressGenerator returns a legacy MAC address generator.
//
// Deprecated: use New() and Faker.MACAddress.
func NewMacAddressGenerator() *MacAddressGenerator {
	return &MacAddressGenerator{Separator: ":"}
}

// Generate produces a MAC address using the legacy generator fields.
//
// Deprecated: use Faker.MACAddress.
func (g *MacAddressGenerator) Generate() string {
	opts := []MACAddressOption{}

	if g.Prefix != "" {
		opts = append(opts, WithMACPrefix(g.Prefix))
	}

	if g.Separator != "" && g.Separator != ":" {
		opts = append(opts, WithMACSeparator(g.Separator))
	}

	return defaultFaker().MACAddress(opts...)
}

// MacAddress returns a generated MAC address.
//
// Deprecated: use MACAddress.
func MacAddress() string {
	return MACAddress()
}
