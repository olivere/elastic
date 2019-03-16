package geo

// Examples from RFC7946
// https://tools.ietf.org/html/rfc7946

var point = []byte(`
{
  "type":"point",
  "coordinates":[100.0, 0.0]
}
`)

var multiPoint = []byte(`
{
  "type":"multipoint",
  "coordinates":[
    [100.0, 0.0],
    [101.0, 1.0]
  ]
}
`)

var lineString = []byte(`
{
  "type":"linestring",
  "coordinates":[
    [100.0, 0.0],
    [101.0, 1.0]
  ]
}
`)

var multiLineString = []byte(`

{
  "type":"multilinestring",
  "coordinates":[
    [
      [100.0, 0.0],
      [101.0, 1.0]
    ],
    [
      [102.0, 2.0],
      [103.0, 3.0]
    ]
  ]
}
`)

var polygon = []byte(`
{
  "type":"polygon",
  "coordinates":[
    [
      [100.0, 0.0],
      [101.0, 0.0],
      [101.0, 1.0],
      [100.0, 1.0],
      [100.0, 0.0]
    ]
  ]
}
`)

var polygonWithHoles = []byte(`
{
  "type":"polygon",
  "coordinates":[
    [
      [100.0, 0.0],
      [101.0, 0.0],
      [101.0, 1.0],
      [100.0, 1.0],
      [100.0, 0.0]
    ],
    [
      [100.2, 0.2],
      [100.8, 0.2],
      [100.8, 0.8],
      [100.2, 0.8],
      [100.2, 0.2]
    ]
  ]
}

`)

var multiPolygon = []byte(`
{
  "type":"multipolygon",
  "coordinates":[
    [
      [
        [102.0, 2.0],
        [103.0, 2.0],
        [103.0, 3.0],
        [102.0, 3.0],
        [102.0, 2.0]
      ]
    ],
    [
      [
        [100.0, 0.0],
        [101.0, 0.0],
        [101.0, 1.0],
        [100.0, 1.0],
        [100.0, 0.0]
      ],
      [
        [100.2, 0.2],
        [100.8, 0.2],
        [100.8, 0.8],
        [100.2, 0.8],
        [100.2, 0.2]
      ]
    ]
  ]
}
`)

var geometryCollection = []byte(`
{
  "type":"geometrycollection",
  "geometries":[
    {
      "type":"point",
      "coordinates":[100.0, 0.0]
    },
    {
      "type":"linestring",
      "coordinates":[
        [101.0, 0.0],
        [102.0, 1.0]
      ]
    }
  ]
}
`)

var envelope = []byte(`
{
  "type":"envelope",
  "coordinates":[
	[100.0, 1.0],
	[101.0, 0.0]
  ]
}
`)

var circle = []byte(`
{
  "type":"circle",
  "radius":"25m",
  "coordinates":[-109.874838, 44.439550]
}
`)
