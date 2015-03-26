// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

/*
Package elastic provides an interface to the Elasticsearch server
(http://www.elasticsearch.org/).

Notice: This is version 1 of Elastic. There are newer versions of Elastic
available on GitHub at https://github.com/olivere/elastic. Version 1 is
maintained, but new development happens in newer versions.

The first thing you do is to create a Client. The client takes a http.Client
and (optionally) a list of URLs to the Elasticsearch servers as arguments.
If the list of URLs is empty, http://localhost:9200 is used by default.
You typically create one client for your app.

	client, err := elastic.NewClient(http.DefaultClient)
	if err != nil {
		// Handle error
	}

Notice that you can pass your own http.Client implementation here. You can
also pass more than one URL to a client. Elastic pings the URLs periodically
and takes the first to succeed. By doing this periodically, Elastic provides
automatic failover, e.g. when an Elasticsearch server goes down during
updates.

If no Elasticsearch server is available, services will fail when creating
a new request and will return ErrNoClient. While this method is not very
sophisticated and might result in timeouts, it is robust enough for our
use cases. Pull requests are welcome.

	client, err := elastic.NewClient(http.DefaultClient, "http://1.2.3.4:9200", "http://1.2.3.5:9200")
	if err != nil {
		// Handle error
	}

A Client provides services. The services usually come with a variety of
methods to prepare the query and a Do function to execute it against the
Elasticsearch REST interface and return a response. Here is an example
of the IndexExists service that checks if a given index already exists.

	exists, err := client.IndexExists("twitter").Do()
	if err != nil {
		// Handle error
	}
	if !exists {
		// Index does not exist yet.
	}

Look up the documentation for Client to get an idea of the services provided
and what kinds of responses you get when executing the Do function of a service.
*/
package elastic
