package elastic

type GeohashGridAggregation struct {
	field           string
	precision       int
	subAggregations map[string]Aggregation
	meta            map[string]interface{}
}

func NewGeohashGridAggregation() *GeohashGridAggregation {
	return &GeoHashAggregation{
		subAggregations: make(map[string]Aggregation),
	}
}

func (a *GeohashGridAggregation) Field(field string) *GeohashGridAggregation {
	a.field = field
	return a
}

func (a *GeohashGridAggregation) Precision(precision int) 8GeohashGridAggregation {
    a.precision = precision
    return a
}

func (a *GeohashGridAggregation) SubAggregation(name string, subAggregation Aggregation) *GeohashGridAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

func (a *GeohashGridAggregation) Meta(metaData map[string]interface{}) *GeohashGridAggregation {
	a.meta = metaData
	return a
}

func (a *GeohashGridAggregation) Source() (interface{}, error) {
	// Example:
	// {
	//     "aggs": {
	//         "new_york": {
	//             "geohash_grid": {
	//                 "field": "location",
	//                 "precision": 5
	//             }
	//         }
	//     }
	// }

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["geohash_grid"] = opts

	if a.field != "" {
		opts["field"] = a.field
	}

    if a.precision != 0 {
        opts["precision"] = a.precision
    }

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

    return source, nil
}
