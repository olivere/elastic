package elastic

// Highlight allows highlighting search results on one or more fields.
// For details, see:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-request-highlighting.html
type Highlight struct {
	fields            []HighlighterField
	tagsSchema        string
	preTags           []string
	postTags          []string
	order             string
	encoder           string
	requireFieldMatch *bool
	highlighterType   string
	fragmenter        string
	highlightQuery    Query
	noMatchSize       *int
	options           map[string]interface{}
	forceSource       *bool
}

func NewHighlight() Highlight {
	hl := Highlight{
		fields:   make([]HighlighterField, 0),
		preTags:  make([]string, 0),
		postTags: make([]string, 0),
		options:  make(map[string]interface{}),
	}
	return hl
}

func (hl Highlight) Fields(fields ...HighlighterField) Highlight {
	hl.fields = append(hl.fields, fields...)
	return hl
}

func (hl Highlight) Field(name string, fragmentSize int) Highlight {
	field := NewHighlighterField(name).FragmentSize(fragmentSize)
	hl.fields = append(hl.fields, field)
	return hl
}

func (hl Highlight) TagsSchema(tagsSchema string) Highlight {
	hl.tagsSchema = tagsSchema
	return hl
}

func (hl Highlight) Encoder(encoder string) Highlight {
	hl.encoder = encoder
	return hl
}

func (hl Highlight) PreTags(preTags ...string) Highlight {
	hl.preTags = append(hl.preTags, preTags...)
	return hl
}

func (hl Highlight) PostTags(postTags ...string) Highlight {
	hl.postTags = append(hl.postTags, postTags...)
	return hl
}

func (hl Highlight) Order(order string) Highlight {
	hl.order = order
	return hl
}

func (hl Highlight) HighlighterType(typ string) Highlight {
	hl.highlighterType = typ
	return hl
}

func (hl Highlight) Fragmenter(fragmenter string) Highlight {
	hl.fragmenter = fragmenter
	return hl
}

func (hl Highlight) HighlighQuery(query Query) Highlight {
	hl.highlightQuery = query
	return hl
}

func (hl Highlight) NoMatchSize(size int) Highlight {
	hl.noMatchSize = &size
	return hl
}

func (hl Highlight) Options(options map[string]interface{}) Highlight {
	hl.options = options
	return hl
}

func (hl Highlight) ForceSource(forceSource bool) Highlight {
	hl.forceSource = &forceSource
	return hl
}

// Creates the query source for the bool query.
func (hl Highlight) Source() interface{} {
	highlightS := make(map[string]interface{})

	if hl.tagsSchema != "" {
		highlightS["tags_schema"] = hl.tagsSchema
	}
	if len(hl.preTags) > 0 {
		highlightS["pre_tags"] = hl.preTags
	}
	if len(hl.postTags) > 0 {
		highlightS["post_tags"] = hl.postTags
	}
	if hl.order != "" {
		highlightS["order"] = hl.order
	}
	if hl.encoder != "" {
		highlightS["encoder"] = hl.encoder
	}
	if hl.requireFieldMatch != nil {
		highlightS["require_field_match"] = *hl.requireFieldMatch
	}
	if hl.highlighterType != "" {
		highlightS["type"] = hl.highlighterType
	}
	if hl.fragmenter != "" {
		highlightS["fragmenter"] = hl.fragmenter
	}
	if hl.highlightQuery != nil {
		highlightS["highlight_query"] = hl.highlightQuery.Source()
	}
	if hl.noMatchSize != nil {
		highlightS["no_match_size"] = *hl.noMatchSize
	}
	if len(hl.options) > 0 {
		highlightS["options"] = hl.options
	}
	if hl.forceSource != nil {
		highlightS["force_source"] = *hl.forceSource
	}
	if len(hl.fields) > 0 {
		fieldsS := make(map[string]interface{})
		for _, field := range hl.fields {
			fieldsS[field.Name] = field.Source()
		}
		highlightS["fields"] = fieldsS
	}

	return highlightS
}

// HighlighterField specifies a highlighted field.
type HighlighterField struct {
	Name              string
	preTags           []string
	postTags          []string
	fragmentSize      int
	numOfFragments    int
	fragmentOffset    int
	highlightFilter   *bool
	order             string
	requireFieldMatch *bool
	boundaryMaxScan   int
	boundaryChars     []rune
	highlighterType   string
	fragmenter        string
	highlightQuery    Query
	noMatchSize       *int
	matchedFields     []string
	options           map[string]interface{}
	forceSource       *bool
}

