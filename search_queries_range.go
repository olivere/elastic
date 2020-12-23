// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

const (
	// (Default) Matches documents with a range field value that intersects the query’s range.
	RelationIntersects string = "INTERSECTS"
	// Matches documents with a range field value that entirely contains the query’s range.
	RelationContains string = "CONTAINS"
	// Matches documents with a range field value entirely within the query’s range.
	RelationWithin string = "WITHIN"
)

// RangeQuery matches documents with fields that have terms within a certain range.
//
// For details, see Elastic Documentation (6.8):
// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-range-query.html
type RangeQuery struct {
	name     string
	gt       interface{}
	gte      interface{}
	lt       interface{}
	lte      interface{}
	format   string
	relation string
	timeZone string
	boost    *float64
}

// NewRangeQuery creates and initializes a new RangeQuery.
func NewRangeQuery(name string) *RangeQuery {
	return &RangeQuery{name: name}
}

// Gt (Optional) Greater than.
func (q *RangeQuery) Gt(gt interface{}) *RangeQuery {
	q.gt = gt
	return q
}

// Gte (Optional) Greater than or equal to.
func (q *RangeQuery) Gte(gte interface{}) *RangeQuery {
	q.gte = gte
	return q
}

// Lt (Optional) Less than.
func (q *RangeQuery) Lt(lt interface{}) *RangeQuery {
	q.lt = lt
	return q
}

// Lte (Optional) Less than or equal to.
func (q *RangeQuery) Lte(lte interface{}) *RangeQuery {
	q.lte = lte
	return q
}

// Boost sets the boost for this query.
func (q *RangeQuery) Boost(boost float64) *RangeQuery {
	q.boost = &boost
	return q
}

// TimeZone (Optional, string) Coordinated Universal Time (UTC) offset or IANA time zone used to convert date values in the query to UTC.
func (q *RangeQuery) TimeZone(timeZone string) *RangeQuery {
	q.timeZone = timeZone
	return q
}

// Format (Optional, string) Date format used to convert date values in the query.
func (q *RangeQuery) Format(format string) *RangeQuery {
	q.format = format
	return q
}

// Relation (Optional, string) Indicates how the range query matches values for range fields. Valid values are: INTERSECTS (Default), CONTAINS and WITHIN
func (q *RangeQuery) Relation(relation string) *RangeQuery {
	q.relation = relation
	return q
}

// Source returns JSON for the query.
func (q *RangeQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})

	rangeQ := make(map[string]interface{})
	source["range"] = rangeQ

	params := make(map[string]interface{})
	rangeQ[q.name] = params

	if q.gt != nil {
		params["gt"] = q.gt
	}
	if q.gte != nil {
		params["gte"] = q.gte
	}
	if q.lt != nil {
		params["lt"] = q.lt
	}
	if q.lte != nil {
		params["lte"] = q.lte
	}
	if q.timeZone != "" {
		params["time_zone"] = q.timeZone
	}
	if q.format != "" {
		params["format"] = q.format
	}
	if q.relation != "" {
		params["relation"] = q.relation
	}
	if q.boost != nil {
		params["boost"] = *q.boost
	}

	return source, nil
}
