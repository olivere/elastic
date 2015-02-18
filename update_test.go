// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestUpdateViaScript(t *testing.T) {
	client := setupTestClient(t)
	update := client.Update().
		Index("test").Type("type1").Id("1").
		Script("ctx._source.tags += tag").
		ScriptParams(map[string]interface{}{"tag": "blue"}).
		ScriptLang("groovy")
	urls, err := update.url()
	if err != nil {
		t.Fatalf("expected to return URL, got: %v", err)
	}
	expected := `/test/type1/1/_update`
	if expected != urls {
		t.Errorf("expected URL\n%s\ngot:\n%s", expected, urls)
	}
	body, err := update.body()
	if err != nil {
		t.Fatalf("expected to return body, got: %v", err)
	}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("expected to marshal body as JSON, got: %v", err)
	}
	got := string(data)
	expected = `{"lang":"groovy","params":{"tag":"blue"},"script":"ctx._source.tags += tag"}`
	if got != expected {
		t.Errorf("expected\n%s\ngot:\n%s", expected, got)
	}
}

func TestUpdateViaScriptId(t *testing.T) {
	client := setupTestClient(t)

	scriptParams := map[string]interface{}{
		"pageViewEvent": map[string]interface{}{
			"url":      "foo.com/bar",
			"response": 404,
			"time":     "2014-01-01 12:32",
		},
	}

	update := client.Update().
		Index("sessions").Type("session").Id("dh3sgudg8gsrgl").
		ScriptId("my_web_session_summariser").
		ScriptedUpsert(true).
		ScriptParams(scriptParams).
		Upsert(map[string]interface{}{})
	urls, err := update.url()
	if err != nil {
		t.Fatalf("expected to return URL, got: %v", err)
	}
	expected := `/sessions/session/dh3sgudg8gsrgl/_update`
	if expected != urls {
		t.Errorf("expected URL\n%s\ngot:\n%s", expected, urls)
	}
	body, err := update.body()
	if err != nil {
		t.Fatalf("expected to return body, got: %v", err)
	}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("expected to marshal body as JSON, got: %v", err)
	}
	got := string(data)
	expected = `{"params":{"pageViewEvent":{"response":404,"time":"2014-01-01 12:32","url":"foo.com/bar"}},"script_id":"my_web_session_summariser","scripted_upsert":true,"upsert":{}}`
	if got != expected {
		t.Errorf("expected\n%s\ngot:\n%s", expected, got)
	}
}

func TestUpdateViaScriptAndUpsert(t *testing.T) {
	client := setupTestClient(t)
	update := client.Update().
		Index("test").Type("type1").Id("1").
		Script("ctx._source.counter += count").
		ScriptParams(map[string]interface{}{"count": 4}).
		Upsert(map[string]interface{}{"counter": 1})
	urls, err := update.url()
	if err != nil {
		t.Fatalf("expected to return URL, got: %v", err)
	}
	expected := `/test/type1/1/_update`
	if expected != urls {
		t.Errorf("expected URL\n%s\ngot:\n%s", expected, urls)
	}
	body, err := update.body()
	if err != nil {
		t.Fatalf("expected to return body, got: %v", err)
	}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("expected to marshal body as JSON, got: %v", err)
	}
	got := string(data)
	expected = `{"params":{"count":4},"script":"ctx._source.counter += count","upsert":{"counter":1}}`
	if got != expected {
		t.Errorf("expected\n%s\ngot:\n%s", expected, got)
	}
}

func TestUpdateViaDoc(t *testing.T) {
	client := setupTestClient(t)
	update := client.Update().
		Index("test").Type("type1").Id("1").
		Doc(map[string]interface{}{"name": "new_name"}).
		DetectNoop(true)
	urls, err := update.url()
	if err != nil {
		t.Fatalf("expected to return URL, got: %v", err)
	}
	expected := `/test/type1/1/_update`
	if expected != urls {
		t.Errorf("expected URL\n%s\ngot:\n%s", expected, urls)
	}
	body, err := update.body()
	if err != nil {
		t.Fatalf("expected to return body, got: %v", err)
	}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("expected to marshal body as JSON, got: %v", err)
	}
	got := string(data)
	expected = `{"detect_noop":true,"doc":{"name":"new_name"}}`
	if got != expected {
		t.Errorf("expected\n%s\ngot:\n%s", expected, got)
	}
}

