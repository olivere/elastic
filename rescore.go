package elastic

type Rescore struct {
	rescorer                 Rescorer
	windowSize               *int
	defaultRescoreWindowSize *int
}

func NewRescore() *Rescore {
	return &Rescore{}
}

func (r *Rescore) WindowSize(windowSize int) *Rescore {
	r.windowSize = &windowSize
	return r
}

func (r *Rescore) IsEmpty() bool {
	return r.rescorer == nil
}

func (r *Rescore) Rescorer(rescorer Rescorer) *Rescore {
	r.rescorer = rescorer
	return r
}

func (r *Rescore) Source() interface{} {
	source := make(map[string]interface{})
	if r.windowSize != nil {
		source["window_size"] = *r.windowSize
	} else if r.defaultRescoreWindowSize != nil {
		source["window_size"] = *r.defaultRescoreWindowSize
	}
	source[r.rescorer.Name()] = r.rescorer.Source()
	return source
}
