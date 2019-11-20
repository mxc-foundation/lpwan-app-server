import React, { Component } from 'react';

import { MAP_LAYER,GATEWAY_ICON } from '../util/Data'
import GatewayStore from '../stores/GatewayStore';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import styled from 'styled-components';
import 'leaflet.markercluster/dist/MarkerCluster.css'
import 'leaflet.markercluster/dist/MarkerCluster.Default.css'
import 'leaflet.markercluster/dist/leaflet.markercluster.js'
import { withRouter } from "react-router-dom";

const Wrapper = styled.div`
    width: ${props => props.width};
    height: ${props => props.height};
`
function loadGatewayData() {
    return new Promise((resolve, reject) => {
        GatewayStore.listLocations(
        resp => {
          resolve(resp);
        })
    });
} 
class MapTileLayerCluster extends Component {
    componentDidMount(){
        this.map = L.map('map', {
            //preferCanvas: true,
            center: [35.8617, 104.1954],
            zoom: 6,
            zoomContorl: false
        });

        L.tileLayer(MAP_LAYER,{
            detectRetina: true,
            maxZoom: 20,
            maxNativeZoom: 17
        }).addTo(this.map);
        
        let greenIcon = L.icon({
          iconUrl: GATEWAY_ICON,
          //shadowUrl: '',
        
          iconSize:     [30, 56], // size of the icon
          shadowSize:   [50, 64], // size of the shadow
          iconAnchor:   [22, 94], // point of the icon which will correspond to marker's location
          shadowAnchor: [4, 62],  // the same for the shadow
          popupAnchor:  [-3, -76] // point from which the popup should open relative to the iconAnchor
        });

        var markers = L.markerClusterGroup({
          showCoverageOnHover: true,
        });
        
        const res = this.loadData();
        res.then(resp => {
            for(let i=0;i<resp.result.length; i++){
                markers.addLayer(L.marker([resp.result[i].location.latitude, resp.result[i].location.longitude],{icon: greenIcon}));
            }
        })
        
        this.map.addLayer(markers);
    }

    loadData = async () => {
        try {
            var result = await loadGatewayData();
            return result;
        } catch (error) {
            this.setState({loading: false})
            console.error(error);
            this.setState({ error });
        }
    }
    
  render() {
    return(
      <Wrapper
        width="100%" height='100%' id='map'
      />
    )
  }
}

export default withRouter(MapTileLayerCluster);
