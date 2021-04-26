package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	olc "github.com/google/open-location-code/go"
	logging "github.com/op/go-logging"
	ipfsShell "github.com/romanblanco/go-ipfs-api"
	extractor "github.com/romanblanco/graffiti-ipfs/extractor"
	exif "github.com/rwcarlsen/goexif/exif"
)

var debugLog = logging.MustGetLogger("main")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} — %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type LatLon struct {
	Value   float64
	NotNull bool
}

// Graffiti structure describes a graffiti photo stored in IPFS
type Graffiti struct {
	Name       string    `json:"name"`
	Ipfs       string    `json:"ipfs"`
	Collection string    `json:"collection"`
	Date       time.Time `json:"date,omitempty"`
	Latitude   LatLon    `json:"latitude,omitempty"`
	Longitude  LatLon    `json:"longitude,omitempty"`
	Olc        string    `json:"olc"`
	Surface    string    `json:"surface"`
	Tags       []string  `json:"tags"`
}

type GraffitiSet []Graffiti

// TODO: Marker properties and geometry should be implemented on frontend side.
//       The only responsibility of backend should be to provide necessary
//       data.
type MarkerProperties struct {
	Ipfs         string    `json:"ipfs"`
	Collection   string    `json:"collection"`
	Surface      string    `json:"surface"`
	Date         time.Time `json:"date"`
	Latitude     LatLon    `json:"latitude"`
	Longitude    LatLon    `json:"longitude"`
	Olc          string    `json:"olc"`
	Tags         []string  `json:"tags"`
	MarkerSymbol string    `json:"marker-symbol"`
	MarkerColor  string    `json:"marker-color"`
	MarkerSize   string    `json:"marker-size"`
}

type MarkerGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type GeoJsonFeature struct {
	Type       string           `json:"type"`
	Geometry   MarkerGeometry   `json:"geometry"`
	Properties MarkerProperties `json:"properties"`
}

type GeoJsonCollection struct {
	Type     string           `json:"type"`
	Features []GeoJsonFeature `json:"features"`
}

// TODO: IPFS_CONTENT should be an array of IPFS content hashes to use as a
//       source of photos.

const IPFS_CONTENT string = "QmYa8Hi5dtahzUvqBN5orjFhsMyxcyQKefoiCGGmezooQ4"
const TFile int = 2 // go-ipfs-api/shell.go constant describing file type

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	debugLog.Infof("parsing graffiti metadata JSON file")
	descriptionJson, err := ioutil.ReadFile("./graffiti.json")
	if err != nil {
		debugLog.Errorf("error reading json file: %s\n", err)
	}

	var descriptionFromJSON GraffitiSet
	err = json.Unmarshal(descriptionJson, &descriptionFromJSON)
	if err != nil {
		debugLog.Criticalf("%s\n", err)
		panic("parsing JSON failed")
	}
	debugLog.Debugf("parsed JSON with %v elements", len(descriptionFromJSON))

	//sh := ipfsShell.NewShell("0.0.0.0:5001")
	sh := ipfsShell.NewShell("ipfs:5001")
	// This works locally
	//sh := ipfsShell.NewShell("127.0.0.1:5001")

	debugLog.Infof("getting IPFS content")
	// TODO: timeout https://github.com/tumregels/Network-Programming-with-Go/blob/master/socket/controlling_tcp_connections.md#timeout
	photoMetadata, err := sh.List(IPFS_CONTENT)
	debugLog.Debugf("got metadata of %v photos from IPFS", len(photoMetadata))
	descriptionFromIPFS := GraffitiSet{}

	rawTar, err := sh.GetRawTar(IPFS_CONTENT)
	debugLog.Debugf("got raw tar of photos from IPFS: %T: %v", rawTar, rawTar)
	extractorInstance := extractor.New(rawTar)
	if err != nil {
		debugLog.Debugf("err:", err)
		panic("error while extracting raw tar")
	}

	debugLog.Info("parsing EXIF data from photos")
	for _, photo := range photoMetadata {

		debugLog.Debugf("parsing EXIF data for photo: %v", photo.Name)
		if photo.Type != TFile {
			debugLog.Errorf("not a file, skipping %s\n", photo.Name)
			continue
		}

		finfo, freader, err := extractorInstance.Next()
		if err != nil {
			panic("error while loading following record")
		}

		if photo.Name != finfo.Name() {
			panic("name mismatch")
		}

		exifData, err := exif.Decode(freader)
		if err != nil {
			debugLog.Errorf("error decoding exif metadata: %s\n", err)
		}

		var latitude, longitude LatLon
		var openLocCode string
		lat, lon, err := exifData.LatLong()
		if err != nil {
			debugLog.Errorf("photo %s has no coordinates\n", photo.Hash)
			latitude = LatLon{Value: 0, NotNull: false}
			longitude = LatLon{Value: 0, NotNull: false}
			openLocCode = ""
		} else {
			latitude = LatLon{Value: lat, NotNull: true}
			longitude = LatLon{Value: lon, NotNull: true}
			openLocCode = olc.Encode(lat, lon, 16)
		}

		date, err := exifData.DateTime()
		if err != nil {
			debugLog.Errorf("can not parse date: %s\n", err)
		}

		meta := Graffiti{
			Name:       photo.Name,
			Ipfs:       photo.Hash,
			Date:       date,
			Olc:        openLocCode,
			Latitude:   latitude,
			Longitude:  longitude,
			Surface:    "",
			Collection: IPFS_CONTENT,
			Tags:       make([]string, 0),
		}
		descriptionFromIPFS = append(descriptionFromIPFS, meta)
	}

	debugLog.Info("merging data from IPFS with metadata from JSON")
	result, ipfsExtra, metadataExtra := merge(descriptionFromIPFS, descriptionFromJSON)

	debugLog.Infof("loaded slice len %d\n", len(result))
	debugLog.Infof("extra records in IPFS: %d\n", len(ipfsExtra))
	debugLog.Infof("extra records in metadata: %d\n", len(metadataExtra))

	export, err := json.Marshal(result)
	if err != nil {
		debugLog.Errorf("error creating JSON: %s", err)
	}

	apiHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(export))
	}

	geoJsonHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, geoJson(result))
	}

	debugLog.Info("serving complete content at :8083/api")
	debugLog.Info("serving geotagged content at :8083/geojson")
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/geojson", geoJsonHandler)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

