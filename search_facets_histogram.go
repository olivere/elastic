// Copyright 2012 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// Histogram Facet
// See: http://www.elasticsearch.org/guide/reference/api/search/facets/histogram-facet.html
type HistogramFacet struct {
	Facet
	global       *bool
	field        string
	keyField     string
	valueField   string
	interval     *interface{}
	timeInterval string
}

func NewHistogramFacet(field string) HistogramFacet {
	return HistogramFacet{
		field:  field,
	}
}

func (f HistogramFacet) Global(global bool) HistogramFacet {
	f.global = &global
	return f
}

func (f HistogramFacet) Interval(interval interface{}) HistogramFacet {
	f.interval = &interval
	return f
}

func (f HistogramFacet) TimeInterval(timeInterval string) HistogramFacet {
	f.timeInterval = timeInterval
	return f
}

func (f HistogramFacet) KeyField(keyField string) HistogramFacet {
	f.keyField = keyField
	return f
}

func (f HistogramFacet) ValueField(valueField string) HistogramFacet {
	f.valueField = valueField
	return f
}

func (f HistogramFacet) Source() interface{} {
	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["histogram"] = opts

	if f.keyField != "" {
		opts["key_field"] = f.keyField
		opts["value_field"] = f.valueField
	} else {
		opts["field"] = f.field
	}

	if f.interval != nil {
		opts["interval"] = *f.interval
	}

	if f.timeInterval != "" {
		opts["time_interval"] = f.timeInterval
	}

	if f.global != nil {
		opts["global"] = *f.global
	}

	return source
}
