# Graffiti

Repository provides a tool for graffiti community to keep track and preserving street artwork in their areas.

<img src="./docs/map.png" width="100%" />

## Features

### Distributed storage

The core of the project are geotagged photos stored in [IPFS](), meaning they are not stored in a central database, allowing users to only maintain the photos they are interested in. The data will be reachable as long as there is someone sharing them.

### Metadata

The maintained graffiti collection can be annotated by the user in [metadata](https://github.com/romanblanco/graffiti/blob/master/collection/graffiti.json), specifying OpenStreetMap node as a surface and tags for photos for better searchability.

## Requirements

- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting started

```
docker-compose up -d --build
docker-compose up
```

The map should be available on http://localhost:4567/

## Troubleshooting

If the `graffiti-collection` container is hanging on `getting IPFS content`, it is probably because the route to content has not been discovered. To help the discovery, run:

```bash
docker exec -it graffiti-ipfs sh
ipfs ping QmZCtha6AHsaNRV1LScMXLewjx9M2imPTxKaP2ty2TJ219
exit
```

The `graffiti-collection` container should now successfully retrieve the photos as blobs from IPFS into `data` folder, and they should be visible on the map.

## Learn more

- [Wiki](https://github.com/romanblanco/graffiti/wiki)
- [Task board](https://github.com/romanblanco/graffiti/projects)
- [#graffiti:matrix.org](https://view.matrix.org/room/!WVlsowqbtSqMWutaCX:matrix.org/)
