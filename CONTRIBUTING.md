# Take a geotagged photo

The initial part of this project are geotagged photos. The date and GPS location are extracted from the photos EXIF data to later plot them on map.

Make sure your device has photo geotagging enabled ([Android](https://support.apple.com/guide/iphone/control-the-location-information-you-share-iph3dd5f9be/14.0/ios/14.0), [iPhone](https://support.apple.com/guide/iphone/control-the-location-information-you-share-iph3dd5f9be/14.0/ios/14.0)).

In case you have photos of graffiti but with no GPS in EXIF metadata, if you know the location, the metadata can be provided additionally using tools described on [OpenStreetMap Wiki](https://wiki.openstreetmap.org/wiki/Geotagging_Source_Photos#Geotagging_photos_from_a_GPS_tracklog).

To make sure your photo has needed metadata, you can use [Jeffrey's Image Metadata Viewer](http://exif.regex.info/exif.cgi).

# Provide photos on distributed network

To make the collection work as a distributed system, the project uses IPFS to for sharing photos.

This project already [sets up an IPFS node](https://github.com/romanblanco/graffiti/blob/master/docker-compose.yml#L12) that can be used for providing photos.
After building the Docker containers, it will accessible on [localhost:5001](http://127.0.0.1:5001/webui).

## Upload photos using IPFS Desktop

1. Navigate to Files
2. Press Import button and choose Folder
3. Select foler to upload and submit.
4. Using three dotcs, you can copy CID of folder to share it.

[read more](https://github.com/ipfs/ipfs-desktop/#quickly-import-files-folders-and-screenshots-to-ipfs)

## Upload photos using CLI

1. Use `ipfs add --recursive /path/to/folder`
2. Share the hash of the uploaded folder

[read more](https://docs.ipfs.io/reference/cli/#ipfs-add)

# Metadata

To help mapping the graffiti in the area of interest, the collection can use [manually provided metadata](https://github.com/romanblanco/graffiti/blob/master/collection/metadata.yaml) describing a OpenStreetMap node of a surface (wall) the graffiti is sprayed on and tags for easier searchability.

# Codebase

The codebase of the project consist of two modules:

- Golang server generating GeoJSON from the data provided by IPFS directory
- React application using mapbox to plot the data

The plans for both modules can be found in https://github.com/romanblanco/graffiti/projects/1.

Any help is very appreciated.
