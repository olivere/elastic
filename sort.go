package elastic

// SortInfo contains information about sorting a field.
type SortInfo struct {
	Field          string
	Ascending      bool
	Missing        interface{}
	IgnoreUnmapped *bool
	SortMode       string
	NestedFilter   Filter
	NestedPath     string
}

func (info SortInfo) Source() interface{} {
	prop := make(map[string]interface{})
	if info.Ascending {
		prop["order"] = "asc"
	} else {
		prop["order"] = "desc"
	}
	if info.Missing != nil {
		prop["missing"] = info.Missing
	}
	if info.IgnoreUnmapped != nil {
		prop["ignore_unmapped"] = *info.IgnoreUnmapped
	}
	if info.SortMode != "" {
		prop["sort_mode"] = info.SortMode
	}
	if info.NestedFilter != nil {
		prop["nested_filter"] = info.NestedFilter
	}
	if info.NestedPath != "" {
		prop["nested_path"] = info.NestedPath
	}
	source := make(map[string]interface{})
	source[info.Field] = prop
	return source
}
