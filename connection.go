// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// defaultHeartbeatDuration is the time a connection to ES is periodically checked.
	defaultHeartbeatDuration = 30 * time.Second
)

// Connection is a single connection to an Elasticsearch server.
type Connection struct {
	sync.Mutex
	c                 *http.Client
	url               string
	broken            bool          // broken indicates if this connection is currently broken
	heartbeatDuration time.Duration // heartbeatDuration is the duration in which the ES server is checked
}

// NewConnection returns a new connection to an Elasticsearch server.
func NewConnection(client *http.Client, url string) *Connection {
	conn := &Connection{c: client, url: url, broken: true, heartbeatDuration: defaultHeartbeatDuration}

	// Check the ES server right away
	conn.checkBroken()

	// Periodically update the broken flag
	go conn.heartbeat()

	return conn
}

// heartbeat periodically checks if the connection to the Elasticsearch server is broken.
func (c *Connection) heartbeat() {
	ticker := time.NewTicker(c.heartbeatDuration)
	for {
		select {
		case <-ticker.C:
			c.checkBroken()
		}
	}
}

// checkBroken performs a request to check if the Elasticsearch server is broken.
// The broken flag of the Connection is updated here.
func (c *Connection) checkBroken() {
	params := make(url.Values)
	params.Set("timeout", "1")
	req, err := NewRequest("HEAD", c.url+"/?"+params.Encode())
	if err == nil {
		res, err := c.c.Do((*http.Request)(req))
		if err == nil {
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				// Not broken.
				c.Lock()
				c.broken = false
				c.Unlock()
			} else {
				// Broken.
				c.Lock()
				c.broken = true
				c.Unlock()
			}
		} else {
			// Request error: mark as broken.
			c.Lock()
			c.broken = true
			c.Unlock()
			log.Printf("elastic: %v", err)
		}
	} else {
		// Request error: mark as broken.
		c.Lock()
		c.broken = true
		c.Unlock()
		log.Printf("elastic: %v", err)
	}
}

// IsBroken returns true if the connection to the Elasticsearch server
// is broken, e.g. because the server is gone or unreachable.
func (c *Connection) IsBroken() bool {
	return c.broken
}
