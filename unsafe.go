// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import "unsafe"

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