func NewHighlighterField(name string) HighlighterField {
	return HighlighterField{
		Name:            name,
		fragmentSize:    -1,
		numOfFragments:  -1,
		fragmentOffset:  -1,
		boundaryMaxScan: -1,
		preTags:         make([]string, 0),
		postTags:        make([]string, 0),
		matchedFields:   make([]string, 0),
		boundaryChars:   make([]rune, 0),
		options:         make(map[string]interface{}),
	}
}

func (f HighlighterField) PreTags(preTags ...string) HighlighterField {
	f.preTags = append(f.preTags, preTags...)
	return f
}

func (f HighlighterField) PostTags(postTags ...string) HighlighterField {
	f.postTags = append(f.postTags, postTags...)
	return f
}

func (f HighlighterField) FragmentSize(size int) HighlighterField {
	f.fragmentSize = size
	return f
}

func (f HighlighterField) FragmentOffset(offset int) HighlighterField {
	f.fragmentOffset = offset
	return f
}

func (f HighlighterField) NumOfFragments(num int) HighlighterField {
	f.numOfFragments = num
	return f
}

func (f HighlighterField) HighlightFilter(highlightFilter bool) HighlighterField {
	f.highlightFilter = &highlightFilter
	return f
}

func (f HighlighterField) Order(order string) HighlighterField {
	f.order = order
	return f
}

func (f HighlighterField) RequireFieldMatch(require bool) HighlighterField {
	f.requireFieldMatch = &require
	return f
}

func (f HighlighterField) BoundaryMaxScan(maxScan int) HighlighterField {
	f.boundaryMaxScan = maxScan
	return f
}

func (f HighlighterField) BoundaryChars(chars []rune) HighlighterField {
	f.boundaryChars = chars
	return f
}

func (f HighlighterField) HighlighterType(typ string) HighlighterField {
	f.highlighterType = typ
	return f
}

func (f HighlighterField) Fragmenter(fragmenter string) HighlighterField {
	f.fragmenter = fragmenter
	return f
}

func (f HighlighterField) HighlightQuery(query Query) HighlighterField {
	f.highlightQuery = query
	return f
}

func (f HighlighterField) NoMatchSize(size int) HighlighterField {
	f.noMatchSize = &size
	return f
}

func (f HighlighterField) MatchedFields(fields ...string) HighlighterField {
	f.matchedFields = append(f.matchedFields, fields...)
	return f
}

func (f HighlighterField) Options(options map[string]interface{}) HighlighterField {
	f.options = options
	return f
}

func (f HighlighterField) ForceSource(forceSource bool) HighlighterField {
	f.forceSource = &forceSource
	return f
}

func (f HighlighterField) Source() interface{} {
	source := make(map[string]interface{})
	if len(f.preTags) > 0 {
		source["pre_tags"] = f.preTags
	}
	if len(f.postTags) > 0 {
		source["post_tags"] = f.postTags
	}
	if f.fragmentSize != -1 {
		source["fragment_size"] = f.fragmentSize
	}
	if f.numOfFragments != -1 {
		source["number_of_fragments"] = f.numOfFragments
	}
	if f.fragmentOffset != -1 {
		source["fragment_offset"] = f.fragmentOffset
	}
	if f.order != "" {
		source["order"] = f.order
	}
	if f.requireFieldMatch != nil {
		source["require_field_match"] = *f.requireFieldMatch
	}
	if f.boundaryMaxScan != -1 {
		source["boundary_max_scan"] = f.boundaryMaxScan
	}
	if len(f.boundaryChars) > 0 {
		source["boundary_chars"] = f.boundaryChars
	}
	if f.highlighterType != "" {
		source["type"] = f.highlighterType
	}
	if f.fragmenter != "" {
		source["fragmenter"] = f.fragmenter
	}
	if f.highlightQuery != nil {
		source["highlight_query"] = f.highlightQuery.Source()
	}
	if f.noMatchSize != nil {
		source["no_match_size"] = *f.noMatchSize
	}
	if len(f.matchedFields) > 0 {
		source["matched_fields"] = f.matchedFields
	}
	if len(f.options) > 0 {
		source["options"] = f.options
	}
	if f.forceSource != nil {
		source["force_source"] = *f.forceSource
	}
	return source
}
