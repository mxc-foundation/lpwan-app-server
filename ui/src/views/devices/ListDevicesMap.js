import L from "leaflet";
import "leaflet.awesome-markers";
import moment from "moment";
import React, { Component } from "react";
import { Map, Marker, Popup } from 'react-leaflet';
import MarkerClusterGroup from "react-leaflet-markercluster";
import { Link } from "react-router-dom";
import MapTileLayer from "../../components/MapTileLayer";
import DeviceStore from "../../stores/DeviceStore";


class ListDevicesMap extends Component {
  constructor() {
    super();

    this.state = {
      items: null,
    };
  }

  componentDidMount() {
    const filters = {
      limit: 999,
      offset: 0,
      organizationID: this.props.organizationID,
      search: "",
    };
    DeviceStore.list(filters, resp => {
      this.setState({
        items: resp.result,
      });
    });
  }

  render() {
    const currentOrgID = this.props.organizationID;
    const currentApplicationID = this.props.applicationID;

    if (this.state.items === null) {
      return null;
    }

    const style = {
      height: 800,
    };


    let bounds = [];
    let markers = [];

    const greenMarker = L.AwesomeMarkers.icon({
      icon: "wifi",
      prefix: "fa",
      markerColor: "green",
    });

    const grayMarker = L.AwesomeMarkers.icon({
      icon: "wifi",
      prefix: "fa",
      markerColor: "gray",
    });

    const redMarker = L.AwesomeMarkers.icon({
      icon: "wifi",
      prefix: "fa",
      markerColor: "red",
    });

    for (const item of this.state.items) {
      const position = [item.location.latitude, item.location.longitude];

      bounds.push(position);

      let marker = greenMarker;
      let lastSeen = "";

      if (item.lastSeenAt === undefined || item.lastSeenAt === null) {
        marker = grayMarker;
        lastSeen = "Never seen online";
      } else {
        const ts = moment(item.lastSeenAt);
        if (ts.isBefore(moment().subtract(5, 'minutes'))) {
          marker = redMarker;
        }

        lastSeen = ts.fromNow();
      }

      markers.push(
        <Marker position={position} key={`gw-${item.devEUI}`} icon={marker}>
          <Popup>
            {currentApplicationID 
              ? <Link to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${item.devEUI}`}>{item.name}</Link>
              : <Link to={`/organizations/${currentOrgID}/devices/${item.devEUI}`}>{item.name}</Link>
            }

            <br />
            {item.devEUI}<br /><br />
            {lastSeen}
          </Popup>
        </Marker>
      );
    }

    return (<React.Fragment>
      {bounds && <Map bounds={bounds} maxZoom={19} style={style} animate={true} scrollWheelZoom={false}>
        <MapTileLayer />
        <MarkerClusterGroup>
          {markers}
        </MarkerClusterGroup>
      </Map>}
    </React.Fragment>
    );
  }
}

export default ListDevicesMap;
