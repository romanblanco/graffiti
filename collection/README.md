### build and run server:

```sh
$ go get github.com/ipfs/go-ipfs-api
$ go get github.com/romanblanco/go-ipfs-api
$ git clone https://github.com/romanblanco/graffiti-ipfs.git
$ cd graffiti-ipfs/
$ go build -o graffiti -ldflags="-s -w" .
$ # TODO: update recources for data and metadata in source.json
$ ./graffiti-ipfs
```

server is running at http://localhost:8083

- `/api` -- JSON output of photos loaded from IPFS enriched by metadata

```json
[
  {
    "name": "IMG_20180610_144211.jpg",
    "ipfs": "QmQPBmvEB32UUwowAHxs7rCFn5q9hAitmowbSB8YuRpbtx",
    "date": "2018-06-10T14:42:12Z",
    "latitude": 49.22980116666667,
    "longitude": 16.529714583333334,
    "olc": "8FXR6GHH+WVG575V",
    "surface": "way/737085945",
    "tags": [
      "bufet"
    ]
  },
  {
    "name": "IMG_20180728_170515.jpg",
    "ipfs": "QmNoQwa8wFKq5gW2v4Bg4eqcNqNi1hPnWHyQjup4PPfoNW",
    "date": "2018-07-28T17:05:16Z",
    "latitude": 49.22021102777778,
    "longitude": 16.640451416666668,
    "olc": "8FXR6JCR+35PF924",
    "surface": "",
    "tags": [
      "gee",
      "head"
    ]
  },
  ...
]
```

- `/geojson` -- GeoJSON source for graffiti-map

```json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          14.460802055555554,
          50.09414288888889
        ]
      },
      "properties": {
        "ipfs": "QmNLfYnfLEVnb13qMoMFpqYybUpowEVVg1K5hp2SfdW2eg",
        "surface": "",
        "date": "2018-04-29T12:48:52Z",
        "latitude": 50.09414288888889,
        "longitude": 14.460802055555554,
        "olc": "9F2P3FV6+M83PGWG",
        "tags": [],
        "marker-symbol": "art-gallery",
        "marker-color": "#0088ce",
        "marker-size": "medium"
      }
    },
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          16.628084166666667,
          49.19617841666666
        ]
      },
      "properties": {
        "ipfs": "QmNP5Bk7PCmX3ynK9fo2UbMpnKVZGbrNpuqzMmM7yDNfbm",
        "surface": "",
        "date": "2018-05-12T11:16:10Z",
        "latitude": 49.19617841666666,
        "longitude": 16.628084166666667,
        "olc": "8FXR5JWH+F6G4QC3",
        "tags": [],
        "marker-symbol": "art-gallery",
        "marker-color": "#0088ce",
        "marker-size": "medium"
      }
    },
    ...

  ]
}
```