// TODO: This should be implemented on frontend.
//
// geoJson function takes JSON data describing photos and enriches them by
// mapbox properties for a marker.
func geoJson(photos GraffitiSet) (jsonData string) {
	markers := []GeoJsonFeature{}

	for _, photo := range photos {

		if photo.Olc == "" {
			continue
		}

		markerColor := "#000000"
		if photo.Surface == "" {
			markerColor = "#0088ce"
		}

		properties := MarkerProperties{
			Ipfs:         photo.Ipfs,
			Surface:      photo.Surface,
			Collection:   photo.Collection,
			Date:         photo.Date,
			Latitude:     photo.Latitude,
			Longitude:    photo.Longitude,
			Olc:          photo.Olc,
			Tags:         photo.Tags,
			MarkerSymbol: "art-gallery",
			MarkerColor:  markerColor,
			MarkerSize:   "medium",
		}

		geometry := MarkerGeometry{
			Type:        "Point",
			Coordinates: []float64{photo.Longitude.Value, photo.Latitude.Value},
		}

		marker := GeoJsonFeature{
			Type:       "Feature",
			Geometry:   geometry,
			Properties: properties,
		}

		markers = append(markers, marker)
	}

	collection := GeoJsonCollection{
		Type:     "FeatureCollection",
		Features: markers,
	}

	jsonMarkers, _ := json.Marshal(collection)
	jsonData = string(jsonMarkers)
	return
}

func merge(ipfsDataSlice, metadataSlice GraffitiSet) (united, ipfsExtras, metadataExtras GraffitiSet) {
	var ipfsCounter int = 0
	var metaCounter int = 0

	ipfsDataSlice.sortByIpfsHash()
	metadataSlice.sortByIpfsHash()

	for {
		if metaCounter == len(metadataSlice) && ipfsCounter == len(ipfsDataSlice) {
			return
		}
		if ipfsCounter == len(ipfsDataSlice) {
			extras := metadataSlice[metaCounter : len(metadataSlice)-1]
			united = append(united, extras...)
			metadataExtras = append(metadataExtras, extras...)
			break
		}
		if metaCounter == len(metadataSlice) {
			extras := ipfsDataSlice[ipfsCounter : len(ipfsDataSlice)-1]
			united = append(united, extras...)
			ipfsExtras = append(ipfsExtras, extras...)
			break
		}

		if ipfsDataSlice[ipfsCounter].Ipfs == metadataSlice[metaCounter].Ipfs {
			// enrich IPFS data with metadata attributes
			ipfsDataSlice[ipfsCounter].Surface = metadataSlice[metaCounter].Surface
			ipfsDataSlice[ipfsCounter].Tags = metadataSlice[metaCounter].Tags
			united = append(united, ipfsDataSlice[ipfsCounter])
			ipfsCounter += 1
			metaCounter += 1
		} else if ipfsDataSlice[ipfsCounter].Ipfs < metadataSlice[metaCounter].Ipfs {
			united = append(united, ipfsDataSlice[ipfsCounter])
			ipfsExtras = append(ipfsExtras, ipfsDataSlice[ipfsCounter])
			ipfsCounter += 1
		} else {
			united = append(united, metadataSlice[metaCounter])
			metadataExtras = append(metadataExtras, metadataSlice[metaCounter])
			metaCounter += 1
		}
	}

	return
}

func (set GraffitiSet) sortByIpfsHash() {
	sort.Slice(set, func(i, j int) bool {
		return set[i].Ipfs < set[j].Ipfs
	})
}

func (i LatLon) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	// the key is not set to null
	var tmp float64
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	i.Value = tmp
	i.NotNull = true
	return nil
}

func (i LatLon) MarshalJSON() ([]byte, error) {
	if i.NotNull {
		return json.Marshal(i.Value)
	}

	return []byte("null"), nil
}
