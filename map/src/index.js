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
      style: "mapbox://styles/mapbox/dark-v10",
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

      // TODO: {{
      // I want to fetch 8083/api and convert it into geojson here.
      //
      // However, I can only resolve these requests when running on
      // local machine.
      //
      // If I try to load the website from another machine, I can reach the
      // api on 8083/geojson endpoint, but the map can not fetch from these:
      //
      // Fails with ERR_CONNECTION_REFUSED

      fetch('http://0.0.0.0:8083/geojson', {
        credentials: 'same-origin'
      })
      .then(res => res.json())
      .then((out) => {
        console.log('Checkout this JSON! ', out);
      })
      .catch(err => { throw err });

      fetch('http://127.0.0.1:8083/geojson', {
        credentials: 'same-origin'
      })
      .then(res => res.json())
      .then((out) => {
        console.log('Checkout this JSON! ', out);
      })
      .catch(err => { throw err });

      // I am not able to access this anywhere.
      //
      // The same happens on
      // <a ... href="http://ipfs:8080/ipfs/${description.ipfs}" ...>
      //
      // Fails with ERR_NAME_NOT_RESOLVED.
      //
      // https://stackoverflow.com/questions/50299051/how-do-you-resolve-a-docker-alias-using-node/50334529

      fetch('http://collection:8083/geojson', {
        credentials: 'same-origin'
      })
      .then(res => res.json())
      .then((out) => {
        console.log('Checkout this JSON! ', out);
      })
      .catch(err => { throw err });

      // }}

      map.addSource("graffiti", {
        type: "geojson",
        data: "http://127.0.0.1:8083/geojson",
        cluster: true,
      });

      map.addLayer({
        id: "clusters",
        type: "circle",
        source: "graffiti",
        filter: ["has", "point_count"],
        paint: {
          "circle-color": [
            "step",
            ["get", "point_count"],
            "#51bbd6",
            100, "#f1f075",
            750, "#f28cb1"
          ],
          "circle-radius": [
            "step",
            ["get", "point_count"],
            20,
            100, 30,
            750, 40
          ]
        }
      });

      map.addLayer({
        id: "cluster-count",
        type: "symbol",
        source: "graffiti",
        filter: ["has", "point_count"],
        layout: {
          "text-field": "{point_count_abbreviated}",
          "text-size": 10
        }
      });

      map.addLayer({
        id: "unclustered-point",
        type: "circle",
        source: "graffiti",
        filter: ["!", ["has", "point_count"]],
        paint: {
          "circle-color": "#51bbd6",
          "circle-radius": 8,
          "circle-stroke-width": 1,
          "circle-stroke-color": "#51bbd6"
        }
      });

      map.on("click", "clusters", function(e) {
        var features = map.queryRenderedFeatures(e.point, {
          layers: ["clusters"]
        });
        var clusterId = features[0].properties.cluster_id;
        map.getSource("graffiti").getClusterExpansionZoom(
          clusterId,
          function(err, zoom) {
            if (err) return;

            map.easeTo({
              center: features[0].geometry.coordinates,
              zoom: zoom
            });
          }
        );
      });

      map.on("click", "unclustered-point", function(e) {
        var coordinates = e.features[0].geometry.coordinates.slice();
        var description = e.features[0].properties;

        var surface = description.surface === "" ?
            `query?lat=${description.latitude}&lon=${description.longitude}#map=17/${description.latitude}/${description.longitude}` :
          description.surface;
        var surfaceTitle = description.surface === "" ?
          `→ …` :
          description.surface;

        const tags = JSON.parse(description.tags);
        var tagsTitle;
        if (tags === []) {
          tagsTitle = ""
        } else {
          const tagsDiv = tags.map(tag => `<div class="tag tag-gray">${tag}</div>`);
          tagsTitle = tagsDiv.join(" ");
        }

        var tagsElem;
        if (tagsTitle === "") {
          tagsElem = ""
        } else {
          tagsElem = "<i class=\"fas fa-tags\"></i><div class=\"attr\"> " + tagsTitle + "</div>";
        }

        var popupContent = `
        <a target="_blank"
           class="popup"
           href="http://127.0.0.1:8081/ipfs/${description.ipfs}">
          <picture>
            <img src="http://127.0.0.1:8081/ipfs/${description.ipfs}" />
          </picture>
        </a>
        <div class="attr">
          <abbr title="Date"><i class="fas fa-calendar-alt"></i></abbr>
          <time title="${new Date(description.date)}" datetime="${description.date}">
            ${new Date(description.date).toDateString()}
          </time>
        </div>
        <div class="attr">
          <abbr title="InterPlanetary File System"><i class="fas fa-database"></i></abbr>
          <a target="_blank"
             class="tag ipfs"
             title="${description.ipfs}"
             rel="noopener noreferrer"
             href="https://explore.ipld.io/#/explore/${description.ipfs}">
            ${description.ipfs}
          </a>
        </div>
        <div class="attr">
          <abbr title="InterPlanetary File Linked Data"><i class="fas fa-sitemap"></i></abbr>
          <a target="_blank"
             class="tag ipfs"
             title="${description.collection}"
             rel="noopener noreferrer"
             href="https://explore.ipld.io/#/explore/${description.collection}">
            ${description.collection}
          </a>
        </div>
        <div class="attr">
          <abbr title="Open Location Code"><i class="fas fa-map-marked-alt"></i></abbr>
          <a target="_blank"
             class="tag"
             rel="noopener noreferrer"
             href="https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(description.olc)}">
            ${description.olc}
          </a>
        </div>
        <div class="attr">
          <abbr title="OpenStreetMap Node"><i class="fas fa-spray-can"></i></abbr>
          <a target="_blank"
             class="tag"
             rel="noopener noreferrer"
             href="https://osm.org/${surface}">
            ${surfaceTitle}
          </a>
        </div>
        ${tagsElem}
        `

        while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
          coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
        }

        new mapboxgl.Popup()
          .setLngLat(coordinates)
          .setHTML(popupContent)
          .addTo(map);
      });

      map.on("mouseenter", "clusters", function() {
          map.getCanvas().style.cursor = "pointer";
      });
      map.on("mouseleave", "clusters", function() {
          map.getCanvas().style.cursor = "";
      });

      map.addControl(
        new mapboxgl.GeolocateControl({
          positionOptions: {
            enableHighAccuracy: true
          },
          trackUserLocation: true
        })
      );

      map.addControl(new mapboxgl.FullscreenControl());
    });
  }

  render() {
    return (
      <div>
        <div ref={el => this.mapContainer = el} className="mapContainer"/>
      </div>
    );
  }
}

ReactDOM.render(<Application />, document.getElementById("app"));
