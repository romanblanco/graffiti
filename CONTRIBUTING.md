# Share geotagged photos of graffiti and metadata

The core of this project are photos that contain GPS location and date in the photo EXIF metadata.
In case you have photos of graffiti but with no GPS in EXIF metadata, if you know the location, the metadata can be provided additionally using tools described on [OpenStreetMap Wiki](https://wiki.openstreetmap.org/wiki/Geotagging_Source_Photos#Geotagging_photos_from_a_GPS_tracklog).

To make the collection work as a distributed system, the project uses IPFS to for sharing photos.
To share your collection of photos follow these steps:

  1. Install [IPFS](https://ipfs.io/#install)
  2. Upload folder with photos to IPFS (`ipfs add -r /path/to/folder`)
  3. Share the hash of the uploaded folder

To help mapping the graffiti in the area of interest, the collection can use metadata describing a OpenStreetMap node of surface (wall) the graffiti is sprayed on and tags for easier searchability.

# Codebase

The codebase of the project consist of two modules:

- Golang server that provides GeoJSON from provided IPFS directory
- React application using mapbox to plot the data

The plans for both modules can be found in https://github.com/romanblanco/graffiti/projects/1.

Any help is very appreciated.
