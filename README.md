# Graffiti

Graffiti repository provides tools for any graffiti community around the world to keep track of street artwork in their areas.

<img src="./docs/map.png" width="100%" />

## Features

### Distributed storage

The core of the project are geotagged photos stored in IPFS, meaning they are not stored in a central database, allowing users to only maintain the photos they are interested in. The data will be reachable as long as there is someone sharing them.

### Metadata

The maintained graffiti collection can be described by the users in [metadata](https://github.com/romanblanco/graffiti/blob/master/collection/graffiti.json), specifying OpenStreetMap node and tags for photos.

## Requirements

- Go (https://golang.org/)
- node.js (https://nodejs.org/)
- IPFS (https://ipfs.io/)

## Getting started

### build and run server:

```sh
$ cd collection/
$ # TODO: update recources for data and metadata in source.json
$ # go get <libraries>
$ go build -o graffiti -ldflags="-s -w" .
$ ./graffiti
```

server is running at http://localhost:8083

### run map client:

```sh
$ cd map/
$ # update your Mapbox token in src/token.js
$ npm install
$ npm start
```

map is running at http://localhost:3000


## Learn more

- [Wiki](https://github.com/romanblanco/graffiti/wiki)
- [Task board](https://github.com/romanblanco/graffiti/projects)
