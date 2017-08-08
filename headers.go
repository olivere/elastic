// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

type headers map[string][]string

// addHeader adds key, value pair to the existing headers map, or creates a new map with that pair if headers was nil
func addHeader(headers headers, key string, value string) headers {
	if headers == nil {
		headers = make(map[string][]string)
	}

	var values []string
	if v, ok := headers[key]; ok {
		values = append(v, value)
	} else {
		values = []string{value}
	}

	headers[key] = values
	return headers
}
