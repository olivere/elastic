// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Range Facet
// See: http://www.elasticsearch.org/guide/reference/api/search/facets/range-facet.html
type RangeFacet struct {
	Facet
	global     *bool
	field      string
	keyField   string
	valueField string
	ranges     []rangeFacetRange
}

type rangeFacetRange struct {
	From *float64
	To   *float64
}

func NewRangeFacet(field string) RangeFacet {
	return RangeFacet{
		field:  field,
		ranges: make([]rangeFacetRange, 0),
	}
}

func (f RangeFacet) Global(global bool) RangeFacet {
	f.global = &global
	return f
}

func (f RangeFacet) Lt(to float64) RangeFacet {
	f.ranges = append(f.ranges, rangeFacetRange{From: nil, To: &to})
	return f
}

func (f RangeFacet) Between(from, to float64) RangeFacet {
	f.ranges = append(f.ranges, rangeFacetRange{From: &from, To: &to})
	return f
}

func (f RangeFacet) Gt(from float64) RangeFacet {
	f.ranges = append(f.ranges, rangeFacetRange{From: &from, To: nil})
	return f
}

func (f RangeFacet) KeyField(keyField string) RangeFacet {
	f.keyField = keyField
	return f
}

func (f RangeFacet) ValueField(valueField string) RangeFacet {
	f.valueField = valueField
	return f
}

func (f RangeFacet) Source() interface{} {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["range"] = opts

	if f.keyField != "" {
		opts["key_field"] = f.keyField
		opts["value_field"] = f.valueField
	} else {
		opts["field"] = f.field
	}

	if f.global != nil {
		opts["global"] = *f.global
	}

	ranges := make([]interface{}, 0)

	for _, rng := range f.ranges {
		r := make(map[string]interface{})
		if rng.From != nil {
			r["from"] = *rng.From
		}
		if rng.To != nil {
			r["to"] = *rng.To
		}
		ranges = append(ranges, r)
	}

	opts["ranges"] = ranges

	return source
}
