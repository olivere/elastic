// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"strings"
)

// Removes all characters unless they're
// between greater or equal to 0x2d (-) or
// less than or equal to (z) or a tilde (~).
func cleanPathString(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 0x2d && r <= 0x7a || r == '~' {
			return r
		}
		return -1
	}, s)
}
