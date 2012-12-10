// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Terms Facet
// See: http://www.elasticsearch.org/guide/reference/api/search/facets/terms-facet.html
type TermsFacet struct {
	Facet
	global   *bool
	fields   []string
	size     *int
	order    string
	allTerms *bool
	exclude  []string
}

func NewTermsFacet(fields ...string) TermsFacet {
	f := TermsFacet{
		fields:  make([]string, 0),
		exclude: make([]string, 0),
	}
	f.fields = append(f.fields, fields...)
	return f
}

func (f TermsFacet) Global(global bool) TermsFacet {
	f.global = &global
	return f
}

func (f TermsFacet) Size(size int) TermsFacet {
	f.size = &size
	return f
}

// Valid order options are: "count" (default), "term",
// "reverse_count", and "reverse_term".
func (f TermsFacet) Order(order string) TermsFacet {
	f.order = order
	return f
}

func (f TermsFacet) AllTerms(allTerms bool) TermsFacet {
	f.allTerms = &allTerms
	return f
}

func (f TermsFacet) Exclude(exclude ...string) TermsFacet {
	f.exclude = append(f.exclude, exclude...)
	return f
}

func (f TermsFacet) Source() interface{} {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["terms"] = opts

	if len(f.fields) == 1 {
		opts["field"] = f.fields[0]
	} else if len(f.fields) > 1 {
		opts["fields"] = f.fields
	}

	if f.global != nil {
		opts["global"] = *f.global
	}

	if f.size != nil {
		opts["size"] = *f.size
	}

	if f.order != "" {
		opts["order"] = f.order
	}

	if f.allTerms != nil {
		opts["all_terms"] = *f.allTerms
	}

	if len(f.exclude) > 0 {
		opts["exclude"] = f.exclude
	}

	return source
}