func TestUpdateViaDocAndUpsert(t *testing.T) {
	client := setupTestClient(t)
	update := client.Update().
		Index("test").Type("type1").Id("1").
		Doc(map[string]interface{}{"name": "new_name"}).
		DocAsUpsert(true).
		Timeout("1s").
		Refresh(true)
	urls, err := update.url()
	if err != nil {
		t.Fatalf("expected to return URL, got: %v", err)
	}
	expected := `/test/type1/1/_update?refresh=true&timeout=1s`
	if expected != urls {
		t.Errorf("expected URL\n%s\ngot:\n%s", expected, urls)
	}
	body, err := update.body()
	if err != nil {
		t.Fatalf("expected to return body, got: %v", err)
	}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("expected to marshal body as JSON, got: %v", err)
	}
	got := string(data)
	expected = `{"doc":{"name":"new_name"},"doc_as_upsert":true}`
	if got != expected {
		t.Errorf("expected\n%s\ngot:\n%s", expected, got)
	}
}

func TestUpdateViaScriptIntegration(t *testing.T) {
	client := setupTestClientAndCreateIndex(t)

	// Groovy is disabled for [1.3.8,1.4.0) and 1.4.3+, so skip tests in that case
	esversion, err := client.ElasticsearchVersion(defaultUrl)
	if err != nil {
		t.Fatal(err)
	}
	if esversion >= "1.4.3" || (esversion < "1.4.0" && esversion >= "1.3.8") {
		t.Skip("groovy scripting has been disabled as for [1.3.8,1.4.0) and 1.4.3+")
		return
	}

	tweet1 := tweet{User: "olivere", Retweets: 10, Message: "Welcome to Golang and Elasticsearch."}

	// Add a document
	indexResult, err := client.Index().
		Index(testIndexName).
		Type("tweet").
		Id("1").
		BodyJson(&tweet1).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if indexResult == nil {
		t.Errorf("expected result to be != nil; got: %v", indexResult)
	}

	// Update number of retweets
	increment := 1
	update, err := client.Update().Index(testIndexName).Type("tweet").Id("1").
		Script("ctx._source.retweets += num").
		ScriptParams(map[string]interface{}{"num": increment}).
		ScriptLang("groovy"). // Use "groovy" as default language as 1.3 uses MVEL by default
		// Pretty(true).Debug(true).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if update == nil {
		t.Errorf("expected update to be != nil; got %v", update)
	}
	if update.Version != indexResult.Version+1 {
		t.Errorf("expected version to be %d; got %d", indexResult.Version+1, update.Version)
	}

	// Get document
	getResult, err := client.Get().
		Index(testIndexName).
		Type("tweet").
		Id("1").
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if getResult.Index != testIndexName {
		t.Errorf("expected GetResult.Index %q; got %q", testIndexName, getResult.Index)
	}
	if getResult.Type != "tweet" {
		t.Errorf("expected GetResult.Type %q; got %q", "tweet", getResult.Type)
	}
	if getResult.Id != "1" {
		t.Errorf("expected GetResult.Id %q; got %q", "1", getResult.Id)
	}
	if getResult.Source == nil {
		t.Errorf("expected GetResult.Source to be != nil; got nil")
	}

	// Decode the Source field
	var tweetGot tweet
	err = json.Unmarshal(*getResult.Source, &tweetGot)
	if err != nil {
		t.Fatal(err)
	}
	if tweetGot.Retweets != tweet1.Retweets+increment {
		t.Errorf("expected Tweet.Retweets to be %d; got %d", tweet1.Retweets+increment, tweetGot.Retweets)
	}
}
