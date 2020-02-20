import React from "react";
import { geolocated } from "react-geolocated";
import MapTileLayerCluster from "./MapTileLayerCluster";

class FoundLocationMap extends React.Component {
    render() {
        const crd = this.props.isGeolocationEnabled && this.props.coords ? this.props.coords: null;

        return !this.props.isGeolocationAvailable ? (
            <div>Your browser does not support Geolocation</div>
        ) : <MapTileLayerCluster crd={crd} />;
    }
}
 
export default geolocated({
    positionOptions: {
        enableHighAccuracy: false,
    },
    userDecisionTimeout: 5000,
})(FoundLocationMap);
