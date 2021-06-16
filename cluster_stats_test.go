// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestClusterStatsURLs(t *testing.T) {
	fFlag := false
	tFlag := true

	tests := []struct {
		Service        *ClusterStatsService
		ExpectedPath   string
		ExpectedParams url.Values
	}{
		{
			Service: &ClusterStatsService{
				nodeId: []string{},
			},
			ExpectedPath: "/_cluster/stats",
		},
		{
			Service: &ClusterStatsService{
				nodeId: []string{"node1"},
			},
			ExpectedPath: "/_cluster/stats/nodes/node1",
		},
		{
			Service: &ClusterStatsService{
				nodeId: []string{"node1", "node2"},
			},
			ExpectedPath: "/_cluster/stats/nodes/node1%2Cnode2",
		},
		{
			Service: &ClusterStatsService{
				nodeId:       []string{},
				flatSettings: &tFlag,
			},
			ExpectedPath:   "/_cluster/stats",
			ExpectedParams: url.Values{"flat_settings": []string{"true"}},
		},
		{
			Service: &ClusterStatsService{
				nodeId:       []string{"node1"},
				flatSettings: &fFlag,
			},
			ExpectedPath:   "/_cluster/stats/nodes/node1",
			ExpectedParams: url.Values{"flat_settings": []string{"false"}},
		},
	}

	for _, test := range tests {
		gotPath, gotParams, err := test.Service.buildURL()
		if err != nil {
			t.Fatalf("expected no error; got: %v", err)
		}
		if gotPath != test.ExpectedPath {
			t.Errorf("expected URL path = %q; got: %q", test.ExpectedPath, gotPath)
		}
		if gotParams.Encode() != test.ExpectedParams.Encode() {
			t.Errorf("expected URL params = %v; got: %v", test.ExpectedParams, gotParams)
		}
	}
}

func TestClusterStatsErrorResponse(t *testing.T) {
	body := `{
		"_nodes": {
			"total": 2,
			"successful": 1,
			"failed": 1,
			"failures": [
				{
					"type": "failed_node_exception",
					"reason": "Failed node [mhUZF1sPTcu2b-pIJfqQRg]",
					"node_id": "mhUZF1sPTcu2b-pIJfqQRg",
					"caused_by": {
						"type": "node_not_connected_exception",
						"reason": "[es02][172.27.0.2:9300] Node not connected"
					}
				}
			]
		},
		"cluster_name": "es-docker-cluster",
		"cluster_uuid": "r-OkEGlJTFOE8wP36G-VSg",
		"timestamp": 1621834319499,
		"indices": {
			"count": 0,
			"shards": {},
			"docs": {
				"count": 0,
				"deleted": 0
			},
			"store": {
				"size_in_bytes": 0,
				"reserved_in_bytes": 0
			},
			"fielddata": {
				"memory_size_in_bytes": 0,
				"evictions": 0
			},
			"query_cache": {
				"memory_size_in_bytes": 0,
				"total_count": 0,
				"hit_count": 0,
				"miss_count": 0,
				"cache_size": 0,
				"cache_count": 0,
				"evictions": 0
			},
			"completion": {
				"size_in_bytes": 0
			},
			"segments": {
				"count": 0,
				"memory_in_bytes": 0,
				"terms_memory_in_bytes": 0,
				"stored_fields_memory_in_bytes": 0,
				"term_vectors_memory_in_bytes": 0,
				"norms_memory_in_bytes": 0,
				"points_memory_in_bytes": 0,
				"doc_values_memory_in_bytes": 0,
				"index_writer_memory_in_bytes": 0,
				"version_map_memory_in_bytes": 0,
				"fixed_bit_set_memory_in_bytes": 0,
				"max_unsafe_auto_id_timestamp": -9223372036854775808,
				"file_sizes": {}
			},
			"mappings": {
				"field_types": []
			},
			"analysis": {
				"char_filter_types": [],
				"tokenizer_types": [],
				"filter_types": [],
				"analyzer_types": [],
				"built_in_char_filters": [],
				"built_in_tokenizers": [],
				"built_in_filters": [],
				"built_in_analyzers": []
			}
		},
		"nodes": {
			"count": {
				"total": 1,
				"coordinating_only": 0,
				"data": 1,
				"data_cold": 1,
				"data_content": 1,
				"data_hot": 1,
				"data_warm": 1,
				"ingest": 1,
				"master": 1,
				"ml": 1,
				"remote_cluster_client": 1,
				"transform": 1,
				"voting_only": 0
			},
			"versions": [
				"7.10.0"
			],
			"os": {
				"available_processors": 6,
				"allocated_processors": 6,
				"names": [
					{
						"name": "Linux",
						"count": 1
					}
				],
				"pretty_names": [
					{
						"pretty_name": "CentOS Linux 8 (Core)",
						"count": 1
					}
				],
				"mem": {
					"total_in_bytes": 2084679680,
					"free_in_bytes": 590282752,
					"used_in_bytes": 1494396928,
					"free_percent": 28,
					"used_percent": 72
				}
			},
			"process": {
				"cpu": {
					"percent": 0
				},
				"open_file_descriptors": {
					"min": 260,
					"max": 260,
					"avg": 260
				}
			},
			"jvm": {
				"max_uptime_in_millis": 1042623,
				"versions": [
					{
						"version": "15.0.1",
						"vm_name": "OpenJDK 64-Bit Server VM",
						"vm_version": "15.0.1+9",
						"vm_vendor": "AdoptOpenJDK",
						"bundled_jdk": true,
						"using_bundled_jdk": true,
						"count": 1
					}
				],
				"mem": {
					"heap_used_in_bytes": 208299344,
					"heap_max_in_bytes": 314572800
				},
				"threads": 32
			},
			"fs": {
				"total_in_bytes": 62725623808,
				"free_in_bytes": 26955173888,
				"available_in_bytes": 23738458112
			},
			"plugins": [],
			"network_types": {
				"transport_types": {
					"security4": 1
				},
				"http_types": {
					"security4": 1
				}
			},
			"discovery_types": {
				"zen": 1
			},
			"packaging_types": [
				{
					"flavor": "default",
					"type": "docker",
					"count": 1
				}
			],
			"ingest": {
				"number_of_pipelines": 0,
				"processor_stats": {}
			}
		}
	}`

	var resp ClusterStatsResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatal(err)
	}
	if want, have := 1, len(resp.NodesStats.Failures); want != have {
		t.Fatalf("expected %d errors, got %d", want, have)
	}
}
