// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Date Histogram Facet
// See: http://www.elasticsearch.org/guide/reference/api/search/facets/date-histogram-facet.html
type DateHistogramFacet struct {
	Facet
	global     *bool
	field      string
	keyField   string
	valueField string
	interval   string
	postZone   string
	preZone    string
	factor     *int
	preOffset  string
	postOffset string
}

func NewDateHistogramFacet(field string) DateHistogramFacet {
	return DateHistogramFacet{
		field: field,
	}
}

func (f DateHistogramFacet) Global(global bool) DateHistogramFacet {
	f.global = &global
	return f
}

func (f DateHistogramFacet) KeyField(keyField string) DateHistogramFacet {
	f.keyField = keyField
	return f
}

func (f DateHistogramFacet) ValueField(valueField string) DateHistogramFacet {
	f.valueField = valueField
	return f
}

// Allowed values are: "year", "quarter", "month", "week", "day",
// "hour", "minute". It also supports time settings like "1.5h" 
// (up to "w" for weeks).
func (f DateHistogramFacet) Interval(interval string) DateHistogramFacet {
	f.interval = interval
	return f
}

func (f DateHistogramFacet) PostZone(postZone string) DateHistogramFacet {
	f.postZone = postZone
	return f
}

func (f DateHistogramFacet) PreZone(preZone string) DateHistogramFacet {
	f.preZone = preZone
	return f
}

func (f DateHistogramFacet) PostOffset(postOffset string) DateHistogramFacet {
	f.postOffset = postOffset
	return f
}

func (f DateHistogramFacet) PreOffset(preOffset string) DateHistogramFacet {
	f.preOffset = preOffset
	return f
}

func (f DateHistogramFacet) Factor(factor int) DateHistogramFacet {
	f.factor = &factor
	return f
}

func (f DateHistogramFacet) Source() interface{} {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["date_histogram"] = opts

	if f.keyField != "" {
		opts["key_field"] = f.keyField
		opts["value_field"] = f.valueField
	} else {
		opts["field"] = f.field
	}

	opts["interval"] = f.interval

	if f.postZone != "" {
		opts["post_zone"] = f.postZone
	}

	if f.preZone != "" {
		opts["pre_zone"] = f.preZone
	}

	if f.postOffset != "" {
		opts["post_offset"] = f.postOffset
	}

	if f.preOffset != "" {
		opts["pre_offset"] = f.preOffset
	}

	if f.factor != nil {
		opts["factor"] = *f.factor
	}

	if f.global != nil {
		opts["global"] = *f.global
	}

	return source
}
