// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

type Facet interface {
	Name() string
	Type() string
	Source() string
}

type Facets struct {
	Facets []*Facet
}

func (f *Facets) ByName(name string) (*Facet, bool) {
	// TODO implement this
	return nil, false
}
