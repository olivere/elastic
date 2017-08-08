// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// addHeader adds key, value pair to the existing headers map, or creates a new map with that pair if headers was nil
func addHeader(headers map[string][]string, key string, value string) map[string][]string {
	if headers == nil {
		headers = make(map[string][]string)
	}

	var values []string
	if v, ok := headers[key]; ok {
		values = v
	} else {
		values = make([]string, 0)
	}

	headers[key] = append(values, value)
	return headers
}
