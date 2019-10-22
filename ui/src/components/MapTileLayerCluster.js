import React, { Component } from 'react';

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
            center: [51.1657, 10.4515],
            zoom: 6,
            zoomContorl: false
        });

        L.tileLayer('https://cartodb-basemaps-{s}.global.ssl.fastly.net/rastertiles/voyager/{z}/{x}/{y}{r}.png',{
            detectRetina: true,
            maxZoom: 20,
            maxNativeZoom: 17
        }).addTo(this.map);

        var markers = L.markerClusterGroup({
          showCoverageOnHover: true,
        });
        
        const res = this.loadData();
        res.then(resp => {
            for(let i=0;i<resp.result.length; i++){
                markers.addLayer(L.marker([resp.result[i].location.latitude, resp.result[i].location.longitude]));
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
