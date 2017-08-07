package elastic

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
