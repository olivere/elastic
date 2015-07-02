// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Filter defines the contract each filter has to fulfill.
type Filter interface {
	// Source returns a JSON-serializable fragment of the request.
	Source() (interface{}, error)
}
