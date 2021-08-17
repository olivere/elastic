// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// PointInTime is a lightweight view into the state of the data that existed
// when initiated. It can be created with OpenPointInTime API and be used
// when searching, e.g. in Search API or with SearchSource.
type PointInTime struct {
	// Id that uniquely identifies the point in time, as created with the
	// OpenPointInTime API.
	Id string `json:"id,omitempty"`
	// KeepAlive is the time for which this specific PointInTime will be
	// kept alive by Elasticsearch.
	KeepAlive string `json:"keep_alive,omitempty"`
}

// NewPointInTime creates a new PointInTime.
func NewPointInTime(id string, keepAlive ...string) *PointInTime {
	pit := &PointInTime{
		Id:        id,
	}
	if len(keepAlive) > 1 {
		pit.KeepAlive = keepAlive[0]
	}
	return pit
}

// Source generates the JSON serializable fragment for the PointInTime.
func (pit *PointInTime) Source() (interface{}, error) {
	if pit == nil {
		return nil, nil
	}
	source := map[string]interface{}{
		"id": pit.Id,
	}

	if pit.KeepAlive != "" {
		source["keep_alive"] = pit.KeepAlive
	}

	return source, nil
}
