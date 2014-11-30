// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"net/http"
	"sync"
)

// ConnectionPool is a pool of connections. It is used by Client to
// get a connection for requests.
type ConnectionPool struct {
	sync.Mutex
	c     *http.Client
	conns []*Connection
	index int // index into conns that points to the connection to use for the next request
}

// NewConnectionPool creates a new pool of connections to Elasticsearch servers.
func NewConnectionPool(client *http.Client, urls ...string) *ConnectionPool {
	pool := &ConnectionPool{
		c:     client,
		conns: make([]*Connection, 0),
	}

	if len(urls) == 0 {
		pool.Add(NewConnection(client, defaultUrl))
	} else {
		for _, url := range urls {
			pool.Add(NewConnection(client, url))
		}
	}
	return pool
}

// Add adds a new connection to the pool.
func (pool *ConnectionPool) Add(conn *Connection) {
	pool.Lock()
	pool.conns = append(pool.conns, conn)
	pool.Unlock()
}

// GetNextRequestURL returns the URL of a connection from the pool.
// It currently uses simple round-robin.
func (pool *ConnectionPool) GetNextRequestURL() (string, error) {
	if len(pool.conns) == 0 {
		return "", ErrNoClient
	}

	pool.Lock()

	// simple round-robin to get a (non-broken) connection
	var conn *Connection
	for i := 0; i < len(pool.conns); i++ {
		// increment index for next request
		pool.index = (pool.index + 1) % len(pool.conns)
		// get candidate and check if broken
		candidate := pool.conns[pool.index]
		if !candidate.IsBroken() {
			// not broken, let's use it
			conn = candidate
			break
		}
	}

	pool.Unlock()

	if conn == nil {
		return "", ErrNoClient
	}

	return conn.url, nil
}
