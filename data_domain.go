// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package faker

import "github.com/euskadi31/go-reggen"

func getFakeDomain() string {
	tld := getTLD()

	name, _ := reggen.Generate("[a-z0-9]{1,8}([a-z]|-)[a-z0-9]")

	return name + "." + tld
}
