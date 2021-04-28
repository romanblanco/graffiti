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

type CollectionsSource struct {
	Collections []string `json:"collections"`
	Metadata    []string `json:"metadata"`
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

type GraffitiProperties struct {
	Ipfs       string    `json:"ipfs"`
	Collection string    `json:"collection"`
	Surface    string    `json:"surface"`
	Date       time.Time `json:"date"`
	Latitude   LatLon    `json:"latitude"`
	Longitude  LatLon    `json:"longitude"`
	Olc        string    `json:"olc"`
	Tags       []string  `json:"tags"`
}

type GraffitiPoint struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type GraffitiFeature struct {
	Type       string             `json:"type"`
	Geometry   GraffitiPoint      `json:"geometry"`
	Properties GraffitiProperties `json:"properties"`
}

type GraffitiFeatureCollection struct {
	Type     string            `json:"type"`
	Features []GraffitiFeature `json:"features"`
}

const TFile int = 2 // go-ipfs-api/shell.go constant describing file type

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	debugLog.Infof("parsing collection source file source.json")
	sourceJson, err := ioutil.ReadFile("./source.json")
	if err != nil {
		debugLog.Errorf("i/o error reading source.json: %s\n", err)
	}
	var source CollectionsSource
	err = json.Unmarshal(sourceJson, &source)
	if err != nil {
		debugLog.Criticalf("%s\n", err)
		panic("parsing source.json failed")
	}

	// TODO: the metadata file should be obtained "metadata" object in
	// source.json and should have an option to
	// - read local file
	// - fetch file from URL
	// - fetch file from IPFS
	debugLog.Infof("parsing graffiti metadata file metadata.json")
	metadataJson, err := ioutil.ReadFile("./metadata.json")
	if err != nil {
		debugLog.Errorf("i/o error reading metadata.json: %s\n", err)
	}
	var metadata GraffitiSet
	err = json.Unmarshal(metadataJson, &metadata)
	if err != nil {
		debugLog.Criticalf("%s\n", err)
		panic("parsing metadata.json failed")
	}
	debugLog.Debugf("parsed %v records from metadata.json", len(metadata))

	debugLog.Infof("connecting to IPFS")
	sh := ipfsShell.NewShell("0.0.0.0:5001")
	// TODO: timeout
	// https://github.com/tumregels/Network-Programming-with-Go/blob/master/socket/controlling_tcp_connections.md#timeout

	parsedCollection := GraffitiSet{}

	for _, collection := range source.Collections {
		collectionContent, err := sh.List(collection)
		collectionRawTar, err := sh.GetRawTar(collection)
		extractorInstance := extractor.New(collectionRawTar)
		if err != nil {
			debugLog.Debugf("err:", err)
			panic("error while extracting raw tar")
		}

		for _, photo := range collectionContent {
			debugLog.Debugf("parsing EXIF data for photo: %v", photo.Hash)
			if photo.Type != TFile {
				debugLog.Errorf("not a file, skipping %s\n", photo.Hash)
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
				Collection: collection,
				Tags:       make([]string, 0),
			}
			parsedCollection = append(parsedCollection, meta)
		}
	}

	parsedCollection = unique(parsedCollection)

	debugLog.Info("merging data from IPFS with metadata from JSON")
	result, ipfsExtra, metadataExtra := merge(parsedCollection, metadata)

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
	http.HandleFunc("/api", apiHandler)
	debugLog.Info("serving geotagged content at :8083/geojson")
	http.HandleFunc("/geojson", geoJsonHandler)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func geoJson(photos GraffitiSet) (jsonData string) {
	markers := []GraffitiFeature{}

	for _, photo := range photos {

		if photo.Olc == "" {
			continue
		}

		properties := GraffitiProperties{
			Ipfs:       photo.Ipfs,
			Surface:    photo.Surface,
			Collection: photo.Collection,
			Date:       photo.Date,
			Latitude:   photo.Latitude,
			Longitude:  photo.Longitude,
			Olc:        photo.Olc,
			Tags:       photo.Tags,
		}

		geometry := GraffitiPoint{
			Type:        "Point",
			Coordinates: []float64{photo.Longitude.Value, photo.Latitude.Value},
		}

		marker := GraffitiFeature{
			Type:       "Feature",
			Geometry:   geometry,
			Properties: properties,
		}

		markers = append(markers, marker)
	}

	collection := GraffitiFeatureCollection{
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

func unique(items GraffitiSet) GraffitiSet {
	keys := make(map[string]bool)
	uniqueList := GraffitiSet{}
	for _, item := range items {
		if _, value := keys[item.Ipfs]; !value {
			keys[item.Ipfs] = true
			uniqueList = append(uniqueList, item)
		}
	}
	return uniqueList
}
