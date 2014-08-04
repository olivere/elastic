package elastic

import (
	"fmt"
)

// -- Bulkable request (index/update/delete) --

// Generic interface to bulkable requests.
type BulkableRequest interface {
	fmt.Stringer
	Source() ([]string, error)
}
