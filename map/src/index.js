import React from "react";
import ReactDOM from "react-dom";
import mapboxgl from "mapbox-gl";
import { TOKEN } from "./token.js";

mapboxgl.accessToken = TOKEN;

class Application extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      lng: 5,
      lat: 34,
      zoom: 2
    };
  }

  componentDidMount() {
    const map = new mapboxgl.Map({
      container: this.mapContainer,
      style: "mapbox://styles/mapbox/light-v10",
      center: [this.state.lng, this.state.lat],
      zoom: this.state.zoom
    });

    map.on("move", () => {
      this.setState({
        lng: map.getCenter().lng.toFixed(4),
        lat: map.getCenter().lat.toFixed(4),
        zoom: map.getZoom().toFixed(2)
      });
    });

    map.on("load", () => {
      map.addSource("graffiti", {
        type: "geojson",
        data: "http://127.0.0.1:8083/geojson",
      });

      map.addLayer({
        "id": "graffiti",
        "type": "symbol",
        "source": "graffiti",
        "layout": {
          "icon-image": "art-gallery-15",
        },
      });

      map.on("click", "graffiti", function(e) {
        var coordinates = e.features[0].geometry.coordinates.slice();
        var description = e.features[0].properties;

        var popupContent = `
        <a target="_blank" class="popup" href="http://ipfs:8080/ipfs/${description.ipfs}">
          <img src="http://127.0.0.1:8080/ipfs/${description.ipfs}" height="140" width="200" />
        </a>
        <p>Date: ${description.date}</p>
        <p>IPFS:
          <a target="_blank"
             rel="noopener noreferrer"
             href="https://ipfs.io/ipfs/${description.ipfs}">
            ipfs.io
          </a>
        </p>
        <p>Plus: ${description.olc}</p>
        <p>Google Maps:
          <a target="_blank"
             rel="noopener noreferrer"
             href="https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(description.olc)}">
            Plus Code
          </a>,
          <a target="_blank"
             rel="noopener noreferrer"
             href="https://www.google.com/maps/search/?api=1&query=${description.latitude},${description.longitude}">
            GPS
          </a>
        </p>
        <p>Tags: ${description.tags}</p>
        <p>Surface:
          <a target="_blank"
             rel="noopener noreferrer"
             href="https://osm.org/${description.surface}">
            ${description.surface}
          </a>
        </p>
        `

        // Ensure that if the map is zoomed out such that multiple
        // copies of the feature are visible, the popup appears
        // over the copy being pointed to.
        while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
          coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
        }

        new mapboxgl.Popup()
          .setLngLat(coordinates)
          .setHTML(popupContent)
          .addTo(map);
      });

      map.addControl(
        new mapboxgl.GeolocateControl({
          positionOptions: {
            enableHighAccuracy: true
          },
          trackUserLocation: true
        })
      );
    });
  }

  render() {
    return (
      <div>
        <div ref={el => this.mapContainer = el} className="mapContainer"/>
      </div>
    )
  }
}

ReactDOM.render(<Application />, document.getElementById("app"));
