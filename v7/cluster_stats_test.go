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

// -- TestClusterStatsErrorResponse --

var clusterStatsErrorResponseTests = []struct {
	Body                       string
	ExpectedNodesStatsFailures int
}{
	// #0
	{
		Body: `{
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
		}`,
		ExpectedNodesStatsFailures: 1,
	},
	// #1 7.13.2 happy path
	{
		Body: `{
			"_nodes" : {
			  "total" : 3,
			  "successful" : 3,
			  "failed" : 0
			},
			"cluster_name" : "elasticsearch",
			"cluster_uuid" : "8TTeQMxRSZmffmYcTjP21w",
			"timestamp" : 1625645280402,
			"status" : "green",
			"indices" : {
			  "count" : 0,
			  "shards" : { },
			  "docs" : {
				"count" : 0,
				"deleted" : 0
			  },
			  "store" : {
				"size_in_bytes" : 0,
				"total_data_set_size_in_bytes" : 0,
				"reserved_in_bytes" : 0
			  },
			  "fielddata" : {
				"memory_size_in_bytes" : 0,
				"evictions" : 0
			  },
			  "query_cache" : {
				"memory_size_in_bytes" : 0,
				"total_count" : 0,
				"hit_count" : 0,
				"miss_count" : 0,
				"cache_size" : 0,
				"cache_count" : 0,
				"evictions" : 0
			  },
			  "completion" : {
				"size_in_bytes" : 0
			  },
			  "segments" : {
				"count" : 0,
				"memory_in_bytes" : 0,
				"terms_memory_in_bytes" : 0,
				"stored_fields_memory_in_bytes" : 0,
				"term_vectors_memory_in_bytes" : 0,
				"norms_memory_in_bytes" : 0,
				"points_memory_in_bytes" : 0,
				"doc_values_memory_in_bytes" : 0,
				"index_writer_memory_in_bytes" : 0,
				"version_map_memory_in_bytes" : 0,
				"fixed_bit_set_memory_in_bytes" : 0,
				"max_unsafe_auto_id_timestamp" : -9223372036854775808,
				"file_sizes" : { }
			  },
			  "mappings" : {
				"field_types" : [ ],
				"runtime_field_types" : [ ]
			  },
			  "analysis" : {
				"char_filter_types" : [ ],
				"tokenizer_types" : [ ],
				"filter_types" : [ ],
				"analyzer_types" : [ ],
				"built_in_char_filters" : [ ],
				"built_in_tokenizers" : [ ],
				"built_in_filters" : [ ],
				"built_in_analyzers" : [ ]
			  },
			  "versions" : [ ]
			},
			"nodes" : {
			  "count" : {
				"total" : 3,
				"coordinating_only" : 0,
				"data" : 3,
				"data_cold" : 3,
				"data_content" : 3,
				"data_frozen" : 3,
				"data_hot" : 3,
				"data_warm" : 3,
				"ingest" : 3,
				"master" : 3,
				"ml" : 3,
				"remote_cluster_client" : 3,
				"transform" : 3,
				"voting_only" : 0
			  },
			  "versions" : [
				"7.13.2"
			  ],
			  "os" : {
				"available_processors" : 12,
				"allocated_processors" : 12,
				"names" : [
				  {
					"name" : "Linux",
					"count" : 3
				  }
				],
				"pretty_names" : [
				  {
					"pretty_name" : "CentOS Linux 8",
					"count" : 3
				  }
				],
				"architectures" : [
				  {
					"arch" : "aarch64",
					"count" : 3
				  }
				],
				"mem" : {
				  "total_in_bytes" : 25013551104,
				  "free_in_bytes" : 250650624,
				  "used_in_bytes" : 24762900480,
				  "free_percent" : 1,
				  "used_percent" : 99
				}
			  },
			  "process" : {
				"cpu" : {
				  "percent" : 0
				},
				"open_file_descriptors" : {
				  "min" : 329,
				  "max" : 330,
				  "avg" : 329
				}
			  },
			  "jvm" : {
				"max_uptime_in_millis" : 119028,
				"versions" : [
				  {
					"version" : "16",
					"vm_name" : "OpenJDK 64-Bit Server VM",
					"vm_version" : "16+36",
					"vm_vendor" : "AdoptOpenJDK",
					"bundled_jdk" : true,
					"using_bundled_jdk" : true,
					"count" : 3
				  }
				],
				"mem" : {
				  "heap_used_in_bytes" : 748080592,
				  "heap_max_in_bytes" : 3221225472
				},
				"threads" : 87
			  },
			  "fs" : {
				"total_in_bytes" : 481133838336,
				"free_in_bytes" : 459472416768,
				"available_in_bytes" : 434940936192
			  },
			  "plugins" : [ ],
			  "network_types" : {
				"transport_types" : {
				  "netty4" : 3
				},
				"http_types" : {
				  "netty4" : 3
				}
			  },
			  "discovery_types" : {
				"zen" : 3
			  },
			  "packaging_types" : [
				{
				  "flavor" : "default",
				  "type" : "docker",
				  "count" : 3
				}
			  ],
			  "ingest" : {
				"number_of_pipelines" : 2,
				"processor_stats" : {
				  "gsub" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  },
				  "script" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  }
				}
			  }
			}
		  }`,
		ExpectedNodesStatsFailures: 0,
	},
	// #2 7.13.2 failed nodes
	{
		Body: `{
			"_nodes" : {
			  "total" : 2,
			  "successful" : 1,
			  "failed" : 1,
			  "failures" : [
				{
				  "type" : "failed_node_exception",
				  "reason" : "Failed node [agq-C6ZPSYmPgNePwb26_Q]",
				  "node_id" : "agq-C6ZPSYmPgNePwb26_Q",
				  "caused_by" : {
					"type" : "node_not_connected_exception",
					"reason" : "[es2][172.23.0.4:9300] Node not connected"
				  }
				}
			  ]
			},
			"cluster_name" : "elasticsearch",
			"cluster_uuid" : "8TTeQMxRSZmffmYcTjP21w",
			"timestamp" : 1625645352594,
			"indices" : {
			  "count" : 0,
			  "shards" : { },
			  "docs" : {
				"count" : 0,
				"deleted" : 0
			  },
			  "store" : {
				"size_in_bytes" : 0,
				"total_data_set_size_in_bytes" : 0,
				"reserved_in_bytes" : 0
			  },
			  "fielddata" : {
				"memory_size_in_bytes" : 0,
				"evictions" : 0
			  },
			  "query_cache" : {
				"memory_size_in_bytes" : 0,
				"total_count" : 0,
				"hit_count" : 0,
				"miss_count" : 0,
				"cache_size" : 0,
				"cache_count" : 0,
				"evictions" : 0
			  },
			  "completion" : {
				"size_in_bytes" : 0
			  },
			  "segments" : {
				"count" : 0,
				"memory_in_bytes" : 0,
				"terms_memory_in_bytes" : 0,
				"stored_fields_memory_in_bytes" : 0,
				"term_vectors_memory_in_bytes" : 0,
				"norms_memory_in_bytes" : 0,
				"points_memory_in_bytes" : 0,
				"doc_values_memory_in_bytes" : 0,
				"index_writer_memory_in_bytes" : 0,
				"version_map_memory_in_bytes" : 0,
				"fixed_bit_set_memory_in_bytes" : 0,
				"max_unsafe_auto_id_timestamp" : -9223372036854775808,
				"file_sizes" : { }
			  },
			  "mappings" : {
				"field_types" : [ ],
				"runtime_field_types" : [ ]
			  },
			  "analysis" : {
				"char_filter_types" : [ ],
				"tokenizer_types" : [ ],
				"filter_types" : [ ],
				"analyzer_types" : [ ],
				"built_in_char_filters" : [ ],
				"built_in_tokenizers" : [ ],
				"built_in_filters" : [ ],
				"built_in_analyzers" : [ ]
			  },
			  "versions" : [ ]
			},
			"nodes" : {
			  "count" : {
				"total" : 1,
				"coordinating_only" : 0,
				"data" : 1,
				"data_cold" : 1,
				"data_content" : 1,
				"data_frozen" : 1,
				"data_hot" : 1,
				"data_warm" : 1,
				"ingest" : 1,
				"master" : 1,
				"ml" : 1,
				"remote_cluster_client" : 1,
				"transform" : 1,
				"voting_only" : 0
			  },
			  "versions" : [
				"7.13.2"
			  ],
			  "os" : {
				"available_processors" : 4,
				"allocated_processors" : 4,
				"names" : [
				  {
					"name" : "Linux",
					"count" : 1
				  }
				],
				"pretty_names" : [
				  {
					"pretty_name" : "CentOS Linux 8",
					"count" : 1
				  }
				],
				"architectures" : [
				  {
					"arch" : "aarch64",
					"count" : 1
				  }
				],
				"mem" : {
				  "total_in_bytes" : 8337850368,
				  "free_in_bytes" : 2963259392,
				  "used_in_bytes" : 5374590976,
				  "free_percent" : 36,
				  "used_percent" : 64
				}
			  },
			  "process" : {
				"cpu" : {
				  "percent" : 0
				},
				"open_file_descriptors" : {
				  "min" : 279,
				  "max" : 279,
				  "avg" : 279
				}
			  },
			  "jvm" : {
				"max_uptime_in_millis" : 189845,
				"versions" : [
				  {
					"version" : "16",
					"vm_name" : "OpenJDK 64-Bit Server VM",
					"vm_version" : "16+36",
					"vm_vendor" : "AdoptOpenJDK",
					"bundled_jdk" : true,
					"using_bundled_jdk" : true,
					"count" : 1
				  }
				],
				"mem" : {
				  "heap_used_in_bytes" : 315085776,
				  "heap_max_in_bytes" : 1073741824
				},
				"threads" : 34
			  },
			  "fs" : {
				"total_in_bytes" : 160377946112,
				"free_in_bytes" : 153157689344,
				"available_in_bytes" : 144980529152
			  },
			  "plugins" : [ ],
			  "network_types" : {
				"transport_types" : {
				  "netty4" : 1
				},
				"http_types" : {
				  "netty4" : 1
				}
			  },
			  "discovery_types" : {
				"zen" : 1
			  },
			  "packaging_types" : [
				{
				  "flavor" : "default",
				  "type" : "docker",
				  "count" : 1
				}
			  ],
			  "ingest" : {
				"number_of_pipelines" : 1,
				"processor_stats" : {
				  "gsub" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  },
				  "script" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  }
				}
			  }
			}
		  }`,
		ExpectedNodesStatsFailures: 1,
	},
	// #3 7.12.1 happy path
	{
		Body: `{
			"_nodes" : {
			  "total" : 3,
			  "successful" : 3,
			  "failed" : 0
			},
			"cluster_name" : "elasticsearch",
			"cluster_uuid" : "XmKLDCnmRqqK3Ve01uJU0Q",
			"timestamp" : 1625645497836,
			"status" : "green",
			"indices" : {
			  "count" : 0,
			  "shards" : { },
			  "docs" : {
				"count" : 0,
				"deleted" : 0
			  },
			  "store" : {
				"size_in_bytes" : 0,
				"reserved_in_bytes" : 0
			  },
			  "fielddata" : {
				"memory_size_in_bytes" : 0,
				"evictions" : 0
			  },
			  "query_cache" : {
				"memory_size_in_bytes" : 0,
				"total_count" : 0,
				"hit_count" : 0,
				"miss_count" : 0,
				"cache_size" : 0,
				"cache_count" : 0,
				"evictions" : 0
			  },
			  "completion" : {
				"size_in_bytes" : 0
			  },
			  "segments" : {
				"count" : 0,
				"memory_in_bytes" : 0,
				"terms_memory_in_bytes" : 0,
				"stored_fields_memory_in_bytes" : 0,
				"term_vectors_memory_in_bytes" : 0,
				"norms_memory_in_bytes" : 0,
				"points_memory_in_bytes" : 0,
				"doc_values_memory_in_bytes" : 0,
				"index_writer_memory_in_bytes" : 0,
				"version_map_memory_in_bytes" : 0,
				"fixed_bit_set_memory_in_bytes" : 0,
				"max_unsafe_auto_id_timestamp" : -9223372036854775808,
				"file_sizes" : { }
			  },
			  "mappings" : {
				"field_types" : [ ]
			  },
			  "analysis" : {
				"char_filter_types" : [ ],
				"tokenizer_types" : [ ],
				"filter_types" : [ ],
				"analyzer_types" : [ ],
				"built_in_char_filters" : [ ],
				"built_in_tokenizers" : [ ],
				"built_in_filters" : [ ],
				"built_in_analyzers" : [ ]
			  },
			  "versions" : [ ]
			},
			"nodes" : {
			  "count" : {
				"total" : 3,
				"coordinating_only" : 0,
				"data" : 3,
				"data_cold" : 3,
				"data_content" : 3,
				"data_frozen" : 3,
				"data_hot" : 3,
				"data_warm" : 3,
				"ingest" : 3,
				"master" : 3,
				"ml" : 3,
				"remote_cluster_client" : 3,
				"transform" : 3,
				"voting_only" : 0
			  },
			  "versions" : [
				"7.12.1"
			  ],
			  "os" : {
				"available_processors" : 12,
				"allocated_processors" : 12,
				"names" : [
				  {
					"name" : "Linux",
					"count" : 3
				  }
				],
				"pretty_names" : [
				  {
					"pretty_name" : "CentOS Linux 8",
					"count" : 3
				  }
				],
				"architectures" : [
				  {
					"arch" : "aarch64",
					"count" : 3
				  }
				],
				"mem" : {
				  "total_in_bytes" : 25013551104,
				  "free_in_bytes" : 277807104,
				  "used_in_bytes" : 24735744000,
				  "free_percent" : 1,
				  "used_percent" : 99
				}
			  },
			  "process" : {
				"cpu" : {
				  "percent" : 1
				},
				"open_file_descriptors" : {
				  "min" : 332,
				  "max" : 333,
				  "avg" : 332
				}
			  },
			  "jvm" : {
				"max_uptime_in_millis" : 36043,
				"versions" : [
				  {
					"version" : "16",
					"vm_name" : "OpenJDK 64-Bit Server VM",
					"vm_version" : "16+36",
					"vm_vendor" : "AdoptOpenJDK",
					"bundled_jdk" : true,
					"using_bundled_jdk" : true,
					"count" : 3
				  }
				],
				"mem" : {
				  "heap_used_in_bytes" : 754473472,
				  "heap_max_in_bytes" : 3221225472
				},
				"threads" : 91
			  },
			  "fs" : {
				"total_in_bytes" : 481133838336,
				"free_in_bytes" : 457492058112,
				"available_in_bytes" : 432960577536
			  },
			  "plugins" : [ ],
			  "network_types" : {
				"transport_types" : {
				  "netty4" : 3
				},
				"http_types" : {
				  "netty4" : 3
				}
			  },
			  "discovery_types" : {
				"zen" : 3
			  },
			  "packaging_types" : [
				{
				  "flavor" : "default",
				  "type" : "docker",
				  "count" : 3
				}
			  ],
			  "ingest" : {
				"number_of_pipelines" : 2,
				"processor_stats" : {
				  "gsub" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  },
				  "script" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  }
				}
			  }
			}
		  }`,
		ExpectedNodesStatsFailures: 0,
	},
	// #4 7.12.1 failed nodes
	{
		Body: `{
			"_nodes" : {
			  "total" : 3,
			  "successful" : 1,
			  "failed" : 2,
			  "failures" : [
				{
				  "type" : "failed_node_exception",
				  "reason" : "Failed node [eE61wq63TxW0yN1ILeoZ5g]",
				  "node_id" : "eE61wq63TxW0yN1ILeoZ5g",
				  "caused_by" : {
					"type" : "node_not_connected_exception",
					"reason" : "[es3][172.24.0.4:9300] Node not connected"
				  }
				},
				{
				  "type" : "failed_node_exception",
				  "reason" : "Failed node [AAQlfCf0TJ644wJvjaOSyA]",
				  "node_id" : "AAQlfCf0TJ644wJvjaOSyA",
				  "caused_by" : {
					"type" : "node_not_connected_exception",
					"reason" : "[es2][172.24.0.3:9300] Node not connected"
				  }
				}
			  ]
			},
			"cluster_name" : "elasticsearch",
			"cluster_uuid" : "XmKLDCnmRqqK3Ve01uJU0Q",
			"timestamp" : 1625645535819,
			"indices" : {
			  "count" : 0,
			  "shards" : { },
			  "docs" : {
				"count" : 0,
				"deleted" : 0
			  },
			  "store" : {
				"size_in_bytes" : 0,
				"reserved_in_bytes" : 0
			  },
			  "fielddata" : {
				"memory_size_in_bytes" : 0,
				"evictions" : 0
			  },
			  "query_cache" : {
				"memory_size_in_bytes" : 0,
				"total_count" : 0,
				"hit_count" : 0,
				"miss_count" : 0,
				"cache_size" : 0,
				"cache_count" : 0,
				"evictions" : 0
			  },
			  "completion" : {
				"size_in_bytes" : 0
			  },
			  "segments" : {
				"count" : 0,
				"memory_in_bytes" : 0,
				"terms_memory_in_bytes" : 0,
				"stored_fields_memory_in_bytes" : 0,
				"term_vectors_memory_in_bytes" : 0,
				"norms_memory_in_bytes" : 0,
				"points_memory_in_bytes" : 0,
				"doc_values_memory_in_bytes" : 0,
				"index_writer_memory_in_bytes" : 0,
				"version_map_memory_in_bytes" : 0,
				"fixed_bit_set_memory_in_bytes" : 0,
				"max_unsafe_auto_id_timestamp" : -9223372036854775808,
				"file_sizes" : { }
			  },
			  "mappings" : {
				"field_types" : [ ]
			  },
			  "analysis" : {
				"char_filter_types" : [ ],
				"tokenizer_types" : [ ],
				"filter_types" : [ ],
				"analyzer_types" : [ ],
				"built_in_char_filters" : [ ],
				"built_in_tokenizers" : [ ],
				"built_in_filters" : [ ],
				"built_in_analyzers" : [ ]
			  },
			  "versions" : [ ]
			},
			"nodes" : {
			  "count" : {
				"total" : 1,
				"coordinating_only" : 0,
				"data" : 1,
				"data_cold" : 1,
				"data_content" : 1,
				"data_frozen" : 1,
				"data_hot" : 1,
				"data_warm" : 1,
				"ingest" : 1,
				"master" : 1,
				"ml" : 1,
				"remote_cluster_client" : 1,
				"transform" : 1,
				"voting_only" : 0
			  },
			  "versions" : [
				"7.12.1"
			  ],
			  "os" : {
				"available_processors" : 4,
				"allocated_processors" : 4,
				"names" : [
				  {
					"name" : "Linux",
					"count" : 1
				  }
				],
				"pretty_names" : [
				  {
					"pretty_name" : "CentOS Linux 8",
					"count" : 1
				  }
				],
				"architectures" : [
				  {
					"arch" : "aarch64",
					"count" : 1
				  }
				],
				"mem" : {
				  "total_in_bytes" : 8337850368,
				  "free_in_bytes" : 2912473088,
				  "used_in_bytes" : 5425377280,
				  "free_percent" : 35,
				  "used_percent" : 65
				}
			  },
			  "process" : {
				"cpu" : {
				  "percent" : 0
				},
				"open_file_descriptors" : {
				  "min" : 308,
				  "max" : 308,
				  "avg" : 308
				}
			  },
			  "jvm" : {
				"max_uptime_in_millis" : 74067,
				"versions" : [
				  {
					"version" : "16",
					"vm_name" : "OpenJDK 64-Bit Server VM",
					"vm_version" : "16+36",
					"vm_vendor" : "AdoptOpenJDK",
					"bundled_jdk" : true,
					"using_bundled_jdk" : true,
					"count" : 1
				  }
				],
				"mem" : {
				  "heap_used_in_bytes" : 219877888,
				  "heap_max_in_bytes" : 1073741824
				},
				"threads" : 32
			  },
			  "fs" : {
				"total_in_bytes" : 160377946112,
				"free_in_bytes" : 152497557504,
				"available_in_bytes" : 144320397312
			  },
			  "plugins" : [ ],
			  "network_types" : {
				"transport_types" : {
				  "netty4" : 1
				},
				"http_types" : {
				  "netty4" : 1
				}
			  },
			  "discovery_types" : {
				"zen" : 1
			  },
			  "packaging_types" : [
				{
				  "flavor" : "default",
				  "type" : "docker",
				  "count" : 1
				}
			  ],
			  "ingest" : {
				"number_of_pipelines" : 1,
				"processor_stats" : {
				  "gsub" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  },
				  "script" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  }
				}
			  }
			}
		  }`,
		ExpectedNodesStatsFailures: 2,
	},
	// #5 7.10.0 happy path
	{
		Body: `{
			"_nodes" : {
			  "total" : 3,
			  "successful" : 3,
			  "failed" : 0
			},
			"cluster_name" : "elasticsearch",
			"cluster_uuid" : "XvRExZcTTqiiYa8Qa0E6LA",
			"timestamp" : 1625645743243,
			"status" : "green",
			"indices" : {
			  "count" : 0,
			  "shards" : { },
			  "docs" : {
				"count" : 0,
				"deleted" : 0
			  },
			  "store" : {
				"size_in_bytes" : 0,
				"reserved_in_bytes" : 0
			  },
			  "fielddata" : {
				"memory_size_in_bytes" : 0,
				"evictions" : 0
			  },
			  "query_cache" : {
				"memory_size_in_bytes" : 0,
				"total_count" : 0,
				"hit_count" : 0,
				"miss_count" : 0,
				"cache_size" : 0,
				"cache_count" : 0,
				"evictions" : 0
			  },
			  "completion" : {
				"size_in_bytes" : 0
			  },
			  "segments" : {
				"count" : 0,
				"memory_in_bytes" : 0,
				"terms_memory_in_bytes" : 0,
				"stored_fields_memory_in_bytes" : 0,
				"term_vectors_memory_in_bytes" : 0,
				"norms_memory_in_bytes" : 0,
				"points_memory_in_bytes" : 0,
				"doc_values_memory_in_bytes" : 0,
				"index_writer_memory_in_bytes" : 0,
				"version_map_memory_in_bytes" : 0,
				"fixed_bit_set_memory_in_bytes" : 0,
				"max_unsafe_auto_id_timestamp" : -9223372036854775808,
				"file_sizes" : { }
			  },
			  "mappings" : {
				"field_types" : [ ]
			  },
			  "analysis" : {
				"char_filter_types" : [ ],
				"tokenizer_types" : [ ],
				"filter_types" : [ ],
				"analyzer_types" : [ ],
				"built_in_char_filters" : [ ],
				"built_in_tokenizers" : [ ],
				"built_in_filters" : [ ],
				"built_in_analyzers" : [ ]
			  }
			},
			"nodes" : {
			  "count" : {
				"total" : 3,
				"coordinating_only" : 0,
				"data" : 3,
				"data_cold" : 3,
				"data_content" : 3,
				"data_hot" : 3,
				"data_warm" : 3,
				"ingest" : 3,
				"master" : 3,
				"ml" : 3,
				"remote_cluster_client" : 3,
				"transform" : 3,
				"voting_only" : 0
			  },
			  "versions" : [
				"7.10.0"
			  ],
			  "os" : {
				"available_processors" : 12,
				"allocated_processors" : 12,
				"names" : [
				  {
					"name" : "Linux",
					"count" : 3
				  }
				],
				"pretty_names" : [
				  {
					"pretty_name" : "CentOS Linux 8 (Core)",
					"count" : 3
				  }
				],
				"mem" : {
				  "total_in_bytes" : 25013551104,
				  "free_in_bytes" : 251179008,
				  "used_in_bytes" : 24762372096,
				  "free_percent" : 1,
				  "used_percent" : 99
				}
			  },
			  "process" : {
				"cpu" : {
				  "percent" : 1
				},
				"open_file_descriptors" : {
				  "min" : 304,
				  "max" : 304,
				  "avg" : 304
				}
			  },
			  "jvm" : {
				"max_uptime_in_millis" : 31334,
				"versions" : [
				  {
					"version" : "15.0.1",
					"vm_name" : "OpenJDK 64-Bit Server VM",
					"vm_version" : "15.0.1+9",
					"vm_vendor" : "AdoptOpenJDK",
					"bundled_jdk" : true,
					"using_bundled_jdk" : true,
					"count" : 3
				  }
				],
				"mem" : {
				  "heap_used_in_bytes" : 746485240,
				  "heap_max_in_bytes" : 3221225472
				},
				"threads" : 93
			  },
			  "fs" : {
				"total_in_bytes" : 481133838336,
				"free_in_bytes" : 454960951296,
				"available_in_bytes" : 430429470720
			  },
			  "plugins" : [ ],
			  "network_types" : {
				"transport_types" : {
				  "netty4" : 3
				},
				"http_types" : {
				  "netty4" : 3
				}
			  },
			  "discovery_types" : {
				"zen" : 3
			  },
			  "packaging_types" : [
				{
				  "flavor" : "default",
				  "type" : "docker",
				  "count" : 3
				}
			  ],
			  "ingest" : {
				"number_of_pipelines" : 2,
				"processor_stats" : {
				  "gsub" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  },
				  "script" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  }
				}
			  }
			}
		  }`,
		ExpectedNodesStatsFailures: 0,
	},
	// #6 7.10.1 failed nodes
	{
		Body: `{
			"_nodes" : {
			  "total" : 3,
			  "successful" : 1,
			  "failed" : 2,
			  "failures" : [
				{
				  "type" : "failed_node_exception",
				  "reason" : "Failed node [fkiRoRjkQY-oij3qaBW9CQ]",
				  "node_id" : "fkiRoRjkQY-oij3qaBW9CQ",
				  "caused_by" : {
					"type" : "node_not_connected_exception",
					"reason" : "[es2][172.25.0.3:9300] Node not connected"
				  }
				},
				{
				  "type" : "failed_node_exception",
				  "reason" : "Failed node [-fP8gimVQLu2tQV8luneBQ]",
				  "node_id" : "-fP8gimVQLu2tQV8luneBQ",
				  "caused_by" : {
					"type" : "node_not_connected_exception",
					"reason" : "[es3][172.25.0.4:9300] Node not connected"
				  }
				}
			  ]
			},
			"cluster_name" : "elasticsearch",
			"cluster_uuid" : "XvRExZcTTqiiYa8Qa0E6LA",
			"timestamp" : 1625645771492,
			"indices" : {
			  "count" : 0,
			  "shards" : { },
			  "docs" : {
				"count" : 0,
				"deleted" : 0
			  },
			  "store" : {
				"size_in_bytes" : 0,
				"reserved_in_bytes" : 0
			  },
			  "fielddata" : {
				"memory_size_in_bytes" : 0,
				"evictions" : 0
			  },
			  "query_cache" : {
				"memory_size_in_bytes" : 0,
				"total_count" : 0,
				"hit_count" : 0,
				"miss_count" : 0,
				"cache_size" : 0,
				"cache_count" : 0,
				"evictions" : 0
			  },
			  "completion" : {
				"size_in_bytes" : 0
			  },
			  "segments" : {
				"count" : 0,
				"memory_in_bytes" : 0,
				"terms_memory_in_bytes" : 0,
				"stored_fields_memory_in_bytes" : 0,
				"term_vectors_memory_in_bytes" : 0,
				"norms_memory_in_bytes" : 0,
				"points_memory_in_bytes" : 0,
				"doc_values_memory_in_bytes" : 0,
				"index_writer_memory_in_bytes" : 0,
				"version_map_memory_in_bytes" : 0,
				"fixed_bit_set_memory_in_bytes" : 0,
				"max_unsafe_auto_id_timestamp" : -9223372036854775808,
				"file_sizes" : { }
			  },
			  "mappings" : {
				"field_types" : [ ]
			  },
			  "analysis" : {
				"char_filter_types" : [ ],
				"tokenizer_types" : [ ],
				"filter_types" : [ ],
				"analyzer_types" : [ ],
				"built_in_char_filters" : [ ],
				"built_in_tokenizers" : [ ],
				"built_in_filters" : [ ],
				"built_in_analyzers" : [ ]
			  }
			},
			"nodes" : {
			  "count" : {
				"total" : 1,
				"coordinating_only" : 0,
				"data" : 1,
				"data_cold" : 1,
				"data_content" : 1,
				"data_hot" : 1,
				"data_warm" : 1,
				"ingest" : 1,
				"master" : 1,
				"ml" : 1,
				"remote_cluster_client" : 1,
				"transform" : 1,
				"voting_only" : 0
			  },
			  "versions" : [
				"7.10.0"
			  ],
			  "os" : {
				"available_processors" : 4,
				"allocated_processors" : 4,
				"names" : [
				  {
					"name" : "Linux",
					"count" : 1
				  }
				],
				"pretty_names" : [
				  {
					"pretty_name" : "CentOS Linux 8 (Core)",
					"count" : 1
				  }
				],
				"mem" : {
				  "total_in_bytes" : 8337850368,
				  "free_in_bytes" : 2907955200,
				  "used_in_bytes" : 5429895168,
				  "free_percent" : 35,
				  "used_percent" : 65
				}
			  },
			  "process" : {
				"cpu" : {
				  "percent" : 0
				},
				"open_file_descriptors" : {
				  "min" : 253,
				  "max" : 253,
				  "avg" : 253
				}
			  },
			  "jvm" : {
				"max_uptime_in_millis" : 59624,
				"versions" : [
				  {
					"version" : "15.0.1",
					"vm_name" : "OpenJDK 64-Bit Server VM",
					"vm_version" : "15.0.1+9",
					"vm_vendor" : "AdoptOpenJDK",
					"bundled_jdk" : true,
					"using_bundled_jdk" : true,
					"count" : 1
				  }
				],
				"mem" : {
				  "heap_used_in_bytes" : 308600832,
				  "heap_max_in_bytes" : 1073741824
				},
				"threads" : 35
			  },
			  "fs" : {
				"total_in_bytes" : 160377946112,
				"free_in_bytes" : 151653855232,
				"available_in_bytes" : 143476695040
			  },
			  "plugins" : [ ],
			  "network_types" : {
				"transport_types" : {
				  "netty4" : 1
				},
				"http_types" : {
				  "netty4" : 1
				}
			  },
			  "discovery_types" : {
				"zen" : 1
			  },
			  "packaging_types" : [
				{
				  "flavor" : "default",
				  "type" : "docker",
				  "count" : 1
				}
			  ],
			  "ingest" : {
				"number_of_pipelines" : 1,
				"processor_stats" : {
				  "gsub" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  },
				  "script" : {
					"count" : 0,
					"failed" : 0,
					"current" : 0,
					"time_in_millis" : 0
				  }
				}
			  }
			}
		  }`,
		ExpectedNodesStatsFailures: 2,
	},
}

func TestClusterStatsErrorResponse(t *testing.T) {
	for i, tt := range clusterStatsErrorResponseTests {
		var resp ClusterStatsResponse
		if err := json.Unmarshal([]byte(tt.Body), &resp); err != nil {
			t.Fatal(err)
		}
		if want, have := tt.ExpectedNodesStatsFailures, len(resp.NodesStats.Failures); want != have {
			t.Fatalf("case #%d: expected %d errors, got %d", i, want, have)
		}
	}
}
