// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"errors"
)

// Reindexer simplifies the process of reindexing an index. You typically
// reindex a source index to a target index. However, you can also specify
// a query that filters out documents from the source index before bulk
// indexing them into the target index. The caller may also specify a
// different client for the target, e.g. when copying indices from one
// Elasticsearch cluster to another.
//
// Internally, the Reindex users a scan and scroll operation on the source
// index and bulk indexing to push data into the target index.
//
// The caller is responsible for setting up and/or clearing the target index
// before starting the reindex process.
//
// See http://www.elastic.co/guide/en/elasticsearch/guide/current/reindex.html
// for more information about reindexing.
type Reindexer struct {
	sourceClient, targetClient *Client
	sourceIndex, targetIndex   string
	query                      Query
	bulkSize                   int
	scroll                     string
	progress                   ReindexerProgressFunc
	statsOnly                  bool
}

// ReindexerProgressFunc is a callback that can be used with Reindexer
// to report progress while reindexing data.
type ReindexerProgressFunc func(current, total int64)

// ReindexerResponse is returned from the Do func in a Reindexer.
// By default, it returns the number of succeeded and failed bulk operations.
// To return details about all failed items, set StatsOnly to false in
// Reindexer.
type ReindexerResponse struct {
	Success int64
	Failed  int64
	Errors  []*BulkResponseItem
}

// NewReindexer returns a new Reindexer.
func NewReindexer(client *Client, source, target string) *Reindexer {
	return &Reindexer{
		sourceClient: client,
		sourceIndex:  source,
		targetIndex:  target,
		statsOnly:    true,
	}
}

// TargetClient specifies a different client for the target. This is
// necessary when the target index is in a different Elasticsearch cluster.
// By default, the source and target clients are the same.
func (ix *Reindexer) TargetClient(c *Client) *Reindexer {
	ix.targetClient = c
	return ix
}

// Query specifies the query to apply to the source. It filters out those
// documents to be indexed into target. A nil query does not filter out any
// documents.
func (ix *Reindexer) Query(q Query) *Reindexer {
	ix.query = q
	return ix
}

// BulkSize returns the number of documents to send to Elasticsearch per chunk.
// The default is 500.
func (ix *Reindexer) BulkSize(size int) *Reindexer {
	ix.bulkSize = size
	return ix
}

// Scroll specifies for how long the scroll operation on the source index
// should be maintained. The default is 5m.
func (ix *Reindexer) Scroll(timeout string) *Reindexer {
	ix.scroll = timeout
	return ix
}

// Progress indicates a callback that will be called while indexing.
func (ix *Reindexer) Progress(f ReindexerProgressFunc) *Reindexer {
	ix.progress = f
	return ix
}

// StatsOnly indicates whether the Do method should return details e.g. about
// the documents that failed while indexing. It is true by default, i.e. only
// the number of documents that succeeded/failed are returned. Set to false
// if you want all the details.
func (ix *Reindexer) StatsOnly(statsOnly bool) *Reindexer {
	ix.statsOnly = statsOnly
	return ix
}

// Do starts the reindexing process.
func (ix *Reindexer) Do() (*ReindexerResponse, error) {
	if ix.sourceClient == nil {
		return nil, errors.New("no source client")
	}
	if ix.sourceIndex == "" {
		return nil, errors.New("no source index")
	}
	if ix.targetIndex == "" {
		return nil, errors.New("no target index")
	}
	if ix.targetClient == nil {
		ix.targetClient = ix.sourceClient
	}
	if ix.bulkSize <= 0 {
		ix.bulkSize = 500
	}
	if ix.scroll == "" {
		ix.scroll = "5m"
	}

	// Count total to report progress (if necessary)
	var err error
	var current, total int64
	if ix.progress != nil {
		total, err = ix.count()
		if err != nil {
			return nil, err
		}
	}

	// Prepare scan and scroll to iterate through the source index
	scanner := ix.sourceClient.Scan(ix.sourceIndex).Scroll(ix.scroll)
	if ix.query != nil {
		scanner = scanner.Query(ix.query)
	}
	cursor, err := scanner.Do()

	bulk := ix.targetClient.Bulk().Index(ix.targetIndex)

	ret := &ReindexerResponse{
		Errors: make([]*BulkResponseItem, 0),
	}

	// Main loop iterates through the source index and bulk indexes into target.
	for {
		docs, err := cursor.Next()
		if err == EOS {
			break
		}
		if err != nil {
			return ret, err
		}

		if docs.TotalHits() > 0 {
			for _, hit := range docs.Hits.Hits {
				if ix.progress != nil {
					current++
					ix.progress(current, total)
				}

				// TODO(oe) Do we need to deserialize here?
				source := make(map[string]interface{})
				if err := json.Unmarshal(*hit.Source, &source); err != nil {
					return ret, err
				}

				// Enqueue and write into target index
				req := NewBulkIndexRequest().Index(ix.targetIndex).Type(hit.Type).Id(hit.Id).Doc(source)
				bulk.Add(req)
				if bulk.NumberOfActions() >= ix.bulkSize {
					bulk, err = ix.commit(bulk, ret)
					if err != nil {
						return ret, err
					}
				}
			}
		}
	}

	// Final flush
	if bulk.NumberOfActions() > 0 {
		bulk, err = ix.commit(bulk, ret)
		if err != nil {
			return ret, err
		}
		bulk = nil
	}

	return ret, nil
}

// count returns the number of documents in the source index.
// The query is taken into account, if specified.
func (ix *Reindexer) count() (int64, error) {
	service := ix.sourceClient.Count(ix.sourceIndex)
	if ix.query != nil {
		service = service.Query(ix.query)
	}
	return service.Do()
}

// commit commits a bulk, updates the stats, and returns a fresh bulk service.
func (ix *Reindexer) commit(bulk *BulkService, ret *ReindexerResponse) (*BulkService, error) {
	bres, err := bulk.Do()
	if err != nil {
		return nil, err
	}
	ret.Success += int64(len(bres.Succeeded()))
	failed := bres.Failed()
	ret.Failed += int64(len(failed))
	if !ix.statsOnly {
		ret.Errors = append(ret.Errors, failed...)
	}
	bulk = ix.targetClient.Bulk().Index(ix.targetIndex)
	return bulk, nil
}
